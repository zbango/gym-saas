package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zbango/gym-saas/go/core/domain"
	"github.com/zbango/gym-saas/go/core/ports"
)

var ErrPaymentMembershipMemberMismatch = errors.New("payment membership belongs to another member")

type RecordPaymentInput struct {
	MemberID      string
	MembershipID  string
	Kind          string
	AmountCents   int64
	Currency      string
	PaymentMethod string
	Reference     string
	Notes         string
}

// PaymentService coordinates standalone receipt lifecycle operations for one
// gym. Atomic membership purchases use MembershipPurchaseService instead.
type PaymentService struct {
	memberships ports.MembershipRepository
	members     ports.MemberRepository
	payments    ports.PaymentRepository
	gymID       domain.GymID
	now         func() time.Time
}

func NewPaymentService(members ports.MemberRepository, memberships ports.MembershipRepository, payments ports.PaymentRepository, gymID domain.GymID, now func() time.Time) (*PaymentService, error) {
	if members == nil {
		return nil, fmt.Errorf("member repository is required")
	}
	if memberships == nil {
		return nil, fmt.Errorf("membership repository is required")
	}
	if payments == nil {
		return nil, fmt.Errorf("payment repository is required")
	}
	if _, err := domain.ParseGymID(string(gymID)); err != nil {
		return nil, err
	}
	if now == nil {
		now = time.Now
	}
	return &PaymentService{members: members, memberships: memberships, payments: payments, gymID: gymID, now: now}, nil
}

func (s *PaymentService) Record(ctx context.Context, input RecordPaymentInput) (domain.Payment, error) {
	return s.create(ctx, input, false)
}

func (s *PaymentService) CreatePending(ctx context.Context, input RecordPaymentInput) (domain.Payment, error) {
	return s.create(ctx, input, true)
}

func (s *PaymentService) create(ctx context.Context, input RecordPaymentInput, pending bool) (domain.Payment, error) {
	memberID, err := domain.ParseMemberID(input.MemberID)
	if err != nil {
		return domain.Payment{}, err
	}
	if _, err := s.members.Get(ctx, s.gymID, memberID); err != nil {
		return domain.Payment{}, fmt.Errorf("find member: %w", err)
	}
	membershipID, err := s.validateMembership(ctx, input.MembershipID, memberID)
	if err != nil {
		return domain.Payment{}, err
	}
	kind, err := domain.ParsePaymentKind(input.Kind)
	if err != nil {
		return domain.Payment{}, err
	}
	amount, err := domain.NewMoney(input.AmountCents, input.Currency)
	if err != nil {
		return domain.Payment{}, err
	}
	method, err := domain.ParsePaymentMethod(input.PaymentMethod)
	if err != nil {
		return domain.Payment{}, err
	}
	createdAt := s.now()
	var payment domain.Payment
	if pending {
		payment, err = domain.CreatePendingPayment(s.gymID, memberID, membershipID, kind, amount, method, input.Reference, input.Notes, createdAt)
	} else {
		payment, err = domain.RecordPayment(s.gymID, memberID, membershipID, kind, amount, method, input.Reference, input.Notes, createdAt)
	}
	if err != nil {
		return domain.Payment{}, err
	}
	if err := s.payments.Create(ctx, payment); err != nil {
		return domain.Payment{}, fmt.Errorf("save payment: %w", err)
	}
	return payment, nil
}

func (s *PaymentService) Post(ctx context.Context, id string) (domain.Payment, error) {
	payment, err := s.find(ctx, id)
	if err != nil {
		return domain.Payment{}, err
	}
	updated, err := payment.Post(s.now())
	if err != nil {
		return domain.Payment{}, err
	}
	if err := s.payments.Update(ctx, updated); err != nil {
		return domain.Payment{}, fmt.Errorf("post payment: %w", err)
	}
	return updated, nil
}

func (s *PaymentService) Void(ctx context.Context, id string) (domain.Payment, error) {
	payment, err := s.find(ctx, id)
	if err != nil {
		return domain.Payment{}, err
	}
	updated, err := payment.Void(s.now())
	if err != nil {
		return domain.Payment{}, err
	}
	if err := s.payments.Update(ctx, updated); err != nil {
		return domain.Payment{}, fmt.Errorf("void payment: %w", err)
	}
	return updated, nil
}

func (s *PaymentService) Refund(ctx context.Context, id string) (domain.Payment, error) {
	payment, err := s.find(ctx, id)
	if err != nil {
		return domain.Payment{}, err
	}
	updated, err := payment.Refund(s.now())
	if err != nil {
		return domain.Payment{}, err
	}
	if err := s.payments.Update(ctx, updated); err != nil {
		return domain.Payment{}, fmt.Errorf("refund payment: %w", err)
	}
	return updated, nil
}

func (s *PaymentService) Get(ctx context.Context, id string) (domain.Payment, error) {
	return s.find(ctx, id)
}

func (s *PaymentService) ListForMember(ctx context.Context, memberID string) ([]domain.Payment, error) {
	parsedMemberID, err := domain.ParseMemberID(memberID)
	if err != nil {
		return nil, err
	}
	if _, err := s.members.Get(ctx, s.gymID, parsedMemberID); err != nil {
		return nil, fmt.Errorf("find member: %w", err)
	}
	payments, err := s.payments.ListForMember(ctx, s.gymID, parsedMemberID)
	if err != nil {
		return nil, fmt.Errorf("list member payments: %w", err)
	}
	return payments, nil
}

func (s *PaymentService) find(ctx context.Context, id string) (domain.Payment, error) {
	paymentID, err := domain.ParsePaymentID(id)
	if err != nil {
		return domain.Payment{}, err
	}
	payment, err := s.payments.Get(ctx, s.gymID, paymentID)
	if err != nil {
		return domain.Payment{}, fmt.Errorf("find payment: %w", err)
	}
	return payment, nil
}

func (s *PaymentService) validateMembership(ctx context.Context, id string, memberID domain.MemberID) (domain.MembershipID, error) {
	if id == "" {
		return "", nil
	}
	membershipID, err := domain.ParseMembershipID(id)
	if err != nil {
		return "", err
	}
	membership, err := s.memberships.Get(ctx, s.gymID, membershipID)
	if err != nil {
		return "", fmt.Errorf("find membership: %w", err)
	}
	if membership.MemberID() != memberID {
		return "", ErrPaymentMembershipMemberMismatch
	}
	return membershipID, nil
}
