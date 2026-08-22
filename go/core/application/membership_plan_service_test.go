package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/zbango/gym-saas/go/core/domain"
	"github.com/zbango/gym-saas/go/core/ports"
)

func TestMembershipPlanServiceCreatesUpdatesAndArchivesPlan(t *testing.T) {
	repository := &memoryMembershipPlanRepository{plans: map[domain.MembershipPlanID]domain.MembershipPlan{}}
	times := []time.Time{
		time.Date(2026, time.August, 20, 10, 0, 0, 0, time.UTC),
		time.Date(2026, time.August, 20, 11, 0, 0, 0, time.UTC),
		time.Date(2026, time.August, 20, 12, 0, 0, 0, time.UTC),
	}
	service, err := NewMembershipPlanService(repository, mustGymID(t), func() time.Time {
		value := times[0]
		times = times[1:]
		return value
	})
	if err != nil {
		t.Fatalf("NewMembershipPlanService returned error: %v", err)
	}
	created, err := service.Create(context.Background(), planInput("10 visits / 2 months", "visits", 10))
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	updatedInput := planInput("20 visits / 2 months", "visits", 20)
	updated, err := service.Update(context.Background(), string(created.ID()), updatedInput)
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if updated.Name() != updatedInput.Name || updated.VisitLimit() != 20 || !updated.CreatedAt().Equal(created.CreatedAt()) {
		t.Fatalf("Update = %#v", updated)
	}
	if err := service.Archive(context.Background(), string(created.ID())); err != nil {
		t.Fatalf("Archive returned error: %v", err)
	}
	plans, err := service.List(context.Background())
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(plans) != 0 {
		t.Fatalf("List length = %d, want 0 after archive", len(plans))
	}
}

func TestMembershipPlanServiceRejectsVisitPlanWithoutValidityWindow(t *testing.T) {
	repository := &memoryMembershipPlanRepository{plans: map[domain.MembershipPlanID]domain.MembershipPlan{}}
	service, err := NewMembershipPlanService(repository, mustGymID(t), time.Now)
	if err != nil {
		t.Fatalf("NewMembershipPlanService returned error: %v", err)
	}
	input := planInput("10 visits", "visits", 10)
	input.DurationValue = 0
	if _, err := service.Create(context.Background(), input); !errors.Is(err, domain.ErrInvalidDurationValue) {
		t.Fatalf("Create error = %v, want %v", err, domain.ErrInvalidDurationValue)
	}
	if len(repository.plans) != 0 {
		t.Fatal("repository persisted an invalid plan")
	}
}

func planInput(name, kind string, visitLimit int) MembershipPlanInput {
	return MembershipPlanInput{Name: name, ValidityKind: kind, DurationValue: 2, DurationUnit: "months", VisitLimit: visitLimit, PriceCents: 5000, Currency: "USD", Status: "active"}
}

type memoryMembershipPlanRepository struct {
	plans    map[domain.MembershipPlanID]domain.MembershipPlan
	archived map[domain.MembershipPlanID]bool
}

func (r *memoryMembershipPlanRepository) Create(_ context.Context, plan domain.MembershipPlan) error {
	r.plans[plan.ID()] = plan
	return nil
}

func (r *memoryMembershipPlanRepository) Get(_ context.Context, gymID domain.GymID, planID domain.MembershipPlanID) (domain.MembershipPlan, error) {
	plan, ok := r.plans[planID]
	if !ok || plan.GymID() != gymID || r.archived[planID] {
		return domain.MembershipPlan{}, ports.ErrMembershipPlanNotFound
	}
	return plan, nil
}

func (r *memoryMembershipPlanRepository) List(_ context.Context, gymID domain.GymID) ([]domain.MembershipPlan, error) {
	var plans []domain.MembershipPlan
	for id, plan := range r.plans {
		if plan.GymID() == gymID && !r.archived[id] {
			plans = append(plans, plan)
		}
	}
	return plans, nil
}

func (r *memoryMembershipPlanRepository) Update(_ context.Context, plan domain.MembershipPlan) error {
	if _, ok := r.plans[plan.ID()]; !ok || r.archived[plan.ID()] {
		return ports.ErrMembershipPlanNotFound
	}
	r.plans[plan.ID()] = plan
	return nil
}

func (r *memoryMembershipPlanRepository) Archive(_ context.Context, _ domain.GymID, planID domain.MembershipPlanID, _ time.Time) error {
	if _, ok := r.plans[planID]; !ok || r.archived[planID] {
		return ports.ErrMembershipPlanNotFound
	}
	if r.archived == nil {
		r.archived = map[domain.MembershipPlanID]bool{}
	}
	r.archived[planID] = true
	return nil
}
