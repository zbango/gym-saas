package domain

import (
	"errors"
	"testing"
	"time"
)

func TestRecordExpenseIsSeparateFromMemberPayment(t *testing.T) {
	amount, err := NewMoney(1250, "USD")
	if err != nil {
		t.Fatalf("NewMoney returned error: %v", err)
	}
	paidAt := time.Date(2026, time.August, 20, 10, 0, 0, 0, time.UTC)
	expense, err := RecordExpense(mustGymID(t), amount, PaymentMethodBankTransfer, " invoice-22 ", " cleaning supplies ", paidAt)
	if err != nil {
		t.Fatalf("RecordExpense returned error: %v", err)
	}
	if expense.Status() != ExpenseStatusPosted || expense.Reference() != "invoice-22" || expense.Notes() != "cleaning supplies" || expense.Amount().Cents() != 1250 {
		t.Fatalf("expense = %#v", expense)
	}
}

func TestPendingExpensePostsAndVoids(t *testing.T) {
	now := time.Date(2026, time.August, 20, 10, 0, 0, 0, time.UTC)
	amount, err := NewMoney(500, "USD")
	if err != nil {
		t.Fatalf("NewMoney returned error: %v", err)
	}
	pending, err := CreatePendingExpense(mustGymID(t), amount, PaymentMethodCash, "", "", now)
	if err != nil {
		t.Fatalf("CreatePendingExpense returned error: %v", err)
	}
	if _, err := pending.Void(now.Add(time.Hour)); err != nil {
		t.Fatalf("Void returned error: %v", err)
	}
	posted, err := pending.Post(now.Add(time.Hour))
	if err != nil {
		t.Fatalf("Post returned error: %v", err)
	}
	if posted.Status() != ExpenseStatusPosted || posted.PaidAt().IsZero() {
		t.Fatalf("posted expense = %#v", posted)
	}
	if _, err := posted.Void(now.Add(2 * time.Hour)); !errors.Is(err, ErrExpenseCannotVoid) {
		t.Fatalf("Void posted expense error = %v, want %v", err, ErrExpenseCannotVoid)
	}
}
