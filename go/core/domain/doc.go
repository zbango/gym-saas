// Package domain contains gym-saas business rules and value objects.
//
// Domain identifiers are UUID strings generated offline. Persisted timestamps
// use UTC RFC 3339 with nanosecond precision. Domain code must not import
// SQLite, HTTP, Wails, or device SDKs.
package domain
