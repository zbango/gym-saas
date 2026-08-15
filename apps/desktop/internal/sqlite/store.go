package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

type HelloRecord struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type Store struct {
	db *sql.DB
}

func Open(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create db directory: %w", err)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	store := &Store{db: db}
	if err := store.bootstrap(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return store, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) bootstrap() error {
	const statement = `
CREATE TABLE IF NOT EXISTS hello_records (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);`

	if _, err := s.db.Exec(statement); err != nil {
		return fmt.Errorf("bootstrap hello_records: %w", err)
	}

	return nil
}

func (s *Store) List() ([]HelloRecord, error) {
	rows, err := s.db.Query(`
SELECT id, name, created_at, updated_at
FROM hello_records
ORDER BY created_at DESC
`)
	if err != nil {
		return nil, fmt.Errorf("list hello records: %w", err)
	}
	defer rows.Close()

	records := []HelloRecord{}
	for rows.Next() {
		var record HelloRecord
		if err := rows.Scan(&record.ID, &record.Name, &record.CreatedAt, &record.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan hello record: %w", err)
		}
		records = append(records, record)
	}

	return records, rows.Err()
}

func (s *Store) Create(name string) (HelloRecord, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	record := HelloRecord{
		ID:        uuid.NewString(),
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err := s.db.Exec(`
INSERT INTO hello_records (id, name, created_at, updated_at)
VALUES (?, ?, ?, ?)
`, record.ID, record.Name, record.CreatedAt, record.UpdatedAt)
	if err != nil {
		return HelloRecord{}, fmt.Errorf("create hello record: %w", err)
	}

	return record, nil
}

func (s *Store) Update(id, name string) (HelloRecord, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := s.db.Exec(`
UPDATE hello_records
SET name = ?, updated_at = ?
WHERE id = ?
`, name, now, id)
	if err != nil {
		return HelloRecord{}, fmt.Errorf("update hello record: %w", err)
	}

	var record HelloRecord
	err = s.db.QueryRow(`
SELECT id, name, created_at, updated_at
FROM hello_records
WHERE id = ?
`, id).Scan(&record.ID, &record.Name, &record.CreatedAt, &record.UpdatedAt)
	if err != nil {
		return HelloRecord{}, fmt.Errorf("reload hello record: %w", err)
	}

	return record, nil
}

func (s *Store) Delete(id string) error {
	_, err := s.db.Exec(`DELETE FROM hello_records WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete hello record: %w", err)
	}

	return nil
}
