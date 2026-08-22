package main

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	dbsqlite "github.com/zbango/gym-saas/apps/desktop/internal/sqlite"
	"github.com/zbango/gym-saas/go/core/application"
	"github.com/zbango/gym-saas/go/core/domain"
)

func TestAppMembershipPlanMethodsUseThePlanService(t *testing.T) {
	store, err := dbsqlite.Open(filepath.Join(t.TempDir(), "gym-saas.db"))
	if err != nil {
		t.Fatalf("Open returned error: %v", err)
	}
	defer store.Close()
	gymID, err := domain.NewGymID()
	if err != nil {
		t.Fatalf("NewGymID returned error: %v", err)
	}
	now := time.Date(2026, time.August, 20, 10, 0, 0, 0, time.UTC)
	gym, err := domain.NewGym(gymID, "Test gym", "UTC", now, now)
	if err != nil {
		t.Fatalf("NewGym returned error: %v", err)
	}
	if err := store.EnsureGym(context.Background(), gym); err != nil {
		t.Fatalf("EnsureGym returned error: %v", err)
	}
	service, err := application.NewMembershipPlanService(dbsqlite.NewMembershipPlanRepository(store), gymID, func() time.Time { return now })
	if err != nil {
		t.Fatalf("NewMembershipPlanService returned error: %v", err)
	}
	api := NewMembershipPlanAPI(&DesktopRuntime{}, service)
	created, err := api.CreateMembershipPlan(MembershipPlanInput{
		Name: "10 visits / 2 months", ValidityKind: "visits", DurationValue: 2, DurationUnit: "months",
		VisitLimit: 10, PriceCents: 5000, Currency: "USD", Status: "active",
	})
	if err != nil {
		t.Fatalf("CreateMembershipPlan returned error: %v", err)
	}
	if created.ID == "" || created.VisitLimit != 10 || created.DurationValue != 2 {
		t.Fatalf("CreateMembershipPlan = %#v", created)
	}
	updated, err := api.UpdateMembershipPlan(created.ID, MembershipPlanInput{
		Name: "20 visits / 2 months", ValidityKind: "visits", DurationValue: 2, DurationUnit: "months",
		VisitLimit: 20, PriceCents: 7000, Currency: "USD", Status: "inactive",
	})
	if err != nil {
		t.Fatalf("UpdateMembershipPlan returned error: %v", err)
	}
	if updated.Name != "20 visits / 2 months" || updated.VisitLimit != 20 || updated.Status != "inactive" {
		t.Fatalf("UpdateMembershipPlan = %#v", updated)
	}
	if err := api.ArchiveMembershipPlan(created.ID); err != nil {
		t.Fatalf("ArchiveMembershipPlan returned error: %v", err)
	}
	plans, err := api.ListMembershipPlans()
	if err != nil {
		t.Fatalf("ListMembershipPlans returned error: %v", err)
	}
	if len(plans) != 0 {
		t.Fatalf("ListMembershipPlans length = %d, want 0 after archive", len(plans))
	}
}
