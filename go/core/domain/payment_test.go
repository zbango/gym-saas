package domain

import (
	"errors"
	"testing"
	"time"
)

func TestRecordPaymentCreatesPostedReceipt(t *testing.T) {
	gymID := mustGymID(t)
	memberID := mustMemberID(t)
	membershipID := mustMembershipID(t)
	amount, err := NewMoney(4500, "USD")
	if err != nil {
		t.Fatalf("NewMoney returned error: %v", err)
	}
	paidAt := time.Date(2026, time.August, 20, 10, 0, 0, 0, time.FixedZone("ECT", -5*60*60))
	payment, err := RecordPayment(gymID, memberID, membershipID, PaymentKindInitial, amount, PaymentMethodCash, "  receipt-100  ", "  first month  ", paidAt)
	if err != nil {
		t.Fatalf("RecordPayment returned error: %v", err)
	}
	if err := ValidateUUID(string(payment.ID())); err != nil {
		t.Fatalf("payment ID is invalid: %v", err)
	}
	if payment.Status() != PaymentStatusPosted || payment.MembershipID() != membershipID || payment.Reference() != "receipt-100" || payment.Notes() != "first month" || !payment.PaidAt().Equal(paidAt.UTC()) {
		t.Fatalf("payment = %#v", payment)
	}
}

func TestPendingPaymentPostsVoidsAndRefundsWithoutNegativeMoney(t *testing.T) {
	gymID := mustGymID(t)
	memberID := mustMemberID(t)
	amount, err := NewMoney(4500, "USD")
	if err != nil {
		t.Fatalf("NewMoney returned error: %v", err)
	}
	createdAt := time.Date(2026, time.August, 20, 10, 0, 0, 0, time.UTC)
	pending, err := CreatePendingPayment(gymID, memberID, "", PaymentKindRenewal, amount, PaymentMethodBankTransfer, "", "", createdAt)
	if err != nil {
		t.Fatalf("CreatePendingPayment returned error: %v", err)
	}
	posted, err := pending.Post(createdAt.Add(time.Hour))
	if err != nil {
		t.Fatalf("Post returned error: %v", err)
	}
	refunded, err := posted.Refund(createdAt.Add(2 * time.Hour))
	if err != nil {
		t.Fatalf("Refund returned error: %v", err)
	}
	if refunded.Status() != PaymentStatusRefunded || refunded.Amount().Cents() != 4500 || !refunded.PaidAt().Equal(posted.PaidAt()) {
		t.Fatalf("refunded payment = %#v", refunded)
	}
	if _, err := refunded.Refund(createdAt.Add(3 * time.Hour)); !errors.Is(err, ErrPaymentCannotRefund) {
		t.Fatalf("repeat Refund error = %v, want %v", err, ErrPaymentCannotRefund)
	}

	voided, err := pending.Void(createdAt.Add(time.Hour))
	if err != nil {
		t.Fatalf("Void returned error: %v", err)
	}
	if voided.Status() != PaymentStatusVoided || !voided.PaidAt().IsZero() {
		t.Fatalf("voided payment = %#v", voided)
	}
}

func TestPaymentRejectsInvalidMoneyAndLifecycle(t *testing.T) {
	gymID := mustGymID(t)
	memberID := mustMemberID(t)
	zero, err := NewMoney(0, "USD")
	if err != nil {
		t.Fatalf("NewMoney returned error: %v", err)
	}
	now := time.Date(2026, time.August, 20, 10, 0, 0, 0, time.UTC)
	if _, err := RecordPayment(gymID, memberID, "", PaymentKindInitial, zero, PaymentMethodCash, "", "", now); !errors.Is(err, ErrPaymentAmountMustBePositive) {
		t.Fatalf("RecordPayment zero amount error = %v, want %v", err, ErrPaymentAmountMustBePositive)
	}
	amount, err := NewMoney(100, "USD")
	if err != nil {
		t.Fatalf("NewMoney returned error: %v", err)
	}
	if _, err := NewPayment(mustPaymentID(t), gymID, memberID, "", PaymentStatusPosted, PaymentKindInitial, amount, PaymentMethodCash, "", "", time.Time{}, now, now); !errors.Is(err, ErrInvalidPaymentLifecycle) {
		t.Fatalf("NewPayment posted without paid_at error = %v, want %v", err, ErrInvalidPaymentLifecycle)
	}
}

func mustPaymentID(t *testing.T) PaymentID {
	t.Helper()
	id, err := NewPaymentID()
	if err != nil {
		t.Fatalf("NewPaymentID returned error: %v", err)
	}
	return id
}

func mustMembershipID(t *testing.T) MembershipID {
	t.Helper()
	id, err := NewMembershipID()
	if err != nil {
		t.Fatalf("NewMembershipID returned error: %v", err)
	}
	return id
}
