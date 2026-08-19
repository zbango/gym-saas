package domain

import (
	"errors"
	"math"
)

var (
	ErrNegativeAmount   = errors.New("money amount cannot be negative")
	ErrInvalidCurrency  = errors.New("currency must be a three-letter uppercase currency code")
	ErrCurrencyMismatch = errors.New("money currencies must match")
	ErrAmountOverflow   = errors.New("money amount overflow")
)

// Money represents a non-negative amount in integer minor units.
//
// Domain workflows decide whether an amount is a payment, an expense, or a
// refund; they never represent those concepts with floating point values.
type Money struct {
	cents    int64
	currency string
}

func NewMoney(cents int64, currency string) (Money, error) {
	money := Money{cents: cents, currency: currency}
	if err := money.validate(); err != nil {
		return Money{}, err
	}

	return money, nil
}

func (m Money) Cents() int64 {
	return m.cents
}

func (m Money) Currency() string {
	return m.currency
}

func (m Money) Add(other Money) (Money, error) {
	if err := m.validate(); err != nil {
		return Money{}, err
	}
	if err := other.validate(); err != nil {
		return Money{}, err
	}
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}
	if m.cents > math.MaxInt64-other.cents {
		return Money{}, ErrAmountOverflow
	}

	return NewMoney(m.cents+other.cents, m.currency)
}

// Compare returns -1, 0, or 1 for smaller, equal, or larger values.
func (m Money) Compare(other Money) (int, error) {
	if err := m.validate(); err != nil {
		return 0, err
	}
	if err := other.validate(); err != nil {
		return 0, err
	}
	if m.currency != other.currency {
		return 0, ErrCurrencyMismatch
	}
	switch {
	case m.cents < other.cents:
		return -1, nil
	case m.cents > other.cents:
		return 1, nil
	default:
		return 0, nil
	}
}

func (m Money) validate() error {
	if m.cents < 0 {
		return ErrNegativeAmount
	}
	if !isCurrency(m.currency) {
		return ErrInvalidCurrency
	}
	return nil
}

func isCurrency(value string) bool {
	if len(value) != 3 {
		return false
	}
	for _, character := range value {
		if character < 'A' || character > 'Z' {
			return false
		}
	}
	return true
}
