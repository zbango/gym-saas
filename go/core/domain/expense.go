package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidExpenseID            = errors.New("invalid expense ID")
	ErrInvalidExpenseStatus        = errors.New("invalid expense status")
	ErrExpenseAmountMustBePositive = errors.New("expense amount must be greater than zero")
	ErrInvalidExpenseLifecycle     = errors.New("invalid expense lifecycle")
	ErrExpenseCannotPost           = errors.New("expense cannot be posted in its current status")
	ErrExpenseCannotVoid           = errors.New("expense cannot be voided in its current status")
)

type ExpenseID string

func NewExpenseID() (ExpenseID, error) {
	value, err := NewUUID()
	if err != nil {
		return "", err
	}
	return ExpenseID(value), nil
}

func ParseExpenseID(value string) (ExpenseID, error) {
	if err := ValidateUUID(value); err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidExpenseID, err)
	}
	return ExpenseID(value), nil
}

type ExpenseStatus string

const (
	ExpenseStatusPending ExpenseStatus = "pending"
	ExpenseStatusPosted  ExpenseStatus = "posted"
	ExpenseStatusVoided  ExpenseStatus = "voided"
)

func ParseExpenseStatus(value string) (ExpenseStatus, error) {
	switch status := ExpenseStatus(value); status {
	case ExpenseStatusPending, ExpenseStatusPosted, ExpenseStatusVoided:
		return status, nil
	default:
		return "", fmt.Errorf("%w: %q", ErrInvalidExpenseStatus, value)
	}
}

// Expense is an outgoing business payment. It intentionally does not refer to
// a member or membership, keeping business expenses separate from receipts.
type Expense struct {
	id        ExpenseID
	gymID     GymID
	status    ExpenseStatus
	amount    Money
	method    PaymentMethod
	reference string
	notes     string
	paidAt    time.Time
	createdAt time.Time
	updatedAt time.Time
}

func RecordExpense(gymID GymID, amount Money, method PaymentMethod, reference, notes string, paidAt time.Time) (Expense, error) {
	id, err := NewExpenseID()
	if err != nil {
		return Expense{}, fmt.Errorf("create expense ID: %w", err)
	}
	return NewExpense(id, gymID, ExpenseStatusPosted, amount, method, reference, notes, paidAt, paidAt, paidAt)
}

func CreatePendingExpense(gymID GymID, amount Money, method PaymentMethod, reference, notes string, createdAt time.Time) (Expense, error) {
	id, err := NewExpenseID()
	if err != nil {
		return Expense{}, fmt.Errorf("create expense ID: %w", err)
	}
	return NewExpense(id, gymID, ExpenseStatusPending, amount, method, reference, notes, time.Time{}, createdAt, createdAt)
}

func NewExpense(id ExpenseID, gymID GymID, status ExpenseStatus, amount Money, method PaymentMethod, reference, notes string, paidAt, createdAt, updatedAt time.Time) (Expense, error) {
	if _, err := ParseExpenseID(string(id)); err != nil {
		return Expense{}, err
	}
	if _, err := ParseGymID(string(gymID)); err != nil {
		return Expense{}, err
	}
	status, err := ParseExpenseStatus(string(status))
	if err != nil {
		return Expense{}, err
	}
	method, err = ParsePaymentMethod(string(method))
	if err != nil {
		return Expense{}, err
	}
	if err := validatePositiveMoney(amount, ErrExpenseAmountMustBePositive); err != nil {
		return Expense{}, err
	}
	createdAt, updatedAt, err = normalizeLifecycleTimestamps(createdAt, updatedAt)
	if err != nil {
		return Expense{}, err
	}
	paidAt, err = validateExpenseLifecycle(status, paidAt, createdAt, updatedAt)
	if err != nil {
		return Expense{}, err
	}
	return Expense{id: id, gymID: gymID, status: status, amount: amount, method: method, reference: strings.TrimSpace(reference), notes: strings.TrimSpace(notes), paidAt: paidAt, createdAt: createdAt, updatedAt: updatedAt}, nil
}

func (e Expense) ID() ExpenseID         { return e.id }
func (e Expense) GymID() GymID          { return e.gymID }
func (e Expense) Status() ExpenseStatus { return e.status }
func (e Expense) Amount() Money         { return e.amount }
func (e Expense) Method() PaymentMethod { return e.method }
func (e Expense) Reference() string     { return e.reference }
func (e Expense) Notes() string         { return e.notes }
func (e Expense) PaidAt() time.Time     { return e.paidAt }
func (e Expense) CreatedAt() time.Time  { return e.createdAt }
func (e Expense) UpdatedAt() time.Time  { return e.updatedAt }

func (e Expense) Post(at time.Time) (Expense, error) {
	if e.status != ExpenseStatusPending {
		return Expense{}, ErrExpenseCannotPost
	}
	at, err := expenseTransitionTime(at, e.updatedAt)
	if err != nil {
		return Expense{}, err
	}
	e.status = ExpenseStatusPosted
	e.paidAt = at
	e.updatedAt = at
	return e, nil
}

func (e Expense) Void(at time.Time) (Expense, error) {
	if e.status != ExpenseStatusPending {
		return Expense{}, ErrExpenseCannotVoid
	}
	at, err := expenseTransitionTime(at, e.updatedAt)
	if err != nil {
		return Expense{}, err
	}
	e.status = ExpenseStatusVoided
	e.updatedAt = at
	return e, nil
}

func validateExpenseLifecycle(status ExpenseStatus, paidAt, createdAt, updatedAt time.Time) (time.Time, error) {
	switch status {
	case ExpenseStatusPending, ExpenseStatusVoided:
		if !paidAt.IsZero() {
			return time.Time{}, fmt.Errorf("%w: pending or voided expense cannot have paid_at", ErrInvalidExpenseLifecycle)
		}
		return time.Time{}, nil
	case ExpenseStatusPosted:
		if paidAt.IsZero() {
			return time.Time{}, fmt.Errorf("%w: posted expense requires paid_at", ErrInvalidExpenseLifecycle)
		}
		paidAt = paidAt.UTC()
		if paidAt.Before(createdAt) || paidAt.After(updatedAt) {
			return time.Time{}, fmt.Errorf("%w: paid_at must be within the expense lifecycle", ErrInvalidExpenseLifecycle)
		}
		return paidAt, nil
	default:
		return time.Time{}, ErrInvalidExpenseLifecycle
	}
}

func expenseTransitionTime(at, updatedAt time.Time) (time.Time, error) {
	if at.IsZero() {
		return time.Time{}, ErrInvalidExpenseLifecycle
	}
	at = at.UTC()
	if at.Before(updatedAt) {
		return time.Time{}, fmt.Errorf("%w: transition precedes latest expense update", ErrInvalidExpenseLifecycle)
	}
	return at, nil
}
