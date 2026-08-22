package application

import (
	"context"
	"fmt"
	"time"

	"github.com/zbango/gym-saas/go/core/domain"
	"github.com/zbango/gym-saas/go/core/ports"
)

type MembershipPlanInput struct {
	Name          string
	ValidityKind  string
	DurationValue int
	DurationUnit  string
	VisitLimit    int
	PriceCents    int64
	Currency      string
	Status        string
}

type MembershipPlanService struct {
	repository ports.MembershipPlanRepository
	gymID      domain.GymID
	now        func() time.Time
}

func NewMembershipPlanService(repository ports.MembershipPlanRepository, gymID domain.GymID, now func() time.Time) (*MembershipPlanService, error) {
	if repository == nil {
		return nil, fmt.Errorf("membership plan repository is required")
	}
	if _, err := domain.ParseGymID(string(gymID)); err != nil {
		return nil, err
	}
	if now == nil {
		now = time.Now
	}
	return &MembershipPlanService{repository: repository, gymID: gymID, now: now}, nil
}

func (s *MembershipPlanService) Create(ctx context.Context, input MembershipPlanInput) (domain.MembershipPlan, error) {
	plan, err := membershipPlanFromInput(s.gymID, "", time.Time{}, input, s.now())
	if err != nil {
		return domain.MembershipPlan{}, err
	}
	if err := s.repository.Create(ctx, plan); err != nil {
		return domain.MembershipPlan{}, fmt.Errorf("save membership plan: %w", err)
	}
	return plan, nil
}

func (s *MembershipPlanService) List(ctx context.Context) ([]domain.MembershipPlan, error) {
	plans, err := s.repository.List(ctx, s.gymID)
	if err != nil {
		return nil, fmt.Errorf("list membership plans: %w", err)
	}
	return plans, nil
}

func (s *MembershipPlanService) Update(ctx context.Context, id string, input MembershipPlanInput) (domain.MembershipPlan, error) {
	planID, err := domain.ParseMembershipPlanID(id)
	if err != nil {
		return domain.MembershipPlan{}, err
	}
	existing, err := s.repository.Get(ctx, s.gymID, planID)
	if err != nil {
		return domain.MembershipPlan{}, fmt.Errorf("find membership plan: %w", err)
	}
	plan, err := membershipPlanFromInput(s.gymID, existing.ID(), existing.CreatedAt(), input, s.now())
	if err != nil {
		return domain.MembershipPlan{}, err
	}
	if err := s.repository.Update(ctx, plan); err != nil {
		return domain.MembershipPlan{}, fmt.Errorf("update membership plan: %w", err)
	}
	return plan, nil
}

func (s *MembershipPlanService) Archive(ctx context.Context, id string) error {
	planID, err := domain.ParseMembershipPlanID(id)
	if err != nil {
		return err
	}
	if err := s.repository.Archive(ctx, s.gymID, planID, s.now()); err != nil {
		return fmt.Errorf("archive membership plan: %w", err)
	}
	return nil
}

func membershipPlanFromInput(gymID domain.GymID, existingID domain.MembershipPlanID, createdAt time.Time, input MembershipPlanInput, updatedAt time.Time) (domain.MembershipPlan, error) {
	kind, err := domain.ParseMembershipValidityKind(input.ValidityKind)
	if err != nil {
		return domain.MembershipPlan{}, err
	}
	unit, err := domain.ParseMembershipDurationUnit(input.DurationUnit)
	if err != nil {
		return domain.MembershipPlan{}, err
	}
	price, err := domain.NewMoney(input.PriceCents, input.Currency)
	if err != nil {
		return domain.MembershipPlan{}, err
	}
	status, err := domain.ParseMembershipPlanStatus(input.Status)
	if err != nil {
		return domain.MembershipPlan{}, err
	}
	if existingID == "" {
		return domain.CreateMembershipPlan(gymID, input.Name, kind, input.DurationValue, unit, input.VisitLimit, price, status, updatedAt)
	}
	return domain.NewMembershipPlan(existingID, gymID, input.Name, kind, input.DurationValue, unit, input.VisitLimit, price, status, createdAt, updatedAt)
}
