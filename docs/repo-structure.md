# Repo Structure

This document explains how gym-saas is organized, which folders own which responsibilities, and the constraints developers should preserve as the project grows.

## High-Level Model

gym-saas is a local-first product with four host surfaces:

- `apps/desktop`
  Wails desktop host for the installed app on the gym PC
- `apps/cloud-api`
  Thin Go cloud layer for public endpoints such as update manifests and future remote APIs
- `apps/web`
  Browser host using React + Vite
- `apps/mobile`
  Expo-based mobile placeholder

Two shared layers sit under those hosts:

- `packages/*`
  Shared TypeScript code for frontend hosts
- `go/core`
  Shared Go code for backend/host apps

## Top-Level Layout

```text
.
├── AGENTS.md
├── README.md
├── docs/
│   └── repo-structure.md
├── package.json
├── package-lock.json
├── tsconfig.base.json
├── go.work
├── apps/
├── packages/
└── go/
```

## Directory Guide

### `apps/desktop`

This is the installed desktop product.

Current responsibilities:

- Wails desktop shell
- embedded local SQLite
- update-check flow
- update download and install handoff
- React frontend under `frontend/`

Important files:

- [apps/desktop/main.go](/Users/zbango/Documents/ChatGPT/gym/apps/desktop/main.go)
  Wails app bootstrap
- [apps/desktop/app.go](/Users/zbango/Documents/ChatGPT/gym/apps/desktop/app.go)
  Wails-bound desktop methods
- [apps/desktop/wails.json](/Users/zbango/Documents/ChatGPT/gym/apps/desktop/wails.json)
  Wails project config
- [apps/desktop/internal/sqlite/store.go](/Users/zbango/Documents/ChatGPT/gym/apps/desktop/internal/sqlite/store.go)
  Local embedded SQLite adapter
- [apps/desktop/internal/updater/client.go](/Users/zbango/Documents/ChatGPT/gym/apps/desktop/internal/updater/client.go)
  Update manifest check logic

Rules:

- Desktop-specific OS/runtime logic stays here.
- Shared business logic should not accumulate here if it is also needed by cloud.
- SQLite adapters belong here or in another adapter layer, not in `go/core/domain`.

### `apps/desktop/frontend`

This is the React frontend mounted inside Wails.

Current responsibilities:

- desktop UI shell
- Wails JS bindings integration
- hello-world CRUD UI for the local SQLite slice
- update-check UI trigger
- update modal/download/install UI
- shared theme provider and switcher integration

Rules:

- Reuse `packages/ui` and `packages/shared`.
- Do not bury desktop-native assumptions inside shared packages.
- Treat this as a host app, not the place where product rules should live long-term.

### `apps/cloud-api`

This is the tiny cloud-facing Go layer.

Current responsibilities:

- Lambda entrypoint for low-cost deployment
- update manifest endpoint
- health/version endpoints

Important files:

- [apps/cloud-api/lambda/main.go](/Users/zbango/Documents/ChatGPT/gym/apps/cloud-api/lambda/main.go)

Rules:

- Keep this layer thin.
- It is not the main operational runtime of gym-saas.
- It should expose public endpoints and cloud adapters, not own all business logic itself.

### `apps/web`

This is the browser host using React + Vite.

Current responsibilities:

- shared frontend host setup
- reuse of shared UI and shared TS utilities

Rules:

- Reuse shared packages.
- Keep web-only hosting/auth/browser concerns local to `apps/web`.
- Do not make web constraints leak into desktop.

### `apps/mobile`

This is the Expo placeholder host.

Current responsibilities:

- keep the mobile slot real and bootable
- prove that shared frontend-safe code can be reused

Rules:

- Keep it intentionally small until mobile becomes an active product lane.
- Avoid building product-specific complexity here too early.

### `packages/ui`

Shared React UI components for multiple frontend hosts.

Current example:

- [packages/ui/src/index.tsx](/Users/zbango/Documents/ChatGPT/gym/packages/ui/src/index.tsx)
- [packages/ui/src/theme.tsx](/Users/zbango/Documents/ChatGPT/gym/packages/ui/src/theme.tsx)

