package domain

import (
	"errors"
	"testing"
	"time"
)

func TestStartMembershipSnapshotsPlanAndUsesGymCalendar(t *testing.T) {
	gym := mustGym(t, "America/Guayaquil")
	plan := mustPlan(t, gym, MembershipValidityTime, 1, MembershipDurationMonths, 0)
	memberID := mustMemberID(t)
	startsAt := time.Date(2026, time.January, 31, 10, 0, 0, 0, time.FixedZone("ECT", -5*60*60))

	membership, err := StartMembership(gym, memberID, plan, startsAt, startsAt)
	if err != nil {
		t.Fatalf("StartMembership returned error: %v", err)
	}
	if membership.Status() != MembershipStatusActive || !membership.ActivatedAt().Equal(startsAt.UTC()) {
		t.Fatalf("membership activation = %q at %s", membership.Status(), membership.ActivatedAt())
	}
	if got, want := membership.EndsAt(), time.Date(2026, time.February, 28, 15, 0, 0, 0, time.UTC); !got.Equal(want) {
		t.Fatalf("end = %s, want %s", got, want)
	}
	if membership.Price().Cents() != plan.Price().Cents() || membership.DurationValue() != 1 || membership.DurationUnit() != MembershipDurationMonths {
		t.Fatalf("membership did not snapshot plan terms: %#v", membership)
	}
	if !membership.IsValidAt(membership.EndsAt().Add(-time.Nanosecond)) || membership.IsValidAt(membership.EndsAt()) {
		t.Fatal("membership end boundary must be exclusive")
	}
}

func TestVisitMembershipExpiresWhenVisitsAreExhausted(t *testing.T) {
	gym := mustGym(t, "UTC")
	plan := mustPlan(t, gym, MembershipValidityVisits, 2, MembershipDurationMonths, 2)
	startsAt := time.Date(2026, time.August, 1, 9, 0, 0, 0, time.UTC)
	membership, err := StartMembership(gym, mustMemberID(t), plan, startsAt, startsAt)
	if err != nil {
		t.Fatalf("StartMembership returned error: %v", err)
	}
	if membership.VisitsRemaining() != 2 || !membership.IsValidAt(startsAt) {
		t.Fatalf("new visit membership = %#v", membership)
	}

	membership, err = membership.ConsumeVisit(startsAt.Add(time.Hour))
	if err != nil || membership.VisitsRemaining() != 1 || !membership.IsValidAt(startsAt.Add(time.Hour)) {
		t.Fatalf("first ConsumeVisit = %#v, %v", membership, err)
	}
	membership, err = membership.ConsumeVisit(startsAt.Add(2 * time.Hour))
	if err != nil {
		t.Fatalf("second ConsumeVisit returned error: %v", err)
	}
	if membership.VisitsRemaining() != 0 || membership.IsValidAt(startsAt.Add(2*time.Hour)) {
		t.Fatalf("exhausted visit membership = %#v", membership)
	}
	if _, err := membership.ConsumeVisit(startsAt.Add(3 * time.Hour)); !errors.Is(err, ErrMembershipNotValid) {
		t.Fatalf("ConsumeVisit after exhaustion error = %v, want %v", err, ErrMembershipNotValid)
	}
}

func TestFutureMembershipRequiresActivationAtStart(t *testing.T) {
	gym := mustGym(t, "UTC")
	plan := mustPlan(t, gym, MembershipValidityTime, 1, MembershipDurationMonths, 0)
	createdAt := time.Date(2026, time.August, 1, 9, 0, 0, 0, time.UTC)
	startsAt := createdAt.AddDate(0, 0, 7)
	membership, err := StartMembership(gym, mustMemberID(t), plan, startsAt, createdAt)
	if err != nil {
		t.Fatalf("StartMembership returned error: %v", err)
	}
	if membership.Status() != MembershipStatusPending || membership.IsValidAt(startsAt) {
		t.Fatalf("future membership = %#v", membership)
	}
	if _, err := membership.Activate(startsAt.Add(-time.Nanosecond)); !errors.Is(err, ErrMembershipActivationTooEarly) {
		t.Fatalf("early Activate error = %v, want %v", err, ErrMembershipActivationTooEarly)
	}
	activated, err := membership.Activate(startsAt)
	if err != nil {
		t.Fatalf("Activate returned error: %v", err)
	}
	if activated.Status() != MembershipStatusActive || !activated.IsValidAt(startsAt) {
		t.Fatalf("activated membership = %#v", activated)
	}
}

