package sqlite

import (
	"path/filepath"
	"testing"
)

func TestStoreCRUD(t *testing.T) {
	store, err := Open(filepath.Join(t.TempDir(), "gym-saas.db"))
	if err != nil {
		t.Fatalf("Open returned error: %v", err)
	}
	defer store.Close()

	created, err := store.Create("first")
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	listed, err := store.List()
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(listed) != 1 {
		t.Fatalf("List length = %d, want 1", len(listed))
	}

	updated, err := store.Update(created.ID, "second")
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if updated.Name != "second" {
		t.Fatalf("Update name = %q, want %q", updated.Name, "second")
	}

	if err := store.Delete(created.ID); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}

	listed, err = store.List()
	if err != nil {
		t.Fatalf("List after delete returned error: %v", err)
	}
	if len(listed) != 0 {
		t.Fatalf("List after delete length = %d, want 0", len(listed))
	}
}
