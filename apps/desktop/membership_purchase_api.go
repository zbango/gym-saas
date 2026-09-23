package main

import (
	"strings"
	"time"

	"github.com/zbango/gym-saas/go/core/application"
	"github.com/zbango/gym-saas/go/core/domain"
)

// MembershipPurchaseAPI is the Wails adapter for the atomic local membership
// purchase workflow.
type MembershipPurchaseAPI struct {
	runtime   *DesktopRuntime
	purchases *application.MembershipPurchaseService
}

func NewMembershipPurchaseAPI(runtime *DesktopRuntime, purchases *application.MembershipPurchaseService) *MembershipPurchaseAPI {
	return &MembershipPurchaseAPI{runtime: runtime, purchases: purchases}
}

type PurchaseMembershipInput struct {
	MemberID         string `json:"memberId"`
	MembershipPlanID string `json:"membershipPlanId"`
	StartsAt         string `json:"startsAt"`
	AmountCents      int64  `json:"amountCents"`
	Currency         string `json:"currency"`
	PaymentMethod    string `json:"paymentMethod"`
	Reference        string `json:"reference"`
	Notes            string `json:"notes"`
}

type MembershipPurchase struct {
	Membership Membership `json:"membership"`
	Payment    Payment    `json:"payment"`
}

func (a *MembershipPurchaseAPI) PurchaseMembership(input PurchaseMembershipInput) (MembershipPurchase, error) {
	startsAt, err := parseOptionalTimestamp(input.StartsAt)
	if err != nil {
		return MembershipPurchase{}, err
	}
	result, err := a.purchases.Purchase(a.runtime.requestContext(), application.PurchaseMembershipInput{
		MemberID: input.MemberID, MembershipPlanID: input.MembershipPlanID, StartsAt: startsAt,
		AmountCents: input.AmountCents, Currency: input.Currency, PaymentMethod: input.PaymentMethod,
		Reference: input.Reference, Notes: input.Notes,
	})
	if err != nil {
		return MembershipPurchase{}, err
	}
	return MembershipPurchase{Membership: membershipView(result.Membership), Payment: paymentView(result.Payment)}, nil
}

type Membership struct {
	ID                 string `json:"id"`
	MemberID           string `json:"memberId"`
	MembershipPlanID   string `json:"membershipPlanId"`
	Status             string `json:"status"`
	StartsAt           string `json:"startsAt"`
	EndsAt             string `json:"endsAt"`
	ActivatedAt        string `json:"activatedAt"`
	CancelledAt        string `json:"cancelledAt"`
	CancellationReason string `json:"cancellationReason"`
	ValidityKind       string `json:"validityKind"`
	DurationValue      int    `json:"durationValue"`
	DurationUnit       string `json:"durationUnit"`
	VisitLimit         int    `json:"visitLimit"`
	VisitsRemaining    int    `json:"visitsRemaining"`
	PriceCents         int64  `json:"priceCents"`
	Currency           string `json:"currency"`
}

func membershipView(membership domain.Membership) Membership {
	return Membership{
		ID: string(membership.ID()), MemberID: string(membership.MemberID()), MembershipPlanID: string(membership.MembershipPlanID()),
		Status: string(membership.Status()), StartsAt: domain.FormatTimestamp(membership.StartsAt()), EndsAt: domain.FormatTimestamp(membership.EndsAt()),
		ActivatedAt: formatOptionalTimestamp(membership.ActivatedAt()), CancelledAt: formatOptionalTimestamp(membership.CancelledAt()),
		CancellationReason: membership.CancellationReason(), ValidityKind: string(membership.ValidityKind()), DurationValue: membership.DurationValue(),
		DurationUnit: string(membership.DurationUnit()), VisitLimit: membership.VisitLimit(), VisitsRemaining: membership.VisitsRemaining(),
		PriceCents: membership.Price().Cents(), Currency: membership.Price().Currency(),
	}
}

func parseOptionalTimestamp(value string) (timeValue time.Time, err error) {
	if strings.TrimSpace(value) == "" {
		return timeValue, nil
	}
	return domain.ParseTimestamp(value)
}