Rules:

- Presentation only.
- No host-specific runtime logic.
- No backend/business rules.
- Theme provider and theme persistence logic for frontend hosts lives here.

### `packages/shared`

Shared frontend-safe TypeScript values and types.

Current contents include:

- app constants
- update manifest types
- hello record DTO
- theme definitions and tokens

Rules:

- Keep this package safe for web, desktop frontend, and mobile.
- Do not place Node-only or Wails-only logic here.
- Do not turn this into a dumping ground for random utilities.

### Shared theming

The current theme system is split like this:

- `packages/shared`
  theme names, labels, and token definitions
- `packages/ui`
  theme provider, CSS variable injection, persistence, and switcher
- frontend hosts
  consume CSS variables and themed UI primitives

Current behavior:

- desktop and web share the same theme engine
- the shell UI no longer depends on hardcoded colors for its main surfaces
- selected theme persists locally per host

### `go/core`

Shared Go module for code that should be reused by multiple Go hosts.

Current layout:

- `domain/`
- `application/`
- `ports/`
- `platform/`

Rules:

- Shared Go code goes here when it is needed by both desktop and cloud.
- Avoid putting importable shared code under a top-level `internal/`, because cross-module imports will break.
- Keep infrastructure adapters out of `domain/`.

## Ownership Boundaries

The most important architectural split is:

- TypeScript presents gym-saas
- Go runs gym-saas

That means:

- React hosts render UI, collect input, and call host/application boundaries
- Go hosts own local runtime work, cloud runtime work, SQLite integration, update checks, and future operational logic

This is not absolute yet because the repo is still at boilerplate stage, but that is the direction to preserve.

## Why `npm workspaces` and `go.work`

The repo is mixed-language, so one tool should not try to own everything.

- `npm workspaces` manages frontend hosts and shared TS packages
- `go.work` coordinates the separate Go modules

This keeps:

- frontend package management simple
- Wails conventional
- cloud and desktop Go code sharable without forcing one giant root Go module

## Local-First Runtime Model

The intended product model is:

- the gym PC is the main operational runtime
- desktop owns the local DB
- cloud stays thin and cheap
- cloud exists mainly for public entrypoints and future remote/product surfaces

That is why:

- local SQLite exists inside desktop
- the cloud API is intentionally small
- update metadata can live in the cloud while the real installed product stays local

## Embedded SQLite

The desktop app embeds SQLite through Go. The gym owner should not install SQL manually.

Current V2 database slice:

- embedded, versioned SQL migrations
- transactional migration records with checksum validation
- foreign keys, WAL journaling, and busy-timeout configuration
- operational core tables for gyms, members, plans, memberships, payments, and visits

Current files:

- [apps/desktop/internal/sqlite/store.go](/Users/zbango/Documents/ChatGPT/gym/apps/desktop/internal/sqlite/store.go)
- [apps/desktop/internal/sqlite/store_test.go](/Users/zbango/Documents/ChatGPT/gym/apps/desktop/internal/sqlite/store_test.go)

What this proves:

- the app can create its own local DB file
- local schema migrations are owned by the product and protected from drift
- the installer does not need an external database dependency

## Desktop Updates

The updater strategy is intentionally split into:

- shared update detection
- per-OS update apply behavior

The cross-platform part:

- desktop fetches a manifest
- compares versions
- filters assets for the current OS/arch

The OS-specific part:

- how the downloaded artifact is actually applied

Current files:

- [apps/desktop/internal/updater/client.go](/Users/zbango/Documents/ChatGPT/gym/apps/desktop/internal/updater/client.go)
- [apps/desktop/internal/updater/client_test.go](/Users/zbango/Documents/ChatGPT/gym/apps/desktop/internal/updater/client_test.go)

Current status:

- update detection exists
- desktop can download its update package itself
- mac-first install handoff exists
- fully polished cross-OS self-replacement does not exist yet

## Development Workflow

### Install

```bash
npm install
go work sync
```

### Shared cloud URL config

