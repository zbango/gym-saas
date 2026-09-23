package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidPaymentID            = errors.New("invalid payment ID")
	ErrInvalidPaymentStatus        = errors.New("invalid payment status")
	ErrInvalidPaymentKind          = errors.New("invalid payment kind")
	ErrInvalidPaymentMethod        = errors.New("invalid payment method")
	ErrPaymentAmountMustBePositive = errors.New("payment amount must be greater than zero")
	ErrInvalidPaymentLifecycle     = errors.New("invalid payment lifecycle")
	ErrPaymentCannotPost           = errors.New("payment cannot be posted in its current status")
	ErrPaymentCannotVoid           = errors.New("payment cannot be voided in its current status")
	ErrPaymentCannotRefund         = errors.New("payment cannot be refunded in its current status")
)

type PaymentID string

func NewPaymentID() (PaymentID, error) {
	value, err := NewUUID()
	if err != nil {
		return "", err
	}
	return PaymentID(value), nil
}

func ParsePaymentID(value string) (PaymentID, error) {
	if err := ValidateUUID(value); err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidPaymentID, err)
	}
	return PaymentID(value), nil
}

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusPosted   PaymentStatus = "posted"
	PaymentStatusVoided   PaymentStatus = "voided"
	PaymentStatusRefunded PaymentStatus = "refunded"
)

func ParsePaymentStatus(value string) (PaymentStatus, error) {
	switch status := PaymentStatus(value); status {
	case PaymentStatusPending, PaymentStatusPosted, PaymentStatusVoided, PaymentStatusRefunded:
		return status, nil
	default:
		return "", fmt.Errorf("%w: %q", ErrInvalidPaymentStatus, value)
	}
}

type PaymentKind string

const (
	PaymentKindInitial    PaymentKind = "initial"
	PaymentKindRenewal    PaymentKind = "renewal"
	PaymentKindAdjustment PaymentKind = "adjustment"
	PaymentKindOther      PaymentKind = "other"
)

func ParsePaymentKind(value string) (PaymentKind, error) {
	switch kind := PaymentKind(value); kind {
	case PaymentKindInitial, PaymentKindRenewal, PaymentKindAdjustment, PaymentKindOther:
		return kind, nil
	default:
		return "", fmt.Errorf("%w: %q", ErrInvalidPaymentKind, value)
	}
}

type PaymentMethod string

const (
	PaymentMethodCash         PaymentMethod = "cash"
	PaymentMethodCreditCard   PaymentMethod = "credit_card"
	PaymentMethodDebitCard    PaymentMethod = "debit_card"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodCheck        PaymentMethod = "check"
	PaymentMethodOther        PaymentMethod = "other"
)

func ParsePaymentMethod(value string) (PaymentMethod, error) {
	switch method := PaymentMethod(value); method {
	case PaymentMethodCash, PaymentMethodCreditCard, PaymentMethodDebitCard, PaymentMethodBankTransfer, PaymentMethodCheck, PaymentMethodOther:
		return method, nil
	default:
		return "", fmt.Errorf("%w: %q", ErrInvalidPaymentMethod, value)
	}
}

// Payment is an incoming member receipt. A refund changes a posted receipt's
// status but never turns its amount negative; the outgoing refund workflow is
// represented separately when accounting requirements are implemented.
type Payment struct {
	id           PaymentID
	gymID        GymID
	memberID     MemberID
	membershipID MembershipID
	status       PaymentStatus
	kind         PaymentKind
	amount       Money
	method       PaymentMethod
	reference    string
	notes        string
	paidAt       time.Time
	createdAt    time.Time
	updatedAt    time.Time
}

func RecordPayment(gymID GymID, memberID MemberID, membershipID MembershipID, kind PaymentKind, amount Money, method PaymentMethod, reference, notes string, paidAt time.Time) (Payment, error) {
	id, err := NewPaymentID()
	if err != nil {
		return Payment{}, fmt.Errorf("create payment ID: %w", err)
	}
	return NewPayment(id, gymID, memberID, membershipID, PaymentStatusPosted, kind, amount, method, reference, notes, paidAt, paidAt, paidAt)
}

func CreatePendingPayment(gymID GymID, memberID MemberID, membershipID MembershipID, kind PaymentKind, amount Money, method PaymentMethod, reference, notes string, createdAt time.Time) (Payment, error) {
	id, err := NewPaymentID()
	if err != nil {
		return Payment{}, fmt.Errorf("create payment ID: %w", err)
	}
	return NewPayment(id, gymID, memberID, membershipID, PaymentStatusPending, kind, amount, method, reference, notes, time.Time{}, createdAt, createdAt)
}

