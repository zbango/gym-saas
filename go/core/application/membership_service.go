package application

import (
	"context"
	"fmt"
	"time"

	"github.com/zbango/gym-saas/go/core/domain"
	"github.com/zbango/gym-saas/go/core/ports"
)

// StartMembershipInput contains only the caller-selected references and an
// optional future start instant. Membership rules remain in domain.Membership.
type StartMembershipInput struct {
	MemberID         string
	MembershipPlanID string
	StartsAt         time.Time
}

// MembershipService coordinates membership lifecycle operations for one gym.
// It receives the full Gym value so calendar calculations cannot fall back to
// a hidden UTC/default timezone.
type MembershipService struct {
	memberships ports.MembershipRepository
	members     ports.MemberRepository
	plans       ports.MembershipPlanRepository
	gym         domain.Gym
	now         func() time.Time
}

func NewMembershipService(members ports.MemberRepository, plans ports.MembershipPlanRepository, memberships ports.MembershipRepository, gym domain.Gym, now func() time.Time) (*MembershipService, error) {
	if members == nil {
		return nil, fmt.Errorf("member repository is required")
	}
	if plans == nil {
		return nil, fmt.Errorf("membership plan repository is required")
	}
	if memberships == nil {
		return nil, fmt.Errorf("membership repository is required")
	}
	if _, err := domain.ParseGymID(string(gym.ID())); err != nil {
		return nil, err
	}
	if now == nil {
		now = time.Now
	}
	return &MembershipService{members: members, plans: plans, memberships: memberships, gym: gym, now: now}, nil
}

func (s *MembershipService) Start(ctx context.Context, input StartMembershipInput) (domain.Membership, error) {
	memberID, err := domain.ParseMemberID(input.MemberID)
	if err != nil {
		return domain.Membership{}, err
	}
	planID, err := domain.ParseMembershipPlanID(input.MembershipPlanID)
	if err != nil {
		return domain.Membership{}, err
	}
	if _, err := s.members.Get(ctx, s.gym.ID(), memberID); err != nil {
		return domain.Membership{}, fmt.Errorf("find member: %w", err)
	}
	plan, err := s.plans.Get(ctx, s.gym.ID(), planID)
	if err != nil {
		return domain.Membership{}, fmt.Errorf("find membership plan: %w", err)
	}
	createdAt := s.now()
	startsAt := input.StartsAt
	if startsAt.IsZero() {
		startsAt = createdAt
	}
	membership, err := domain.StartMembership(s.gym, memberID, plan, startsAt, createdAt)
	if err != nil {
		return domain.Membership{}, err
	}
	if err := s.memberships.Create(ctx, membership); err != nil {
		return domain.Membership{}, fmt.Errorf("save membership: %w", err)
	}
	return membership, nil
}

func (s *MembershipService) Activate(ctx context.Context, id string) (domain.Membership, error) {
	membership, err := s.find(ctx, id)
	if err != nil {
		return domain.Membership{}, err
	}
	updated, err := membership.Activate(s.now())
	if err != nil {
		return domain.Membership{}, err
	}
	if err := s.memberships.Update(ctx, updated); err != nil {
		return domain.Membership{}, fmt.Errorf("activate membership: %w", err)
	}
	return updated, nil
}

func (s *MembershipService) Cancel(ctx context.Context, id, reason string) (domain.Membership, error) {
	membership, err := s.find(ctx, id)
	if err != nil {
		return domain.Membership{}, err
	}
	updated, err := membership.Cancel(s.now(), reason)
	if err != nil {
		return domain.Membership{}, err
	}
	if err := s.memberships.Update(ctx, updated); err != nil {
		return domain.Membership{}, fmt.Errorf("cancel membership: %w", err)
	}
	return updated, nil
}

func (s *MembershipService) Get(ctx context.Context, id string) (domain.Membership, error) {
	return s.find(ctx, id)
}

func (s *MembershipService) ListForMember(ctx context.Context, memberID string) ([]domain.Membership, error) {
	parsedMemberID, err := domain.ParseMemberID(memberID)
	if err != nil {
		return nil, err
	}
	if _, err := s.members.Get(ctx, s.gym.ID(), parsedMemberID); err != nil {
		return nil, fmt.Errorf("find member: %w", err)
	}
	memberships, err := s.memberships.ListForMember(ctx, s.gym.ID(), parsedMemberID)
	if err != nil {
		return nil, fmt.Errorf("list member memberships: %w", err)
	}
	return memberships, nil
}

func (s *MembershipService) find(ctx context.Context, id string) (domain.Membership, error) {
	membershipID, err := domain.ParseMembershipID(id)
	if err != nil {
		return domain.Membership{}, err
	}
	membership, err := s.memberships.Get(ctx, s.gym.ID(), membershipID)
	if err != nil {
		return domain.Membership{}, fmt.Errorf("find membership: %w", err)
	}
	return membership, nil
}
