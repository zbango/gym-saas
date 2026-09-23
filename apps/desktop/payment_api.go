package main

import (
	"github.com/zbango/gym-saas/go/core/application"
	"github.com/zbango/gym-saas/go/core/domain"
)

// PaymentAPI is the Wails delivery adapter for standalone receipt lifecycle
// use cases. Membership purchases use MembershipPurchaseAPI instead.
type PaymentAPI struct {
	runtime  *DesktopRuntime
	payments *application.PaymentService
}

func NewPaymentAPI(runtime *DesktopRuntime, payments *application.PaymentService) *PaymentAPI {
	return &PaymentAPI{runtime: runtime, payments: payments}
}

type PaymentInput struct {
	MemberID      string `json:"memberId"`
	MembershipID  string `json:"membershipId"`
	Kind          string `json:"kind"`
	AmountCents   int64  `json:"amountCents"`
	Currency      string `json:"currency"`
	PaymentMethod string `json:"paymentMethod"`
	Reference     string `json:"reference"`
	Notes         string `json:"notes"`
}

type Payment struct {
	ID            string `json:"id"`
	MemberID      string `json:"memberId"`
	MembershipID  string `json:"membershipId"`
	Status        string `json:"status"`
	Kind          string `json:"kind"`
	AmountCents   int64  `json:"amountCents"`
	Currency      string `json:"currency"`
	PaymentMethod string `json:"paymentMethod"`
	Reference     string `json:"reference"`
	Notes         string `json:"notes"`
	PaidAt        string `json:"paidAt"`
}

func (a *PaymentAPI) ListPaymentsForMember(memberID string) ([]Payment, error) {
	payments, err := a.payments.ListForMember(a.runtime.requestContext(), memberID)
	if err != nil {
		return nil, err
	}
	result := make([]Payment, 0, len(payments))
	for _, payment := range payments {
		result = append(result, paymentView(payment))
	}
	return result, nil
}

func (a *PaymentAPI) RecordPayment(input PaymentInput) (Payment, error) {
	payment, err := a.payments.Record(a.runtime.requestContext(), application.RecordPaymentInput(input))
	if err != nil {
		return Payment{}, err
	}
	return paymentView(payment), nil
}

func (a *PaymentAPI) CreatePendingPayment(input PaymentInput) (Payment, error) {
	payment, err := a.payments.CreatePending(a.runtime.requestContext(), application.RecordPaymentInput(input))
	if err != nil {
		return Payment{}, err
	}
	return paymentView(payment), nil
}

func (a *PaymentAPI) PostPayment(id string) (Payment, error) {
	payment, err := a.payments.Post(a.runtime.requestContext(), id)
	if err != nil {
		return Payment{}, err
	}
	return paymentView(payment), nil
}

func (a *PaymentAPI) VoidPayment(id string) (Payment, error) {
	payment, err := a.payments.Void(a.runtime.requestContext(), id)
	if err != nil {
		return Payment{}, err
	}
	return paymentView(payment), nil
}

func (a *PaymentAPI) RefundPayment(id string) (Payment, error) {
	payment, err := a.payments.Refund(a.runtime.requestContext(), id)
	if err != nil {
		return Payment{}, err
	}
	return paymentView(payment), nil
}

func paymentView(payment domain.Payment) Payment {
	return Payment{
		ID: string(payment.ID()), MemberID: string(payment.MemberID()), MembershipID: string(payment.MembershipID()),
		Status: string(payment.Status()), Kind: string(payment.Kind()), AmountCents: payment.Amount().Cents(),
		Currency: payment.Amount().Currency(), PaymentMethod: string(payment.Method()), Reference: payment.Reference(),
		Notes: payment.Notes(), PaidAt: formatOptionalTimestamp(payment.PaidAt()),
	}
}
