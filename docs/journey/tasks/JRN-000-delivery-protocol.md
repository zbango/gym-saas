# JRN-000 — Human-in-the-loop delivery protocol

**Status:** Accepted
**Roadmap:** Supporting delivery task
**Owner:** Product owner + Codex
**Opened:** 2026-08-19

## Objective

Establish one durable record for each V2 task so implementation cannot move
forward without explicit, reproducible human acceptance.

## Scope

### In scope

- A task template containing plan, testable acceptance criteria, evidence,
  risks, and owner sign-off.
- A status gate that stops the next task until the owner accepts the current
  task.
- Historical entries for the V2 work completed before this protocol.

### Out of scope

- Automating product-owner approval.
- Rewriting the existing roadmap.
- Retroactively changing completed code.

## Plan

1. Add a central journey index and repeatable task template.
2. Document statuses and the explicit human acceptance gate.
3. Record the completed migration and domain-foundation tasks with repeatable
   verification evidence.
4. Ask the owner to verify the protocol before opening the next implementation
   task.

## Acceptance criteria

- [x] Every future task has a place to record scope, plan, acceptance
  criteria, risks, changed files, verification, and owner acceptance.
- [x] The workflow explicitly forbids opening the next implementation task
  until the current one is accepted.
- [x] Completed pre-protocol V2 tasks are traceable with their verification
  evidence and manual check commands.
- [x] The owner confirmed the workflow and template are suitable for use.

## Risks and decisions

| Risk or decision | Mitigation or outcome |
|---|---|
| Journey files become generic status reports. | Each page requires testable criteria, exact commands, and an expected manual result. |
| A pure-domain task has no visual UI to inspect. | Its manual check is a named command with expected test output and an owner review of the task page. |
| Historical work could incorrectly block new work. | It is marked historical; the acceptance gate begins with JRN-000. |

## Implementation record

### Changed files

- `docs/journey/README.md`
- `docs/journey/TEMPLATE.md`
- `docs/journey/tasks/JRN-000-delivery-protocol.md`
- Historical task pages under `docs/journey/tasks/`
- `docs/README.md`

### Simplifications made

- Markdown files are the source of truth; no workflow service or dependency
  was introduced.

### Automated verification

```text
command: git diff --check
result: passed
```

### Known limitations / follow-ups

- The first post-protocol product task must be planned in its own page only
  after the owner accepts JRN-000.

## Human acceptance

### Manual check

1. Open [the journey index](../README.md) and this task page.
2. Confirm the template captures the information you need to approve or send a
   task back: objective, scope, plan, acceptance criteria, evidence, risks,
   and owner result.
3. Confirm that **Awaiting human acceptance** is a clear stop point for the
   next implementation task.

**Expected result:** The process is understandable without referring to chat
history, and you can approve or request changes directly on the task page.

### Owner result

- [x] Accepted
- [ ] Changes requested

**Date:** 2026-08-19
**Notes:** Owner approved the protocol in chat.
