package sqlite

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/zbango/gym-saas/go/core/domain"
	"github.com/zbango/gym-saas/go/core/ports"
)

func TestMemberRepositoryPersistsAcrossRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "gym-saas.db")
	store := openTestStore(t, path)
	gymID := seedGym(t, store)
	repository := NewMemberRepository(store)
	member := mustMember(t, gymID, "Ada", "Lovelace")
	if err := repository.Create(context.Background(), member); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}

	store = openTestStore(t, path)
	defer store.Close()
	found, err := NewMemberRepository(store).Get(context.Background(), gymID, member.ID())
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	assertMemberEqual(t, found, member)
}

func TestMemberRepositoryListsOnlyRequestedGym(t *testing.T) {
	store := openTestStore(t, filepath.Join(t.TempDir(), "gym-saas.db"))
	defer store.Close()
	firstGymID := seedGym(t, store)
	secondGymID := seedGym(t, store)
	repository := NewMemberRepository(store)
	first := mustMember(t, firstGymID, "Ada", "Lovelace")
	second := mustMember(t, secondGymID, "Grace", "Hopper")
	for _, member := range []domain.Member{first, second} {
		if err := repository.Create(context.Background(), member); err != nil {
			t.Fatalf("Create returned error: %v", err)
		}
	}

	members, err := repository.List(context.Background(), firstGymID)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(members) != 1 {
		t.Fatalf("List length = %d, want 1", len(members))
	}
	assertMemberEqual(t, members[0], first)
}

func TestMemberRepositoryReturnsNotFound(t *testing.T) {
	store := openTestStore(t, filepath.Join(t.TempDir(), "gym-saas.db"))
	defer store.Close()
	gymID := seedGym(t, store)
	memberID, err := domain.NewMemberID()
	if err != nil {
		t.Fatalf("NewMemberID returned error: %v", err)
	}

	_, err = NewMemberRepository(store).Get(context.Background(), gymID, memberID)
	if !errors.Is(err, ports.ErrMemberNotFound) {
		t.Fatalf("Get error = %v, want %v", err, ports.ErrMemberNotFound)
	}
}

func TestMemberRepositoryUpdatesAndArchivesMember(t *testing.T) {
	store := openTestStore(t, filepath.Join(t.TempDir(), "gym-saas.db"))
	defer store.Close()
	gymID := seedGym(t, store)
	repository := NewMemberRepository(store)
	member := mustMember(t, gymID, "Ada", "Lovelace")
	if err := repository.Create(context.Background(), member); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	updatedAt := member.UpdatedAt().Add(time.Hour)
	updated, err := domain.NewMember(member.ID(), gymID, "Augusta", "King", "augusta@example.com", "555-0100", "123", "1815-12-10", "New address", domain.MemberStatusInactive, member.CreatedAt(), updatedAt)
	if err != nil {
		t.Fatalf("NewMember returned error: %v", err)
	}
	if err := repository.Update(context.Background(), updated); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	found, err := repository.Get(context.Background(), gymID, member.ID())
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	assertMemberEqual(t, found, updated)

	if err := repository.Archive(context.Background(), gymID, member.ID(), updatedAt.Add(time.Hour)); err != nil {
		t.Fatalf("Archive returned error: %v", err)
	}
	if _, err := repository.Get(context.Background(), gymID, member.ID()); !errors.Is(err, ports.ErrMemberNotFound) {
		t.Fatalf("Get after archive error = %v, want %v", err, ports.ErrMemberNotFound)
	}
	members, err := repository.List(context.Background(), gymID)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(members) != 0 {
		t.Fatalf("List length = %d, want 0 after archive", len(members))
	}
}

func seedGym(t *testing.T, store *Store) domain.GymID {
	t.Helper()
	gymID, err := domain.NewGymID()
	if err != nil {
		t.Fatalf("NewGymID returned error: %v", err)
	}
	now := domain.FormatTimestamp(time.Date(2026, time.August, 19, 10, 0, 0, 0, time.UTC))
	if _, err := store.db.Exec(`INSERT INTO gyms (id, name, timezone, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`, gymID, "Zeus", "UTC", now, now); err != nil {
		t.Fatalf("seed gym: %v", err)
	}
	return gymID
}

func mustMember(t *testing.T, gymID domain.GymID, firstName, lastName string) domain.Member {
	t.Helper()
	member, err := domain.CreateMember(gymID, firstName, lastName, firstName+"@example.com", "+593 99 123 4567", "0102030405", "1815-12-10", "12 St. James's Square", domain.MemberStatusActive, time.Date(2026, time.August, 19, 10, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("CreateMember returned error: %v", err)
	}
	return member
}

func assertMemberEqual(t *testing.T, got, want domain.Member) {
	t.Helper()
	if got.ID() != want.ID() || got.GymID() != want.GymID() || got.FirstName() != want.FirstName() || got.LastName() != want.LastName() || got.Email() != want.Email() || got.Phone() != want.Phone() || got.IdentificationNumber() != want.IdentificationNumber() || got.DateOfBirth() != want.DateOfBirth() || got.Address() != want.Address() || got.Status() != want.Status() || !got.CreatedAt().Equal(want.CreatedAt()) || !got.UpdatedAt().Equal(want.UpdatedAt()) {
		t.Fatalf("member = %#v, want %#v", got, want)
	}
}
