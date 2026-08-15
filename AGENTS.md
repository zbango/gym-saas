# gym-saas Engineering Rules

- Go owns all business logic.
- React must not implement business rules.
- Domain packages cannot import SQLite, HTTP, Wails, or device SDKs.
- Infrastructure depends on domain, never the opposite.
- Money is stored as integer cents.
- IDs are strings for now and should converge on UUIDs.
- Every domain rule requires unit tests.
- Every local mutation that requires cloud synchronization must write an outbox event in the same SQLite transaction.
- Device operations must be asynchronous and persisted.
- Do not introduce new dependencies without justification.
- Do not modify architecture without explicit approval.
- Prefer simple Go over frameworks.