func NewPayment(id PaymentID, gymID GymID, memberID MemberID, membershipID MembershipID, status PaymentStatus, kind PaymentKind, amount Money, method PaymentMethod, reference, notes string, paidAt, createdAt, updatedAt time.Time) (Payment, error) {
	if _, err := ParsePaymentID(string(id)); err != nil {
		return Payment{}, err
	}
	if _, err := ParseGymID(string(gymID)); err != nil {
		return Payment{}, err
	}
	if _, err := ParseMemberID(string(memberID)); err != nil {
		return Payment{}, err
	}
	if membershipID != "" {
		if _, err := ParseMembershipID(string(membershipID)); err != nil {
			return Payment{}, err
		}
	}
	status, err := ParsePaymentStatus(string(status))
	if err != nil {
		return Payment{}, err
	}
	kind, err = ParsePaymentKind(string(kind))
	if err != nil {
		return Payment{}, err
	}
	method, err = ParsePaymentMethod(string(method))
	if err != nil {
		return Payment{}, err
	}
	if err := validatePositiveMoney(amount, ErrPaymentAmountMustBePositive); err != nil {
		return Payment{}, err
	}
	createdAt, updatedAt, err = normalizeLifecycleTimestamps(createdAt, updatedAt)
	if err != nil {
		return Payment{}, err
	}
	paidAt, err = validatePaymentLifecycle(status, paidAt, createdAt, updatedAt)
	if err != nil {
		return Payment{}, err
	}
	return Payment{id: id, gymID: gymID, memberID: memberID, membershipID: membershipID, status: status, kind: kind, amount: amount, method: method, reference: strings.TrimSpace(reference), notes: strings.TrimSpace(notes), paidAt: paidAt, createdAt: createdAt, updatedAt: updatedAt}, nil
}

func (p Payment) ID() PaymentID              { return p.id }
func (p Payment) GymID() GymID               { return p.gymID }
func (p Payment) MemberID() MemberID         { return p.memberID }
func (p Payment) MembershipID() MembershipID { return p.membershipID }
func (p Payment) Status() PaymentStatus      { return p.status }
func (p Payment) Kind() PaymentKind          { return p.kind }
func (p Payment) Amount() Money              { return p.amount }
func (p Payment) Method() PaymentMethod      { return p.method }
func (p Payment) Reference() string          { return p.reference }
func (p Payment) Notes() string              { return p.notes }
func (p Payment) PaidAt() time.Time          { return p.paidAt }
func (p Payment) CreatedAt() time.Time       { return p.createdAt }
func (p Payment) UpdatedAt() time.Time       { return p.updatedAt }

func (p Payment) Post(at time.Time) (Payment, error) {
	if p.status != PaymentStatusPending {
		return Payment{}, ErrPaymentCannotPost
	}
	at, err := paymentTransitionTime(at, p.updatedAt)
	if err != nil {
		return Payment{}, err
	}
	p.status = PaymentStatusPosted
	p.paidAt = at
	p.updatedAt = at
	return p, nil
}

func (p Payment) Void(at time.Time) (Payment, error) {
	if p.status != PaymentStatusPending {
		return Payment{}, ErrPaymentCannotVoid
	}
	at, err := paymentTransitionTime(at, p.updatedAt)
	if err != nil {
		return Payment{}, err
	}
	p.status = PaymentStatusVoided
	p.updatedAt = at
	return p, nil
}

func (p Payment) Refund(at time.Time) (Payment, error) {
	if p.status != PaymentStatusPosted {
		return Payment{}, ErrPaymentCannotRefund
	}
	at, err := paymentTransitionTime(at, p.updatedAt)
	if err != nil {
		return Payment{}, err
	}
	p.status = PaymentStatusRefunded
	p.updatedAt = at
	return p, nil
}

func validatePositiveMoney(amount Money, zeroError error) error {
	if _, err := NewMoney(amount.Cents(), amount.Currency()); err != nil {
		return err
	}
	if amount.Cents() == 0 {
		return zeroError
	}
	return nil
}

func validatePaymentLifecycle(status PaymentStatus, paidAt, createdAt, updatedAt time.Time) (time.Time, error) {
	switch status {
	case PaymentStatusPending, PaymentStatusVoided:
		if !paidAt.IsZero() {
			return time.Time{}, fmt.Errorf("%w: pending or voided payment cannot have paid_at", ErrInvalidPaymentLifecycle)
		}
		return time.Time{}, nil
	case PaymentStatusPosted, PaymentStatusRefunded:
		if paidAt.IsZero() {
			return time.Time{}, fmt.Errorf("%w: posted or refunded payment requires paid_at", ErrInvalidPaymentLifecycle)
		}
		paidAt = paidAt.UTC()
		if paidAt.Before(createdAt) || paidAt.After(updatedAt) {
			return time.Time{}, fmt.Errorf("%w: paid_at must be within the payment lifecycle", ErrInvalidPaymentLifecycle)
		}
		return paidAt, nil
	default:
		return time.Time{}, ErrInvalidPaymentLifecycle
	}
}

func paymentTransitionTime(at, updatedAt time.Time) (time.Time, error) {
	if at.IsZero() {
		return time.Time{}, ErrInvalidPaymentLifecycle
	}
	at = at.UTC()
	if at.Before(updatedAt) {
		return time.Time{}, fmt.Errorf("%w: transition precedes latest payment update", ErrInvalidPaymentLifecycle)
	}
	return at, nil
}