The web host and the Wails frontend read the cloud URL from the root env file:

```bash
VITE_CLOUD_API_BASE_URL=...
```

The root `.env` file holds stable cloud configuration such as `VITE_CLOUD_API_BASE_URL`, `R2_AWS_PROFILE`, `AWS_DEPLOY_PROFILE`, `R2_ACCOUNT_ID`, `R2_BUCKET_NAME`, and `R2_PUBLIC_BASE_URL`. Release-specific metadata like notes and generated artifact URLs are intended to be handled by the release tooling rather than stored manually in a second env file.

For updater testing in development, `.env` may also contain:

```bash
GYM_SAAS_DESKTOP_VERSION_OVERRIDE=...
```

When set, the desktop app reports that override as its runtime version during
`dev:desktop` and `dev:all`. This is useful for simulating an older installed
app while the deployed cloud manifest advertises a newer release.

The canonical app version lives in [VERSION](/Users/zbango/Documents/ChatGPT/gym/VERSION). The release tooling bumps that version automatically and syncs the desktop-visible version files before building and publishing a release.

### Frontend hosts

```bash
npm run dev:desktop
npm run dev:web
npm run dev:all
npm run dev:desktop-frontend
npm run dev:mobile
```

### Desktop development

From [apps/desktop](/Users/zbango/Documents/ChatGPT/gym/apps/desktop):

```bash
wails dev
```

If `wails` is not globally installed, the current verified fallback is:

```bash
go run github.com/wailsapp/wails/v2/cmd/wails@v2.10.1 build
```

### Recommended stack loops

```bash
npm run dev:desktop
```

This starts:

- the Wails desktop development host from `apps/desktop`

Note:

- if the updater install handoff is triggered in dev mode, the desktop process
  intentionally quits and the dev session ends

```bash
npm run dev:web
```

This starts:

- the web host

```bash
npm run dev:all
```

This starts:

- the Wails desktop development host
- the web host

### Verification commands

```bash
npm run typecheck
npm run build -w @gym-saas/web
npm run build -w @gym-saas/desktop-frontend
go test ./apps/cloud-api/... ./apps/desktop/... ./go/core/...
go build ./apps/cloud-api/lambda ./apps/desktop
```

### Cloud deployment

The cloud API is set up for AWS Lambda using AWS SAM.

Current deploy entrypoints:

```bash
npm run deploy
npm run deploy:guided
```

Relevant files:

- [apps/cloud-api/template.yaml](/Users/zbango/Documents/ChatGPT/gym/apps/cloud-api/template.yaml)
- [scripts/deploy-cloud-api.mjs](/Users/zbango/Documents/ChatGPT/gym/scripts/deploy-cloud-api.mjs)

The SAM stack exports a Lambda Function URL output. That URL is the remote cloud entrypoint you can feed into app configuration for update checks and future public endpoints.

## Current Verified State

As of August 14, 2026, the following have been verified in this repo:

- npm workspace install succeeds
- TS workspace typecheck succeeds
- web production build succeeds
- desktop frontend production build succeeds
- Go workspace sync succeeds
- Go build succeeds for:
  - desktop host
  - cloud Lambda entrypoint
- Go tests succeed for:
  - desktop SQLite slice
  - desktop updater slice
- Wails desktop packaging build succeeded on macOS and produced a desktop app bundle under
  [apps/desktop/build/bin](/Users/zbango/Documents/ChatGPT/gym/apps/desktop/build/bin)

## What Not To Do

- Do not reintroduce a root catch-all Go module for everything.
- Do not move shared Go code into a top-level `internal/` if multiple modules need it.
- Do not put business logic into React components because it feels faster.
- Do not let `packages/shared` become a random mixed-runtime junk drawer.
- Do not make the cloud API the only place gym-saas can function if the product is still local-first.

## Likely Next Steps

- replace the hello-world SQLite slice with the first real gym-saas local data slice
- define the first real shared Go application boundary in `go/core`
- evolve the cloud layer from update manifest only to the first real remote/public endpoint
- implement platform-specific installer/apply update flow after detection is stable
