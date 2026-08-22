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

func TestMembershipPlanRepositoryPersistsVisitPlanAcrossRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "gym-saas.db")
	store := openTestStore(t, path)
	gymID := seedGym(t, store)
	repository := NewMembershipPlanRepository(store)
	plan := mustMembershipPlan(t, gymID, "10 visits / 2 months", domain.MembershipValidityVisits, 10)
	if err := repository.Create(context.Background(), plan); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}

	store = openTestStore(t, path)
	defer store.Close()
	found, err := NewMembershipPlanRepository(store).Get(context.Background(), gymID, plan.ID())
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	assertMembershipPlanEqual(t, found, plan)
}

func TestMembershipPlanRepositoryUpdatesAndArchivesPlan(t *testing.T) {
	store := openTestStore(t, filepath.Join(t.TempDir(), "gym-saas.db"))
	defer store.Close()
	gymID := seedGym(t, store)
	repository := NewMembershipPlanRepository(store)
	plan := mustMembershipPlan(t, gymID, "Monthly", domain.MembershipValidityTime, 0)
	if err := repository.Create(context.Background(), plan); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	price, err := domain.NewMoney(5500, "USD")
	if err != nil {
		t.Fatalf("NewMoney returned error: %v", err)
	}
	updated, err := domain.NewMembershipPlan(plan.ID(), gymID, "Quarterly", domain.MembershipValidityTime, 3, domain.MembershipDurationMonths, 0, price, domain.MembershipPlanStatusInactive, plan.CreatedAt(), plan.UpdatedAt().Add(time.Hour))
	if err != nil {
		t.Fatalf("NewMembershipPlan returned error: %v", err)
	}
	if err := repository.Update(context.Background(), updated); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	found, err := repository.Get(context.Background(), gymID, plan.ID())
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	assertMembershipPlanEqual(t, found, updated)

	if err := repository.Archive(context.Background(), gymID, plan.ID(), updated.UpdatedAt().Add(time.Hour)); err != nil {
		t.Fatalf("Archive returned error: %v", err)
	}
	if _, err := repository.Get(context.Background(), gymID, plan.ID()); !errors.Is(err, ports.ErrMembershipPlanNotFound) {
		t.Fatalf("Get after archive error = %v, want %v", err, ports.ErrMembershipPlanNotFound)
	}
}

func mustMembershipPlan(t *testing.T, gymID domain.GymID, name string, kind domain.MembershipValidityKind, visitLimit int) domain.MembershipPlan {
	t.Helper()
	price, err := domain.NewMoney(4500, "USD")
	if err != nil {
		t.Fatalf("NewMoney returned error: %v", err)
	}
	plan, err := domain.CreateMembershipPlan(gymID, name, kind, 2, domain.MembershipDurationMonths, visitLimit, price, domain.MembershipPlanStatusActive, time.Date(2026, time.August, 20, 10, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("CreateMembershipPlan returned error: %v", err)
	}
	return plan
}

func assertMembershipPlanEqual(t *testing.T, got, want domain.MembershipPlan) {
	t.Helper()
	if got.ID() != want.ID() || got.GymID() != want.GymID() || got.Name() != want.Name() || got.ValidityKind() != want.ValidityKind() || got.DurationValue() != want.DurationValue() || got.DurationUnit() != want.DurationUnit() || got.VisitLimit() != want.VisitLimit() || got.Price().Cents() != want.Price().Cents() || got.Price().Currency() != want.Price().Currency() || got.Status() != want.Status() || !got.CreatedAt().Equal(want.CreatedAt()) || !got.UpdatedAt().Equal(want.UpdatedAt()) {
		t.Fatalf("plan = %#v, want %#v", got, want)
	}
}
