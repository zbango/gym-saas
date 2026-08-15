# gym-saas

gym-saas is a mixed-language monorepo for a local-first gym product:

- `apps/desktop`: Wails desktop host with embedded SQLite and update checks
- `apps/cloud-api`: tiny Go cloud surface for update manifests and future public endpoints
- `apps/web`: React + Vite web host
- `apps/mobile`: Expo placeholder host
- `packages/ui`: shared React UI
- `packages/shared`: shared frontend-safe types and constants
- `go/core`: shared Go code for host apps

## Requirements

- Node.js 20+
- npm 10+
- Go 1.23+
- Wails CLI for desktop development

## Docs

- [Repo Structure](./docs/repo-structure.md)

## Install

```bash
npm install
go work sync
```

Root env file used by the web host and the Wails frontend:

```bash
VITE_CLOUD_API_BASE_URL=...
```

See [\.env.example](/Users/zbango/Documents/ChatGPT/gym/.env.example). The local workspace also has a [\.env](/Users/zbango/Documents/ChatGPT/gym/.env) pointing at the currently deployed Lambda Function URL.

The single env file also holds the stable R2 settings:

```bash
GYM_SAAS_DESKTOP_VERSION_OVERRIDE
R2_AWS_PROFILE
AWS_DEPLOY_PROFILE
R2_ACCOUNT_ID
R2_BUCKET_NAME
R2_PUBLIC_BASE_URL
```

Profile split:

- `R2_AWS_PROFILE`
  used for Cloudflare R2 uploads through the AWS CLI S3-compatible endpoint
- `AWS_DEPLOY_PROFILE`
  used for AWS SAM / Lambda deployment

Dev-only updater override:

- `GYM_SAAS_DESKTOP_VERSION_OVERRIDE`
  when set, the Wails desktop app reports that value as its runtime version
  during `dev:desktop` or `dev:all`
- use it only to simulate an older installed app while testing the updater

## Run

```bash
npm run dev:desktop
npm run dev:web
npm run dev:all
npm run dev:desktop-frontend
npm run dev:mobile
wails dev
```

Run `wails dev` from `apps/desktop`.

Recommended local loops:

```bash
npm run dev:desktop
```

What `dev:desktop` starts:

- the Wails desktop dev host

```bash
npm run dev:web
```

What `dev:web` starts:

- the web host

```bash
npm run dev:all
```

What `dev:all` starts:

- the Wails desktop dev host
- the web host

If you want to simulate an older installed desktop app while keeping the
current source tree, set for example:

```bash
GYM_SAAS_DESKTOP_VERSION_OVERRIDE=0.1.1
```

Then restart `dev:desktop` or `dev:all`.

## Desktop Packaging

Build a desktop package from `apps/desktop`:

```bash
wails build
```

For Windows NSIS installers:

```bash
wails build -nsis
```

## Repo Build

Build everything that is locally buildable in this repo:

```bash
npm run build:all
```

Target specific desktop packaging:

```bash
npm run build:desktop:mac
npm run build:desktop:linux
npm run build:desktop:windows
```

Target specific full builds:

```bash
npm run build:all:mac
npm run build:all:linux
npm run build:all:windows
```

What `build:all` does:

- runs repo typecheck
- builds the web app
- builds the cloud API Lambda binary into `apps/cloud-api/build`
- builds the desktop app package

Desktop packaging behavior:

- `build:all` builds for the current host target
- `build:all:windows` targets `windows/amd64` and adds `-nsis`
- `build:all:mac` targets `darwin/arm64`
- `build:all:linux` targets `linux/amd64`

Note:

- `build:all` does not produce an APK or IPA for mobile; `apps/mobile` is still a placeholder host
- cross-target desktop packaging still depends on the toolchain and platform support available on the machine running the build
- if you later need a universal macOS desktop bundle, add a separate `darwin/universal` build path after validating the local Apple toolchain

## Cloud Deploy

Deploy the cloud Lambda stack with AWS SAM:

```bash
npm run deploy
```

First-time guided deploy:

```bash
npm run deploy:guided
```

Notes:

- the web host and the Wails frontend read `VITE_CLOUD_API_BASE_URL` from the root env file at dev/build time
- this uses AWS SAM, not the `serverless` npm framework
- the SAM template lives at `apps/cloud-api/template.yaml`
- release metadata overrides are passed in at release/deploy time
- `deploy` automatically uses `AWS_DEPLOY_PROFILE`
- if `samconfig.toml` does not exist yet, `deploy` falls back to guided deployment
- the SAM template exports a Lambda Function URL output
- that output is the cloud URL you can use for app configuration such as update manifests or future public endpoints

## Update Publishing

The desktop app checks a remote version manifest exposed by `apps/cloud-api`.
Publishing a release means:

1. Build the desktop package for the target OS.
2. Upload the artifact to remote storage.
3. Update the cloud manifest payload with the new version, checksum, and asset URL.
4. Deploy the thin cloud layer.

Current R2-first helper for mac:

```bash
npm run build:desktop:mac
npm run publish:desktop:mac
npm run deploy
```

`publish:desktop:mac`:

- zips the built mac app bundle
- uploads it to Cloudflare R2 using the AWS CLI's S3-compatible endpoint
- computes a SHA-256 checksum
- returns the mac release URL/checksum to the release workflow
- automatically uses `R2_AWS_PROFILE`

Interactive release wizard:

```bash
npm run release
```

The wizard currently supports:

- choosing a release target
- auto-bumping the patch version from the repo `VERSION` file
- auto-setting the minimum supported version to the current released version
- running build, publish, and deploy steps

The wizard handles the selected OS target and keeps the target-specific artifact URL/checksum values in memory for the release flow. You should not need to edit those manually in the normal flow.

Versioning behavior:

- current version is read from [VERSION](/Users/zbango/Documents/ChatGPT/gym/VERSION)
- next release version is automatically computed as the next patch version
- current version becomes the minimum supported version for that release
- the release flow updates:
  - [VERSION](/Users/zbango/Documents/ChatGPT/gym/VERSION)
  - [packages/shared/src/index.ts](/Users/zbango/Documents/ChatGPT/gym/packages/shared/src/index.ts)
  - [go/core/platform/version.go](/Users/zbango/Documents/ChatGPT/gym/go/core/platform/version.go)
  - [apps/desktop/wails.json](/Users/zbango/Documents/ChatGPT/gym/apps/desktop/wails.json)

Today the full publish flow is implemented for:

- macOS Apple Silicon

## Updater Status

Current updater behavior on macOS:

- desktop app checks the deployed manifest automatically on startup
- desktop app can also re-check manually with `Check updates`
- if a newer version exists and a matching mac asset is present, the app:
  - downloads the package itself
  - hands off immediately to the install flow
  - attempts to relaunch the updated app

Important caveat for development:

- when testing from `dev:desktop` or `dev:all`, the updater intentionally quits
  the running desktop app during install handoff
- that ends the Wails dev session
- the most realistic updater test is against a built desktop app bundle, not the
  live dev process
