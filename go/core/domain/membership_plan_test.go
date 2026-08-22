package domain

import (
	"errors"
	"testing"
	"time"
)

func TestCreateMembershipPlan(t *testing.T) {
	now := time.Date(2026, time.August, 20, 10, 0, 0, 0, time.UTC)
	price, err := NewMoney(4500, "USD")
	if err != nil {
		t.Fatalf("NewMoney returned error: %v", err)
	}
	plan, err := CreateMembershipPlan(mustGymID(t), "  Monthly  ", MembershipValidityTime, 1, MembershipDurationMonths, 0, price, MembershipPlanStatusActive, now)
	if err != nil {
		t.Fatalf("CreateMembershipPlan returned error: %v", err)
	}
	if err := ValidateUUID(string(plan.ID())); err != nil {
		t.Fatalf("plan ID is invalid: %v", err)
	}
	if plan.Name() != "Monthly" || plan.DurationValue() != 1 || plan.DurationUnit() != MembershipDurationMonths || plan.VisitLimit() != 0 || plan.Price().Cents() != 4500 {
		t.Fatalf("plan = %#v", plan)
	}
}

func TestMembershipPlanRequiresDurationAndVisitLimitForVisitPlans(t *testing.T) {
	price, err := NewMoney(6000, "USD")
	if err != nil {
		t.Fatalf("NewMoney returned error: %v", err)
	}
	now := time.Date(2026, time.August, 20, 10, 0, 0, 0, time.UTC)
	plan, err := CreateMembershipPlan(mustGymID(t), "10 visits / 2 months", MembershipValidityVisits, 2, MembershipDurationMonths, 10, price, MembershipPlanStatusActive, now)
	if err != nil {
		t.Fatalf("CreateMembershipPlan returned error: %v", err)
	}
	if plan.ValidityKind() != MembershipValidityVisits || plan.VisitLimit() != 10 || plan.DurationValue() != 2 {
		t.Fatalf("visit plan = %#v", plan)
	}
}

func TestMembershipPlanRejectsInvalidConfiguration(t *testing.T) {
	price, err := NewMoney(0, "USD")
	if err != nil {
		t.Fatalf("NewMoney returned error: %v", err)
	}
	id := mustMembershipPlanID(t)
	gymID := mustGymID(t)
	now := time.Date(2026, time.August, 20, 10, 0, 0, 0, time.UTC)
	tests := []struct {
		name          string
		kind          MembershipValidityKind
		durationValue int
		durationUnit  MembershipDurationUnit
		visitLimit    int
		wantErr       error
	}{
		{name: "blank name", kind: MembershipValidityTime, durationValue: 1, durationUnit: MembershipDurationMonths, wantErr: ErrInvalidMembershipPlanName},
		{name: "zero duration", kind: MembershipValidityTime, durationUnit: MembershipDurationMonths, wantErr: ErrInvalidDurationValue},
		{name: "bad duration unit", kind: MembershipValidityTime, durationValue: 1, durationUnit: "fortnights", wantErr: ErrInvalidDurationUnit},
		{name: "time plan has visits", kind: MembershipValidityTime, durationValue: 1, durationUnit: MembershipDurationMonths, visitLimit: 1, wantErr: ErrInvalidVisitLimit},
		{name: "visit plan has no visits", kind: MembershipValidityVisits, durationValue: 2, durationUnit: MembershipDurationMonths, wantErr: ErrInvalidVisitLimit},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			name := "Plan"
			if test.name == "blank name" {
				name = " "
			}
			_, err := NewMembershipPlan(id, gymID, name, test.kind, test.durationValue, test.durationUnit, test.visitLimit, price, MembershipPlanStatusActive, now, now)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("NewMembershipPlan error = %v, want %v", err, test.wantErr)
			}
		})
	}
}

func mustMembershipPlanID(t *testing.T) MembershipPlanID {
	t.Helper()
	id, err := NewMembershipPlanID()
	if err != nil {
		t.Fatalf("NewMembershipPlanID returned error: %v", err)
	}
	return id
}
