package sqlite

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/zbango/gym-saas/go/core/domain"
	"github.com/zbango/gym-saas/go/core/ports"
)

func TestMembershipPaymentWriterCommitsMembershipAndPaymentTogether(t *testing.T) {
	store := openTestStore(t, filepath.Join(t.TempDir(), "gym-saas.db"))
	defer store.Close()
	membership, gymID := seedVisitMembership(t, store)
	payment, err := domain.RecordPayment(gymID, membership.MemberID(), membership.ID(), domain.PaymentKindInitial, mustMoney(t, 4500), domain.PaymentMethodCash, "receipt-1", "", membership.StartsAt())
	if err != nil {
		t.Fatalf("RecordPayment returned error: %v", err)
	}
	writer := NewMembershipPaymentWriter(store)
	if err := writer.CreateMembershipAndPayment(context.Background(), membership, payment); err != nil {
		t.Fatalf("CreateMembershipAndPayment returned error: %v", err)
	}
	if _, err := NewMembershipRepository(store).Get(context.Background(), gymID, membership.ID()); err != nil {
		t.Fatalf("Get membership returned error: %v", err)
	}
	if _, err := NewPaymentRepository(store).Get(context.Background(), gymID, payment.ID()); err != nil {
		t.Fatalf("Get payment returned error: %v", err)
	}
}

func TestMembershipPaymentWriterRollsBackMembershipWhenPaymentFails(t *testing.T) {
	store := openTestStore(t, filepath.Join(t.TempDir(), "gym-saas.db"))
	defer store.Close()
	firstMembership, gymID := seedVisitMembership(t, store)
	if err := NewMembershipRepository(store).Create(context.Background(), firstMembership); err != nil {
		t.Fatalf("seed membership: %v", err)
	}
	payment, err := domain.RecordPayment(gymID, firstMembership.MemberID(), firstMembership.ID(), domain.PaymentKindInitial, mustMoney(t, 4500), domain.PaymentMethodCash, "receipt-1", "", firstMembership.StartsAt())
	if err != nil {
		t.Fatalf("RecordPayment returned error: %v", err)
	}
	if err := NewPaymentRepository(store).Create(context.Background(), payment); err != nil {
		t.Fatalf("seed payment: %v", err)
	}
	secondPlan := mustMembershipPlan(t, gymID, "Second visit plan", domain.MembershipValidityVisits, 2)
	if err := NewMembershipPlanRepository(store).Create(context.Background(), secondPlan); err != nil {
		t.Fatalf("seed second plan: %v", err)
	}
	secondMembership, err := domain.StartMembership(mustGymForRepository(t, gymID), firstMembership.MemberID(), secondPlan, firstMembership.StartsAt().AddDate(0, 1, 0), firstMembership.StartsAt().AddDate(0, 1, 0))
	if err != nil {
		t.Fatalf("StartMembership returned error: %v", err)
	}
	duplicatePayment, err := domain.NewPayment(payment.ID(), gymID, secondMembership.MemberID(), secondMembership.ID(), payment.Status(), payment.Kind(), payment.Amount(), payment.Method(), payment.Reference(), payment.Notes(), payment.PaidAt(), payment.CreatedAt(), payment.UpdatedAt())
	if err != nil {
		t.Fatalf("NewPayment returned error: %v", err)
	}

	if err := NewMembershipPaymentWriter(store).CreateMembershipAndPayment(context.Background(), secondMembership, duplicatePayment); err == nil {
		t.Fatal("CreateMembershipAndPayment succeeded with a duplicate payment ID")
	}
	if _, err := NewMembershipRepository(store).Get(context.Background(), gymID, secondMembership.ID()); !errors.Is(err, ports.ErrMembershipNotFound) {
		t.Fatalf("Get rolled-back membership error = %v, want %v", err, ports.ErrMembershipNotFound)
	}
}
