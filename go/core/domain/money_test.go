package domain

import (
	"errors"
	"math"
	"testing"
)

func TestNewMoney(t *testing.T) {
	tests := []struct {
		name     string
		cents    int64
		currency string
		wantErr  error
	}{
		{name: "zero", cents: 0, currency: "USD"},
		{name: "positive", cents: 12345, currency: "USD"},
		{name: "negative", cents: -1, currency: "USD", wantErr: ErrNegativeAmount},
		{name: "lowercase currency", cents: 1, currency: "usd", wantErr: ErrInvalidCurrency},
		{name: "missing currency", cents: 1, wantErr: ErrInvalidCurrency},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			money, err := NewMoney(test.cents, test.currency)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("NewMoney() error = %v, want %v", err, test.wantErr)
			}
			if err == nil && (money.Cents() != test.cents || money.Currency() != test.currency) {
				t.Fatalf("NewMoney() = %#v, want %d %s", money, test.cents, test.currency)
			}
		})
	}
}

func TestMoneyAdd(t *testing.T) {
	left := mustMoney(t, 1250, "USD")
	right := mustMoney(t, 350, "USD")

	total, err := left.Add(right)
	if err != nil {
		t.Fatalf("Add returned error: %v", err)
	}
	if total.Cents() != 1600 || total.Currency() != "USD" {
		t.Fatalf("Add() = %#v, want 1600 USD", total)
	}
}

func TestMoneyAddRejectsCurrencyMismatchAndOverflow(t *testing.T) {
	usd := mustMoney(t, 1, "USD")
	eur := mustMoney(t, 1, "EUR")
	if _, err := usd.Add(eur); !errors.Is(err, ErrCurrencyMismatch) {
		t.Fatalf("Add() error = %v, want %v", err, ErrCurrencyMismatch)
	}

	maximum := mustMoney(t, math.MaxInt64, "USD")
	if _, err := maximum.Add(usd); !errors.Is(err, ErrAmountOverflow) {
		t.Fatalf("Add() error = %v, want %v", err, ErrAmountOverflow)
	}
}

func TestMoneyCompare(t *testing.T) {
	low := mustMoney(t, 100, "USD")
	high := mustMoney(t, 200, "USD")

	for _, test := range []struct {
		left, right Money
		want        int
	}{
		{left: low, right: high, want: -1},
		{left: high, right: low, want: 1},
		{left: low, right: low, want: 0},
	} {
		got, err := test.left.Compare(test.right)
		if err != nil {
			t.Fatalf("Compare returned error: %v", err)
		}
		if got != test.want {
			t.Fatalf("Compare() = %d, want %d", got, test.want)
		}
	}

	if _, err := low.Compare(mustMoney(t, 100, "EUR")); !errors.Is(err, ErrCurrencyMismatch) {
		t.Fatalf("Compare() error = %v, want %v", err, ErrCurrencyMismatch)
	}
}

func TestMoneyOperationsRejectZeroValue(t *testing.T) {
	valid := mustMoney(t, 1, "USD")
	if _, err := (Money{}).Add(valid); !errors.Is(err, ErrInvalidCurrency) {
		t.Fatalf("Add() error = %v, want %v", err, ErrInvalidCurrency)
	}
	if _, err := (Money{}).Compare(valid); !errors.Is(err, ErrInvalidCurrency) {
		t.Fatalf("Compare() error = %v, want %v", err, ErrInvalidCurrency)
	}
}

func mustMoney(t *testing.T, cents int64, currency string) Money {
	t.Helper()
	money, err := NewMoney(cents, currency)
	if err != nil {
		t.Fatalf("NewMoney(%d, %q) returned error: %v", cents, currency, err)
	}
	return money
}
