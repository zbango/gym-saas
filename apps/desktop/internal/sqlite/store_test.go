package sqlite

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
)

func TestOpenAppliesV2SchemaOnce(t *testing.T) {
	path := filepath.Join(t.TempDir(), "gym-saas.db")

	store := openTestStore(t, path)
	assertSchemaVersion(t, store.db, 1)
	assertTableExists(t, store.db, "gyms")
	assertTableExists(t, store.db, "members")
	assertTableExists(t, store.db, "membership_plans")
	assertTableExists(t, store.db, "memberships")
	assertTableExists(t, store.db, "payments")
	assertTableExists(t, store.db, "visits")
	if err := store.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}

	store = openTestStore(t, path)
	defer store.Close()
	assertSchemaVersion(t, store.db, 1)
}

func TestOpenConfiguresSQLiteReliability(t *testing.T) {
	store := openTestStore(t, filepath.Join(t.TempDir(), "gym-saas.db"))
	defer store.Close()

	var foreignKeys, timeout int
	var journalMode string
	if err := store.db.QueryRow(`PRAGMA foreign_keys`).Scan(&foreignKeys); err != nil {
		t.Fatalf("read foreign_keys pragma: %v", err)
	}
	if err := store.db.QueryRow(`PRAGMA busy_timeout`).Scan(&timeout); err != nil {
		t.Fatalf("read busy_timeout pragma: %v", err)
	}
	if err := store.db.QueryRow(`PRAGMA journal_mode`).Scan(&journalMode); err != nil {
		t.Fatalf("read journal_mode pragma: %v", err)
	}

	if foreignKeys != 1 {
		t.Fatalf("foreign_keys = %d, want 1", foreignKeys)
	}
	if timeout != int(busyTimeout.Milliseconds()) {
		t.Fatalf("busy_timeout = %d, want %d", timeout, busyTimeout.Milliseconds())
	}
	if journalMode != "wal" {
		t.Fatalf("journal_mode = %q, want wal", journalMode)
	}

	_, err := store.db.Exec(`
INSERT INTO members (id, gym_id, first_name, last_name, phone, status, created_at, updated_at)
VALUES ('member-1', 'missing-gym', 'Ada', 'Lovelace', '555-0100', 'active', '2026-08-19T00:00:00Z', '2026-08-19T00:00:00Z')
`)
	if err == nil {
		t.Fatal("insert member without gym succeeded; foreign keys are not enforced")
	}
}

func TestApplyMigrationsRollsBackOnFailure(t *testing.T) {
	store := openTestStore(t, filepath.Join(t.TempDir(), "gym-saas.db"))
	defer store.Close()

	err := applyMigrations(context.Background(), store.db, append(migrations, migration{
		Version: 2,
		Name:    "broken",
		SQL: `
CREATE TABLE rollback_probe (id TEXT PRIMARY KEY);
THIS IS NOT VALID SQL;
`,
	}))
	if err == nil {
		t.Fatal("applyMigrations succeeded for invalid SQL")
	}

	assertSchemaVersion(t, store.db, 1)
	assertTableMissing(t, store.db, "rollback_probe")
}

func TestOpenRejectsUnmanagedDatabase(t *testing.T) {
	path := filepath.Join(t.TempDir(), "gym-saas.db")
	db, err := sql.Open(driverName, databaseURL(path))
	if err != nil {
		t.Fatalf("open unmanaged database: %v", err)
	}
	if _, err := db.Exec(`CREATE TABLE legacy_records (id TEXT PRIMARY KEY)`); err != nil {
		_ = db.Close()
		t.Fatalf("create unmanaged table: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close unmanaged database: %v", err)
	}

	if _, err := Open(path); err == nil {
		t.Fatal("Open succeeded for an unmanaged database")
	}
}

func TestOpenMigratesLegacyProofDatabase(t *testing.T) {
	path := filepath.Join(t.TempDir(), "gym-saas.db")
	db, err := sql.Open(driverName, databaseURL(path))
	if err != nil {
		t.Fatalf("open legacy proof database: %v", err)
	}
	if _, err := db.Exec(`CREATE TABLE hello_records (id TEXT PRIMARY KEY, name TEXT NOT NULL)`); err != nil {
		_ = db.Close()
		t.Fatalf("create legacy proof table: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close legacy proof database: %v", err)
	}

	store := openTestStore(t, path)
	defer store.Close()
	assertSchemaVersion(t, store.db, 1)
	assertTableExists(t, store.db, "hello_records")
	assertTableExists(t, store.db, "members")
}

func TestOpenRejectsMigrationDrift(t *testing.T) {
	path := filepath.Join(t.TempDir(), "gym-saas.db")
	store := openTestStore(t, path)
	if _, err := store.db.Exec(`UPDATE schema_migrations SET checksum = 'invalid' WHERE version = 1`); err != nil {
		_ = store.Close()
		t.Fatalf("corrupt migration checksum: %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("close database: %v", err)
	}

	if _, err := Open(path); err == nil {
		t.Fatal("Open succeeded after applied migration checksum drift")
	}
}

func openTestStore(t *testing.T, path string) *Store {
	t.Helper()
	store, err := Open(path)
	if err != nil {
		t.Fatalf("Open returned error: %v", err)
	}
	return store
}

func assertSchemaVersion(t *testing.T, db *sql.DB, want int) {
	t.Helper()
	var count, version int
	if err := db.QueryRow(`SELECT COUNT(*), COALESCE(MAX(version), 0) FROM schema_migrations`).Scan(&count, &version); err != nil {
		t.Fatalf("read schema migrations: %v", err)
	}
	if count != 1 || version != want {
		t.Fatalf("schema migrations = count %d, version %d; want count 1, version %d", count, version, want)
	}
}

func assertTableExists(t *testing.T, db *sql.DB, name string) {
	t.Helper()
	assertTable(t, db, name, true)
}

func assertTableMissing(t *testing.T, db *sql.DB, name string) {
	t.Helper()
	assertTable(t, db, name, false)
}

func assertTable(t *testing.T, db *sql.DB, name string, want bool) {
	t.Helper()
	var exists int
	if err := db.QueryRow(`SELECT EXISTS (SELECT 1 FROM sqlite_master WHERE type = 'table' AND name = ?)`, name).Scan(&exists); err != nil {
		t.Fatalf("check table %q: %v", name, err)
	}
	if (exists == 1) != want {
		t.Fatalf("table %q exists = %t, want %t", name, exists == 1, want)
	}
}
