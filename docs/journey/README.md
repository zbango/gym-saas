# Delivery Journey

This directory is the human-in-the-loop record for every V2 task. It makes a
task's intended behavior, implementation evidence, and manual acceptance
check visible before work moves to the next task.

## Workflow

1. Create a task page from [TEMPLATE.md](./TEMPLATE.md) before implementation.
2. Write a small, feasible plan and concrete acceptance criteria.
3. Implement only that task's approved scope.
4. Record changed files, automated verification, and known limitations.
5. Set the task status to **Awaiting human acceptance**.
6. The owner runs the documented manual check and records **Accepted** or
   **Changes requested** on the task page.
7. Open the next implementation task only after the previous task is
   accepted.

Each task must remain a small, reviewable change set. A failed manual check
returns the same task to **Planned** with the observed result recorded; it
does not silently advance the roadmap.

## Status meanings

| Status | Meaning |
|---|---|
| Proposed | Scope is being clarified; no implementation has started. |
| Planned | Scope, acceptance criteria, risks, and tests are ready for implementation. |
| Implementing | Code or documentation is actively changing. |
| Awaiting human acceptance | Automated checks passed; the owner must perform the stated manual check. |
| Accepted | The owner accepted the result; the next task may be opened. |
| Changes requested | The owner found a gap; record it and return the task to Planned. |
| Blocked | A required decision or external dependency prevents safe progress. |

## Task index

| Task | Status | Purpose |
|---|---|---|
| [JRN-000](./tasks/JRN-000-delivery-protocol.md) | Accepted | Establish this delivery and acceptance process. |
| [ZV2-022](./tasks/ZV2-022-gym-branch-identity.md) | Awaiting human acceptance | Gym and Branch domain identity model. |
| [ZV2-043](./tasks/ZV2-043-sqlite-migration-runner.md) | Historical—reviewable | Embedded, transactional SQLite migration runner. |
| [ZV2-020-021](./tasks/ZV2-020-021-domain-conventions-money.md) | Historical—reviewable | UUID/UTC conventions and integer-cent Money. |

Tasks completed before this protocol are recorded as historical: their
automated evidence and a repeatable manual check are preserved, but they do
not retroactively block JRN-000. Every task opened after JRN-000 must use the
status gate above.
