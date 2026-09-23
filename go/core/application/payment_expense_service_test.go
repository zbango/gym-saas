package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/zbango/gym-saas/go/core/domain"
	"github.com/zbango/gym-saas/go/core/ports"
)

func TestPaymentServiceRecordsRefundsAndListsMemberReceipts(t *testing.T) {
	gym, member, plan := membershipServiceFixture(t)
	members := &memoryMemberRepository{members: map[domain.MemberID]domain.Member{member.ID(): member}}
	membership, err := domain.StartMembership(gym, member.ID(), plan, time.Date(2026, time.August, 1, 9, 0, 0, 0, time.UTC), time.Date(2026, time.August, 1, 9, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("StartMembership returned error: %v", err)
	}
	memberships := &memoryMembershipRepository{memberships: map[domain.MembershipID]domain.Membership{membership.ID(): membership}}
	payments := &memoryPaymentRepository{payments: map[domain.PaymentID]domain.Payment{}}
	times := []time.Time{
		time.Date(2026, time.August, 1, 10, 0, 0, 0, time.UTC),
		time.Date(2026, time.August, 2, 10, 0, 0, 0, time.UTC),
	}
	service, err := NewPaymentService(members, memberships, payments, gym.ID(), func() time.Time {
		value := times[0]
		times = times[1:]
		return value
	})
	if err != nil {
		t.Fatalf("NewPaymentService returned error: %v", err)
	}

	payment, err := service.Record(context.Background(), RecordPaymentInput{
		MemberID: string(member.ID()), MembershipID: string(membership.ID()), Kind: "initial",
		AmountCents: 4500, Currency: "USD", PaymentMethod: "cash", Reference: "receipt-1",
	})
	if err != nil {
		t.Fatalf("Record returned error: %v", err)
	}
	refunded, err := service.Refund(context.Background(), string(payment.ID()))
	if err != nil {
		t.Fatalf("Refund returned error: %v", err)
	}
	if refunded.Status() != domain.PaymentStatusRefunded || payments.payments[payment.ID()].Status() != domain.PaymentStatusRefunded {
		t.Fatalf("refunded payment = %#v", refunded)
	}
	listed, err := service.ListForMember(context.Background(), string(member.ID()))
	if err != nil {
		t.Fatalf("ListForMember returned error: %v", err)
	}
	if len(listed) != 1 || listed[0].ID() != payment.ID() {
		t.Fatalf("ListForMember = %#v", listed)
	}
}

func TestPaymentServiceRejectsMembershipOwnedByAnotherMember(t *testing.T) {
	gym, member, plan := membershipServiceFixture(t)
	otherMember, err := domain.CreateMember(gym.ID(), "Grace", "Hopper", "grace@example.com", "555-0101", "0203040506", "1906-12-09", "", domain.MemberStatusActive, time.Date(2026, time.August, 1, 8, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("CreateMember returned error: %v", err)
	}
	membership, err := domain.StartMembership(gym, otherMember.ID(), plan, time.Date(2026, time.August, 1, 9, 0, 0, 0, time.UTC), time.Date(2026, time.August, 1, 9, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("StartMembership returned error: %v", err)
	}
	members := &memoryMemberRepository{members: map[domain.MemberID]domain.Member{member.ID(): member, otherMember.ID(): otherMember}}
	memberships := &memoryMembershipRepository{memberships: map[domain.MembershipID]domain.Membership{membership.ID(): membership}}
	payments := &memoryPaymentRepository{payments: map[domain.PaymentID]domain.Payment{}}
	service, err := NewPaymentService(members, memberships, payments, gym.ID(), time.Now)
	if err != nil {
		t.Fatalf("NewPaymentService returned error: %v", err)
	}
	_, err = service.Record(context.Background(), RecordPaymentInput{MemberID: string(member.ID()), MembershipID: string(membership.ID()), Kind: "initial", AmountCents: 4500, Currency: "USD", PaymentMethod: "cash"})
	if !errors.Is(err, ErrPaymentMembershipMemberMismatch) {
		t.Fatalf("Record error = %v, want %v", err, ErrPaymentMembershipMemberMismatch)
	}
	if len(payments.payments) != 0 {
		t.Fatal("mismatched membership created a payment")
	}
}

func TestPaymentServicePostsAndVoidsPendingPayments(t *testing.T) {
	gym, member, _ := membershipServiceFixture(t)
	members := &memoryMemberRepository{members: map[domain.MemberID]domain.Member{member.ID(): member}}
	memberships := &memoryMembershipRepository{memberships: map[domain.MembershipID]domain.Membership{}}
	payments := &memoryPaymentRepository{payments: map[domain.PaymentID]domain.Payment{}}
	times := []time.Time{
		time.Date(2026, time.August, 1, 10, 0, 0, 0, time.UTC),
		time.Date(2026, time.August, 1, 11, 0, 0, 0, time.UTC),
		time.Date(2026, time.August, 1, 12, 0, 0, 0, time.UTC),
		time.Date(2026, time.August, 1, 13, 0, 0, 0, time.UTC),
	}
	service, err := NewPaymentService(members, memberships, payments, gym.ID(), func() time.Time {
		value := times[0]
		times = times[1:]
		return value
	})
	if err != nil {
		t.Fatalf("NewPaymentService returned error: %v", err)
	}
	pending, err := service.CreatePending(context.Background(), RecordPaymentInput{MemberID: string(member.ID()), Kind: "other", AmountCents: 1000, Currency: "USD", PaymentMethod: "cash"})
	if err != nil {
		t.Fatalf("CreatePending returned error: %v", err)
	}
	posted, err := service.Post(context.Background(), string(pending.ID()))
	if err != nil {
		t.Fatalf("Post returned error: %v", err)
	}
	if posted.Status() != domain.PaymentStatusPosted {
		t.Fatalf("posted payment = %#v", posted)
	}
	pending, err = service.CreatePending(context.Background(), RecordPaymentInput{MemberID: string(member.ID()), Kind: "other", AmountCents: 1000, Currency: "USD", PaymentMethod: "cash"})
	if err != nil {
		t.Fatalf("CreatePending returned error: %v", err)
	}
	voided, err := service.Void(context.Background(), string(pending.ID()))
	if err != nil {
		t.Fatalf("Void returned error: %v", err)
	}
	if voided.Status() != domain.PaymentStatusVoided {
		t.Fatalf("voided payment = %#v", voided)
	}
}

func TestExpenseServiceRecordsAndVoidsExpense(t *testing.T) {
	gym, _, _ := membershipServiceFixture(t)
	expenses := &memoryExpenseRepository{expenses: map[domain.ExpenseID]domain.Expense{}}
	times := []time.Time{
		time.Date(2026, time.August, 1, 10, 0, 0, 0, time.UTC),
		time.Date(2026, time.August, 2, 10, 0, 0, 0, time.UTC),
	}
	service, err := NewExpenseService(expenses, gym.ID(), func() time.Time {
		value := times[0]
		times = times[1:]
		return value
	})
	if err != nil {
		t.Fatalf("NewExpenseService returned error: %v", err)
	}
	expense, err := service.CreatePending(context.Background(), RecordExpenseInput{AmountCents: 1250, Currency: "USD", PaymentMethod: "bank_transfer", Reference: "vendor-1"})
	if err != nil {
		t.Fatalf("CreatePending returned error: %v", err)
	}
	voided, err := service.Void(context.Background(), string(expense.ID()))
	if err != nil {
		t.Fatalf("Void returned error: %v", err)
	}
	if voided.Status() != domain.ExpenseStatusVoided || expenses.expenses[expense.ID()].Status() != domain.ExpenseStatusVoided {
		t.Fatalf("voided expense = %#v", voided)
	}
}

func TestExpenseServicePostsPendingExpense(t *testing.T) {
	gym, _, _ := membershipServiceFixture(t)
	expenses := &memoryExpenseRepository{expenses: map[domain.ExpenseID]domain.Expense{}}
	times := []time.Time{
		time.Date(2026, time.August, 1, 10, 0, 0, 0, time.UTC),
		time.Date(2026, time.August, 2, 10, 0, 0, 0, time.UTC),
	}
	service, err := NewExpenseService(expenses, gym.ID(), func() time.Time {
		value := times[0]
		times = times[1:]
		return value
	})
	if err != nil {
		t.Fatalf("NewExpenseService returned error: %v", err)
	}
	pending, err := service.CreatePending(context.Background(), RecordExpenseInput{AmountCents: 1250, Currency: "USD", PaymentMethod: "cash"})
	if err != nil {
		t.Fatalf("CreatePending returned error: %v", err)
	}
	posted, err := service.Post(context.Background(), string(pending.ID()))
	if err != nil {
		t.Fatalf("Post returned error: %v", err)
	}
	if posted.Status() != domain.ExpenseStatusPosted {
		t.Fatalf("posted expense = %#v", posted)
	}
}

type memoryPaymentRepository struct {
	payments map[domain.PaymentID]domain.Payment
}

func (r *memoryPaymentRepository) Create(_ context.Context, payment domain.Payment) error {
	r.payments[payment.ID()] = payment
	return nil
}

func (r *memoryPaymentRepository) Get(_ context.Context, gymID domain.GymID, paymentID domain.PaymentID) (domain.Payment, error) {
	payment, ok := r.payments[paymentID]
	if !ok || payment.GymID() != gymID {
		return domain.Payment{}, ports.ErrPaymentNotFound
	}
	return payment, nil
}

func (r *memoryPaymentRepository) ListForMember(_ context.Context, gymID domain.GymID, memberID domain.MemberID) ([]domain.Payment, error) {
	var result []domain.Payment
	for _, payment := range r.payments {
		if payment.GymID() == gymID && payment.MemberID() == memberID {
			result = append(result, payment)
		}
	}
	return result, nil
}

func (r *memoryPaymentRepository) Update(_ context.Context, payment domain.Payment) error {
	if _, ok := r.payments[payment.ID()]; !ok {
		return ports.ErrPaymentNotFound
	}
	r.payments[payment.ID()] = payment
	return nil
}

type memoryExpenseRepository struct {
	expenses map[domain.ExpenseID]domain.Expense
}

func (r *memoryExpenseRepository) Create(_ context.Context, expense domain.Expense) error {
	r.expenses[expense.ID()] = expense
	return nil
}

func (r *memoryExpenseRepository) Get(_ context.Context, gymID domain.GymID, expenseID domain.ExpenseID) (domain.Expense, error) {
	expense, ok := r.expenses[expenseID]
	if !ok || expense.GymID() != gymID {
		return domain.Expense{}, ports.ErrExpenseNotFound
	}
	return expense, nil
}

func (r *memoryExpenseRepository) List(_ context.Context, gymID domain.GymID) ([]domain.Expense, error) {
	var result []domain.Expense
	for _, expense := range r.expenses {
		if expense.GymID() == gymID {
			result = append(result, expense)
		}
	}
	return result, nil
}

func (r *memoryExpenseRepository) Update(_ context.Context, expense domain.Expense) error {
	if _, ok := r.expenses[expense.ID()]; !ok {
		return ports.ErrExpenseNotFound
	}
	r.expenses[expense.ID()] = expense
	return nil
}
