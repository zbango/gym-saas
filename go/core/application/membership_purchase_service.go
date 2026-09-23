package application

import (
	"context"
	"fmt"
	"time"

	"github.com/zbango/gym-saas/go/core/domain"
	"github.com/zbango/gym-saas/go/core/ports"
)

// PurchaseMembershipInput records the explicit initial receipt alongside a
// membership. Amount is deliberately caller-supplied: its relationship to the
// plan price (discount, deposit, or full settlement) is a later policy, never
// an implicit calculation in the UI.
type PurchaseMembershipInput struct {
	MemberID         string
	MembershipPlanID string
	StartsAt         time.Time
	AmountCents      int64
	Currency         string
	PaymentMethod    string
	Reference        string
	Notes            string
}

type MembershipPurchaseResult struct {
	Membership domain.Membership
	Payment    domain.Payment
}

// MembershipPurchaseService coordinates the one atomic R1 mutation that
// creates a membership and its initial receipt. The infrastructure writer owns
// the actual transaction boundary.
type MembershipPurchaseService struct {
	members ports.MemberRepository
	plans   ports.MembershipPlanRepository
	writer  ports.MembershipPaymentWriter
	gym     domain.Gym
	now     func() time.Time
}

func NewMembershipPurchaseService(members ports.MemberRepository, plans ports.MembershipPlanRepository, writer ports.MembershipPaymentWriter, gym domain.Gym, now func() time.Time) (*MembershipPurchaseService, error) {
	if members == nil {
		return nil, fmt.Errorf("member repository is required")
	}
	if plans == nil {
		return nil, fmt.Errorf("membership plan repository is required")
	}
	if writer == nil {
		return nil, fmt.Errorf("membership payment writer is required")
	}
	if _, err := domain.ParseGymID(string(gym.ID())); err != nil {
		return nil, err
	}
	if now == nil {
		now = time.Now
	}
	return &MembershipPurchaseService{members: members, plans: plans, writer: writer, gym: gym, now: now}, nil
}

func (s *MembershipPurchaseService) Purchase(ctx context.Context, input PurchaseMembershipInput) (MembershipPurchaseResult, error) {
	memberID, err := domain.ParseMemberID(input.MemberID)
	if err != nil {
		return MembershipPurchaseResult{}, err
	}
	planID, err := domain.ParseMembershipPlanID(input.MembershipPlanID)
	if err != nil {
		return MembershipPurchaseResult{}, err
	}
	if _, err := s.members.Get(ctx, s.gym.ID(), memberID); err != nil {
		return MembershipPurchaseResult{}, fmt.Errorf("find member: %w", err)
	}
	plan, err := s.plans.Get(ctx, s.gym.ID(), planID)
	if err != nil {
		return MembershipPurchaseResult{}, fmt.Errorf("find membership plan: %w", err)
	}
	amount, err := domain.NewMoney(input.AmountCents, input.Currency)
	if err != nil {
		return MembershipPurchaseResult{}, err
	}
	method, err := domain.ParsePaymentMethod(input.PaymentMethod)
	if err != nil {
		return MembershipPurchaseResult{}, err
	}
	createdAt := s.now()
	startsAt := input.StartsAt
	if startsAt.IsZero() {
		startsAt = createdAt
	}
	membership, err := domain.StartMembership(s.gym, memberID, plan, startsAt, createdAt)
	if err != nil {
		return MembershipPurchaseResult{}, err
	}
	payment, err := domain.RecordPayment(s.gym.ID(), memberID, membership.ID(), domain.PaymentKindInitial, amount, method, input.Reference, input.Notes, createdAt)
	if err != nil {
		return MembershipPurchaseResult{}, err
	}
	if err := s.writer.CreateMembershipAndPayment(ctx, membership, payment); err != nil {
		return MembershipPurchaseResult{}, fmt.Errorf("save membership and payment: %w", err)
	}
	return MembershipPurchaseResult{Membership: membership, Payment: payment}, nil
}
