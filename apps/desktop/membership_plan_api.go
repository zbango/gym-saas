package main

import (
	"github.com/zbango/gym-saas/go/core/application"
	"github.com/zbango/gym-saas/go/core/domain"
)

// MembershipPlanAPI is the Wails delivery adapter for Membership Plan use cases.
type MembershipPlanAPI struct {
	runtime *DesktopRuntime
	plans   *application.MembershipPlanService
}

func NewMembershipPlanAPI(runtime *DesktopRuntime, plans *application.MembershipPlanService) *MembershipPlanAPI {
	return &MembershipPlanAPI{runtime: runtime, plans: plans}
}

type MembershipPlanInput struct {
	Name          string `json:"name"`
	ValidityKind  string `json:"validityKind"`
	DurationValue int    `json:"durationValue"`
	DurationUnit  string `json:"durationUnit"`
	VisitLimit    int    `json:"visitLimit"`
	PriceCents    int64  `json:"priceCents"`
	Currency      string `json:"currency"`
	Status        string `json:"status"`
}

type MembershipPlan struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	ValidityKind  string `json:"validityKind"`
	DurationValue int    `json:"durationValue"`
	DurationUnit  string `json:"durationUnit"`
	VisitLimit    int    `json:"visitLimit"`
	PriceCents    int64  `json:"priceCents"`
	Currency      string `json:"currency"`
	Status        string `json:"status"`
}

func (a *MembershipPlanAPI) ListMembershipPlans() ([]MembershipPlan, error) {
	plans, err := a.plans.List(a.runtime.requestContext())
	if err != nil {
		return nil, err
	}
	result := make([]MembershipPlan, 0, len(plans))
	for _, plan := range plans {
		result = append(result, membershipPlanView(plan))
	}
	return result, nil
}

func (a *MembershipPlanAPI) CreateMembershipPlan(input MembershipPlanInput) (MembershipPlan, error) {
	plan, err := a.plans.Create(a.runtime.requestContext(), application.MembershipPlanInput(input))
	if err != nil {
		return MembershipPlan{}, err
	}
	return membershipPlanView(plan), nil
}

func (a *MembershipPlanAPI) UpdateMembershipPlan(id string, input MembershipPlanInput) (MembershipPlan, error) {
	plan, err := a.plans.Update(a.runtime.requestContext(), id, application.MembershipPlanInput(input))
	if err != nil {
		return MembershipPlan{}, err
	}
	return membershipPlanView(plan), nil
}

func (a *MembershipPlanAPI) ArchiveMembershipPlan(id string) error {
	return a.plans.Archive(a.runtime.requestContext(), id)
}

func membershipPlanView(plan domain.MembershipPlan) MembershipPlan {
	return MembershipPlan{
		ID:            string(plan.ID()),
		Name:          plan.Name(),
		ValidityKind:  string(plan.ValidityKind()),
		DurationValue: plan.DurationValue(),
		DurationUnit:  string(plan.DurationUnit()),
		VisitLimit:    plan.VisitLimit(),
		PriceCents:    plan.Price().Cents(),
		Currency:      plan.Price().Currency(),
		Status:        string(plan.Status()),
	}
}
