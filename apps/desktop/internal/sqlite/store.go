package sqlite

import (
	"context"
	"crypto/sha256"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/zbango/gym-saas/go/core/domain"
	_ "modernc.org/sqlite"
)

const (
	driverName  = "sqlite"
	busyTimeout = 5 * time.Second
)

type migration struct {
	Version int
	Name    string
	SQL     string
}

type appliedMigration struct {
	Name     string
	Checksum string
}

//go:embed migrations/001_v2_operational_core.sql
var operationalCoreSQL string

//go:embed migrations/002_visit_plan_expiry.sql
var visitPlanExpirySQL string

var migrations = []migration{
	{Version: 1, Name: "v2_operational_core", SQL: operationalCoreSQL},
	{Version: 2, Name: "visit_plan_expiry", SQL: visitPlanExpirySQL},
}

type Store struct {
	db *sql.DB
}

func Open(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create db directory: %w", err)
	}

	db, err := sql.Open(driverName, databaseURL(path))
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	db.SetMaxOpenConns(1)

	if err := db.PingContext(context.Background()); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("connect sqlite: %w", err)
	}
	if err := applyMigrations(context.Background(), db, migrations); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

// EnsureGym creates the desktop's local tenant on its first launch. Existing
// records are left untouched so member data always remains associated with the
// same tenant across restarts.
func (s *Store) EnsureGym(ctx context.Context, gym domain.Gym) error {
	_, err := s.db.ExecContext(ctx, `
INSERT INTO gyms (id, name, timezone, created_at, updated_at)
VALUES (?, ?, ?, ?, ?)
ON CONFLICT (id) DO NOTHING
`, gym.ID(), gym.Name(), gym.Timezone(), domain.FormatTimestamp(gym.CreatedAt()), domain.FormatTimestamp(gym.UpdatedAt()))
	if err != nil {
		return fmt.Errorf("ensure gym: %w", err)
	}
	return nil
}

func databaseURL(path string) string {
	u := url.URL{Scheme: "file", Path: path}
	query := u.Query()
	query.Add("_pragma", fmt.Sprintf("busy_timeout(%d)", busyTimeout.Milliseconds()))
	query.Add("_pragma", "foreign_keys(1)")
	query.Add("_pragma", "journal_mode(WAL)")
	query.Add("_pragma", "synchronous(NORMAL)")
	u.RawQuery = query.Encode()
	return u.String()
}

func applyMigrations(ctx context.Context, db *sql.DB, available []migration) error {
	if err := validateMigrations(available); err != nil {
		return err
	}

	applied, err := loadAppliedMigrations(ctx, db)
	if err != nil {
		return err
	}

	knownVersions := make(map[int]struct{}, len(available))
	for _, current := range available {
		knownVersions[current.Version] = struct{}{}
		if appliedMigration, ok := applied[current.Version]; ok {
			if appliedMigration.Name != current.Name {
				return fmt.Errorf("migration %d name mismatch: database has %q, binary has %q", current.Version, appliedMigration.Name, current.Name)
			}
			if appliedMigration.Checksum != migrationChecksum(current.SQL) {
				return fmt.Errorf("migration %d checksum mismatch; applied migrations must not be edited", current.Version)
			}
			continue
		}

		if err := applyMigration(ctx, db, current); err != nil {
			return fmt.Errorf("apply migration %d_%s: %w", current.Version, current.Name, err)
		}
	}
	for version := range applied {
		if _, ok := knownVersions[version]; !ok {
			return fmt.Errorf("database requires unknown migration %d; refusing to run an older binary", version)
		}
	}

	return nil
}

func validateMigrations(available []migration) error {
	if len(available) == 0 {
		return errors.New("no migrations configured")
	}

	for index, current := range available {
		if current.Version <= 0 || current.Name == "" || current.SQL == "" {
			return fmt.Errorf("invalid migration at index %d", index)
		}
		if index > 0 && available[index-1].Version >= current.Version {
			return fmt.Errorf("migration versions must be strictly increasing: %d then %d", available[index-1].Version, current.Version)
		}
	}

	return nil
}

func loadAppliedMigrations(ctx context.Context, db *sql.DB) (map[int]appliedMigration, error) {
	var schemaMigrationsExists int
	err := db.QueryRowContext(ctx, `
SELECT EXISTS (
  SELECT 1
  FROM sqlite_master
  WHERE type = 'table' AND name = 'schema_migrations'
)
`).Scan(&schemaMigrationsExists)
	if err != nil {
		return nil, fmt.Errorf("check schema migrations table: %w", err)
	}

	if schemaMigrationsExists == 0 {
		legacyProofDatabase, err := isLegacyProofDatabase(ctx, db)
		if err != nil {
			return nil, err
		}
		if !legacyProofDatabase {
			return nil, errors.New("database has application tables but no schema_migrations; refusing to overwrite an unmanaged database")
		}
		return map[int]appliedMigration{}, nil
	}

	rows, err := db.QueryContext(ctx, `
SELECT version, name, checksum
FROM schema_migrations
ORDER BY version
`)
	if err != nil {
		return nil, fmt.Errorf("load applied migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[int]appliedMigration)
	for rows.Next() {
		var version int
		var migration appliedMigration
		if err := rows.Scan(&version, &migration.Name, &migration.Checksum); err != nil {
			return nil, fmt.Errorf("scan applied migration: %w", err)
		}
		applied[version] = migration
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate applied migrations: %w", err)
	}

	return applied, nil
}

func isLegacyProofDatabase(ctx context.Context, db *sql.DB) (bool, error) {
	rows, err := db.QueryContext(ctx, `
SELECT name
FROM sqlite_master
WHERE type = 'table' AND name NOT LIKE 'sqlite_%'
ORDER BY name
`)
	if err != nil {
		return false, fmt.Errorf("list unmanaged tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return false, fmt.Errorf("scan unmanaged table: %w", err)
		}
		tables = append(tables, name)
	}
	if err := rows.Err(); err != nil {
		return false, fmt.Errorf("iterate unmanaged tables: %w", err)
	}

	return len(tables) == 0 || (len(tables) == 1 && tables[0] == "hello_records"), nil
}

func applyMigration(ctx context.Context, db *sql.DB, current migration) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, current.SQL); err != nil {
		return fmt.Errorf("execute SQL: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `
INSERT INTO schema_migrations (version, name, checksum, applied_at)
VALUES (?, ?, ?, ?)
`, current.Version, current.Name, migrationChecksum(current.SQL), time.Now().UTC().Format(time.RFC3339Nano)); err != nil {
		return fmt.Errorf("record migration: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func migrationChecksum(sql string) string {
	sum := sha256.Sum256([]byte(sql))
	return fmt.Sprintf("%x", sum)
}