func TestMembershipCancellationRequiresReasonAndCannotCancelExpiredMembership(t *testing.T) {
	gym := mustGym(t, "UTC")
	plan := mustPlan(t, gym, MembershipValidityTime, 1, MembershipDurationDays, 0)
	startsAt := time.Date(2026, time.August, 1, 9, 0, 0, 0, time.UTC)
	membership, err := StartMembership(gym, mustMemberID(t), plan, startsAt, startsAt)
	if err != nil {
		t.Fatalf("StartMembership returned error: %v", err)
	}
	if _, err := membership.Cancel(startsAt.Add(time.Hour), " "); !errors.Is(err, ErrInvalidCancellationReason) {
		t.Fatalf("Cancel without reason error = %v, want %v", err, ErrInvalidCancellationReason)
	}
	if _, err := membership.Cancel(membership.EndsAt(), "member request"); !errors.Is(err, ErrMembershipAlreadyExpired) {
		t.Fatalf("Cancel expired membership error = %v, want %v", err, ErrMembershipAlreadyExpired)
	}
	cancelled, err := membership.Cancel(startsAt.Add(time.Hour), "  member request  ")
	if err != nil {
		t.Fatalf("Cancel returned error: %v", err)
	}
	if cancelled.Status() != MembershipStatusCancelled || cancelled.CancellationReason() != "member request" || cancelled.IsValidAt(startsAt.Add(2*time.Hour)) {
		t.Fatalf("cancelled membership = %#v", cancelled)
	}
}

func TestMembershipMarksExpiredOnlyAfterAValidityBound(t *testing.T) {
	gym := mustGym(t, "UTC")
	plan := mustPlan(t, gym, MembershipValidityTime, 1, MembershipDurationDays, 0)
	startsAt := time.Date(2026, time.August, 1, 9, 0, 0, 0, time.UTC)
	membership, err := StartMembership(gym, mustMemberID(t), plan, startsAt, startsAt)
	if err != nil {
		t.Fatalf("StartMembership returned error: %v", err)
	}
	if _, err := membership.MarkExpired(startsAt.Add(time.Hour)); !errors.Is(err, ErrMembershipNotValid) {
		t.Fatalf("early MarkExpired error = %v, want %v", err, ErrMembershipNotValid)
	}
	expired, err := membership.MarkExpired(membership.EndsAt())
	if err != nil {
		t.Fatalf("MarkExpired returned error: %v", err)
	}
	if expired.Status() != MembershipStatusExpired || expired.IsValidAt(expired.EndsAt()) {
		t.Fatalf("expired membership = %#v", expired)
	}
}

func TestStartMembershipRejectsInactiveOrForeignPlan(t *testing.T) {
	gym := mustGym(t, "UTC")
	otherGym := mustGym(t, "UTC")
	plan := mustPlan(t, otherGym, MembershipValidityTime, 1, MembershipDurationMonths, 0)
	now := time.Date(2026, time.August, 1, 9, 0, 0, 0, time.UTC)
	if _, err := StartMembership(gym, mustMemberID(t), plan, now, now); !errors.Is(err, ErrMembershipPlanGymMismatch) {
		t.Fatalf("foreign plan error = %v, want %v", err, ErrMembershipPlanGymMismatch)
	}

	price, err := NewMoney(4500, "USD")
	if err != nil {
		t.Fatalf("NewMoney returned error: %v", err)
	}
	inactivePlan, err := CreateMembershipPlan(gym.ID(), "Inactive", MembershipValidityTime, 1, MembershipDurationMonths, 0, price, MembershipPlanStatusInactive, now)
	if err != nil {
		t.Fatalf("CreateMembershipPlan returned error: %v", err)
	}
	if _, err := StartMembership(gym, mustMemberID(t), inactivePlan, now, now); !errors.Is(err, ErrMembershipPlanNotActive) {
		t.Fatalf("inactive plan error = %v, want %v", err, ErrMembershipPlanNotActive)
	}
}

func mustGym(t *testing.T, timezone string) Gym {
	t.Helper()
	gym, err := CreateGym("Zeus Gym", timezone, time.Date(2026, time.August, 1, 8, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("CreateGym returned error: %v", err)
	}
	return gym
}

func mustPlan(t *testing.T, gym Gym, kind MembershipValidityKind, duration int, unit MembershipDurationUnit, visits int) MembershipPlan {
	t.Helper()
	price, err := NewMoney(4500, "USD")
	if err != nil {
		t.Fatalf("NewMoney returned error: %v", err)
	}
	plan, err := CreateMembershipPlan(gym.ID(), "Plan", kind, duration, unit, visits, price, MembershipPlanStatusActive, time.Date(2026, time.August, 1, 8, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("CreateMembershipPlan returned error: %v", err)
	}
	return plan
}
