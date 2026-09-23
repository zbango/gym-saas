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

func TestMembershipRepositoryPersistsVisitMembershipAcrossRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "gym-saas.db")
	store := openTestStore(t, path)
	membership, gymID := seedVisitMembership(t, store)
	repository := NewMembershipRepository(store)
	if err := repository.Create(context.Background(), membership); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}

	store = openTestStore(t, path)
	defer store.Close()
	found, err := NewMembershipRepository(store).Get(context.Background(), gymID, membership.ID())
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	assertMembershipEqual(t, found, membership)
}

func TestMembershipRepositoryUpdatesExhaustedVisitAllowance(t *testing.T) {
	store := openTestStore(t, filepath.Join(t.TempDir(), "gym-saas.db"))
	defer store.Close()
	membership, gymID := seedVisitMembership(t, store)
	repository := NewMembershipRepository(store)
	if err := repository.Create(context.Background(), membership); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	updated, err := membership.ConsumeVisit(membership.StartsAt().Add(time.Hour))
	if err != nil {
		t.Fatalf("first ConsumeVisit returned error: %v", err)
	}
	updated, err = updated.ConsumeVisit(membership.StartsAt().Add(2 * time.Hour))
	if err != nil {
		t.Fatalf("second ConsumeVisit returned error: %v", err)
	}
	if err := repository.Update(context.Background(), updated); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	found, err := repository.Get(context.Background(), gymID, membership.ID())
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if found.VisitsRemaining() != 0 || found.IsValidAt(updated.UpdatedAt()) {
		t.Fatalf("updated membership = %#v", found)
	}
	assertMembershipEqual(t, found, updated)
}

func TestMembershipRepositoryPersistsCancellation(t *testing.T) {
	store := openTestStore(t, filepath.Join(t.TempDir(), "gym-saas.db"))
	defer store.Close()
	membership, gymID := seedVisitMembership(t, store)
	cancelled, err := membership.Cancel(membership.StartsAt().Add(time.Hour), "member request")
	if err != nil {
		t.Fatalf("Cancel returned error: %v", err)
	}
	repository := NewMembershipRepository(store)
	if err := repository.Create(context.Background(), cancelled); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	found, err := repository.Get(context.Background(), gymID, cancelled.ID())
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	assertMembershipEqual(t, found, cancelled)
}

func TestMembershipRepositoryListsOnlyRequestedMemberAndReturnsNotFound(t *testing.T) {
	store := openTestStore(t, filepath.Join(t.TempDir(), "gym-saas.db"))
	defer store.Close()
	first, gymID := seedVisitMembership(t, store)
	secondMember, err := domain.CreateMember(gymID, "Grace", "Hopper", "grace@example.com", "+593 99 765 4321", "0203040506", "1906-12-09", "", domain.MemberStatusActive, first.StartsAt())
	if err != nil {
		t.Fatalf("CreateMember returned error: %v", err)
	}
	if err := NewMemberRepository(store).Create(context.Background(), secondMember); err != nil {
		t.Fatalf("create second member: %v", err)
	}
	plan := mustMembershipPlan(t, gymID, "Monthly", domain.MembershipValidityTime, 0)
	if err := NewMembershipPlanRepository(store).Create(context.Background(), plan); err != nil {
		t.Fatalf("create time plan: %v", err)
	}
	gym := mustGymForRepository(t, gymID)
	second, err := domain.StartMembership(gym, secondMember.ID(), plan, first.StartsAt(), first.StartsAt())
	if err != nil {
		t.Fatalf("StartMembership returned error: %v", err)
	}
	repository := NewMembershipRepository(store)
	for _, membership := range []domain.Membership{first, second} {
		if err := repository.Create(context.Background(), membership); err != nil {
			t.Fatalf("Create returned error: %v", err)
		}
	}

	memberships, err := repository.ListForMember(context.Background(), gymID, first.MemberID())
	if err != nil {
		t.Fatalf("ListForMember returned error: %v", err)
	}
	if len(memberships) != 1 {
		t.Fatalf("ListForMember length = %d, want 1", len(memberships))
	}
	assertMembershipEqual(t, memberships[0], first)

	missingID, err := domain.NewMembershipID()
	if err != nil {
		t.Fatalf("NewMembershipID returned error: %v", err)
	}
	if _, err := repository.Get(context.Background(), gymID, missingID); !errors.Is(err, ports.ErrMembershipNotFound) {
		t.Fatalf("Get missing error = %v, want %v", err, ports.ErrMembershipNotFound)
	}
}

func seedVisitMembership(t *testing.T, store *Store) (domain.Membership, domain.GymID) {
	t.Helper()
	gymID := seedGym(t, store)
	member := mustMember(t, gymID, "Ada", "Lovelace")
	if err := NewMemberRepository(store).Create(context.Background(), member); err != nil {
		t.Fatalf("create member: %v", err)
	}
	plan := mustMembershipPlan(t, gymID, "10 visits", domain.MembershipValidityVisits, 2)
	if err := NewMembershipPlanRepository(store).Create(context.Background(), plan); err != nil {
		t.Fatalf("create plan: %v", err)
	}
	startsAt := time.Date(2026, time.August, 20, 10, 0, 0, 0, time.UTC)
	membership, err := domain.StartMembership(mustGymForRepository(t, gymID), member.ID(), plan, startsAt, startsAt)
	if err != nil {
		t.Fatalf("StartMembership returned error: %v", err)
	}
	return membership, gymID
}

func mustGymForRepository(t *testing.T, gymID domain.GymID) domain.Gym {
	t.Helper()
	now := time.Date(2026, time.August, 19, 10, 0, 0, 0, time.UTC)
	gym, err := domain.NewGym(gymID, "Zeus", "UTC", now, now)
	if err != nil {
		t.Fatalf("NewGym returned error: %v", err)
	}
	return gym
}

func assertMembershipEqual(t *testing.T, got, want domain.Membership) {
	t.Helper()
	if got.ID() != want.ID() || got.GymID() != want.GymID() || got.MemberID() != want.MemberID() || got.MembershipPlanID() != want.MembershipPlanID() || got.Status() != want.Status() || !got.StartsAt().Equal(want.StartsAt()) || !got.EndsAt().Equal(want.EndsAt()) || !got.ActivatedAt().Equal(want.ActivatedAt()) || !got.CancelledAt().Equal(want.CancelledAt()) || got.CancellationReason() != want.CancellationReason() || got.ValidityKind() != want.ValidityKind() || got.DurationValue() != want.DurationValue() || got.DurationUnit() != want.DurationUnit() || got.VisitLimit() != want.VisitLimit() || got.VisitsRemaining() != want.VisitsRemaining() || got.Price().Cents() != want.Price().Cents() || got.Price().Currency() != want.Price().Currency() || !got.CreatedAt().Equal(want.CreatedAt()) || !got.UpdatedAt().Equal(want.UpdatedAt()) {
		t.Fatalf("membership = %#v, want %#v", got, want)
	}
}
