package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/zbango/gym-saas/go/core/domain"
	"github.com/zbango/gym-saas/go/core/ports"
)

func TestMembershipPurchaseServiceCreatesMembershipAndInitialReceiptAtomically(t *testing.T) {
	gym, member, plan := membershipServiceFixture(t)
	members := &memoryMemberRepository{members: map[domain.MemberID]domain.Member{member.ID(): member}}
	plans := &memoryMembershipPlanRepository{plans: map[domain.MembershipPlanID]domain.MembershipPlan{plan.ID(): plan}}
	writer := &memoryMembershipPaymentWriter{}
	now := time.Date(2026, time.August, 1, 9, 0, 0, 0, time.UTC)
	service, err := NewMembershipPurchaseService(members, plans, writer, gym, func() time.Time { return now })
	if err != nil {
		t.Fatalf("NewMembershipPurchaseService returned error: %v", err)
	}

	result, err := service.Purchase(context.Background(), PurchaseMembershipInput{
		MemberID:         string(member.ID()),
		MembershipPlanID: string(plan.ID()),
		AmountCents:      4000,
		Currency:         "USD",
		PaymentMethod:    "cash",
		Reference:        "receipt-1",
	})
	if err != nil {
		t.Fatalf("Purchase returned error: %v", err)
	}
	if writer.membership.ID() != result.Membership.ID() || writer.payment.ID() != result.Payment.ID() || result.Payment.MembershipID() != result.Membership.ID() || result.Payment.MemberID() != member.ID() || result.Payment.Amount().Cents() != 4000 || result.Payment.Kind() != domain.PaymentKindInitial {
		t.Fatalf("purchase result = %#v, writer membership/payment = %#v / %#v", result, writer.membership, writer.payment)
	}
}

func TestMembershipPurchaseServiceRejectsUnknownMemberBeforeWriter(t *testing.T) {
	gym, _, plan := membershipServiceFixture(t)
	members := &memoryMemberRepository{members: map[domain.MemberID]domain.Member{}}
	plans := &memoryMembershipPlanRepository{plans: map[domain.MembershipPlanID]domain.MembershipPlan{plan.ID(): plan}}
	writer := &memoryMembershipPaymentWriter{}
	service, err := NewMembershipPurchaseService(members, plans, writer, gym, time.Now)
	if err != nil {
		t.Fatalf("NewMembershipPurchaseService returned error: %v", err)
	}
	memberID, err := domain.NewMemberID()
	if err != nil {
		t.Fatalf("NewMemberID returned error: %v", err)
	}
	_, err = service.Purchase(context.Background(), PurchaseMembershipInput{MemberID: string(memberID), MembershipPlanID: string(plan.ID()), AmountCents: 4500, Currency: "USD", PaymentMethod: "cash"})
	if !errors.Is(err, ports.ErrMemberNotFound) {
		t.Fatalf("Purchase error = %v, want %v", err, ports.ErrMemberNotFound)
	}
	if writer.called {
		t.Fatal("writer ran for an unknown member")
	}
}

type memoryMembershipPaymentWriter struct {
	membership domain.Membership
	payment    domain.Payment
	called     bool
}

func (w *memoryMembershipPaymentWriter) CreateMembershipAndPayment(_ context.Context, membership domain.Membership, payment domain.Payment) error {
	w.membership = membership
	w.payment = payment
	w.called = true
	return nil
}
