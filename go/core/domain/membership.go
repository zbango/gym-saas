package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidMembershipID          = errors.New("invalid membership ID")
	ErrInvalidMembershipStatus      = errors.New("invalid membership status")
	ErrInvalidMembershipDates       = errors.New("invalid membership dates")
	ErrInvalidMembershipTerms       = errors.New("invalid membership terms")
	ErrInvalidCancellationReason    = errors.New("cancellation reason cannot be blank")
	ErrMembershipPlanGymMismatch    = errors.New("membership plan belongs to another gym")
	ErrMembershipPlanNotActive      = errors.New("membership plan is not active")
	ErrMembershipNotPending         = errors.New("membership is not pending")
	ErrMembershipActivationTooEarly = errors.New("membership cannot activate before its start time")
	ErrMembershipNotValid           = errors.New("membership is not valid at this time")
	ErrMembershipAlreadyExpired     = errors.New("membership is already expired")
	ErrMembershipCannotCancel       = errors.New("membership cannot be cancelled in its current status")
	ErrVisitAllowanceNotApplicable  = errors.New("visit allowance does not apply to a time membership")
)

type MembershipID string

func NewMembershipID() (MembershipID, error) {
	value, err := NewUUID()
	if err != nil {
		return "", err
	}
	return MembershipID(value), nil
}

func ParseMembershipID(value string) (MembershipID, error) {
	if err := ValidateUUID(value); err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidMembershipID, err)
	}
	return MembershipID(value), nil
}

type MembershipStatus string

const (
	MembershipStatusPending   MembershipStatus = "pending"
	MembershipStatusActive    MembershipStatus = "active"
	MembershipStatusExpired   MembershipStatus = "expired"
	MembershipStatusCancelled MembershipStatus = "cancelled"
)

func ParseMembershipStatus(value string) (MembershipStatus, error) {
	switch status := MembershipStatus(value); status {
	case MembershipStatusPending, MembershipStatusActive, MembershipStatusExpired, MembershipStatusCancelled:
		return status, nil
	default:
		return "", fmt.Errorf("%w: %q", ErrInvalidMembershipStatus, value)
	}
}

// Membership is a purchased snapshot of a plan. Calendar bounds are computed
// in the owning gym's timezone and persisted in UTC; the end instant is
// exclusive. Visit plans expire at the earlier of their end instant and visit
// exhaustion.
type Membership struct {
	id            MembershipID
	gymID         GymID
	memberID      MemberID
	planID        MembershipPlanID
	status        MembershipStatus
	startsAt      time.Time
	endsAt        time.Time
	activatedAt   time.Time
	cancelledAt   time.Time
	cancelReason  string
	validityKind  MembershipValidityKind
	durationValue int
	durationUnit  MembershipDurationUnit
	visitLimit    int
	visitsRemain  int
	price         Money
	createdAt     time.Time
	updatedAt     time.Time
}

// StartMembership snapshots an active plan for one member. A purchase whose
// start is in the future remains pending until Activate is called at its start
// instant; an immediate purchase becomes active immediately.
func StartMembership(gym Gym, memberID MemberID, plan MembershipPlan, startsAt, createdAt time.Time) (Membership, error) {
	if plan.GymID() != gym.ID() {
		return Membership{}, ErrMembershipPlanGymMismatch
	}
	if plan.Status() != MembershipPlanStatusActive {
		return Membership{}, ErrMembershipPlanNotActive
	}
	if _, err := ParseMemberID(string(memberID)); err != nil {
		return Membership{}, err
	}
	startsAt, err := normalizeMembershipTimestamp(startsAt)
	if err != nil {
		return Membership{}, err
	}
	createdAt, err = normalizeMembershipTimestamp(createdAt)
	if err != nil {
		return Membership{}, err
	}
	if startsAt.Before(createdAt) {
		return Membership{}, fmt.Errorf("%w: start cannot precede purchase", ErrInvalidMembershipDates)
	}
	endsAt, err := membershipEndAt(startsAt, plan.DurationValue(), plan.DurationUnit(), gym.Timezone())
	if err != nil {
		return Membership{}, err
	}
	id, err := NewMembershipID()
	if err != nil {
		return Membership{}, fmt.Errorf("create membership ID: %w", err)
	}

	status := MembershipStatusPending
	activatedAt := time.Time{}
	if startsAt.Equal(createdAt) {
		status = MembershipStatusActive
		activatedAt = startsAt
	}
	visitsRemaining := 0
	if plan.ValidityKind() == MembershipValidityVisits {
		visitsRemaining = plan.VisitLimit()
	}
	return NewMembership(
		id, gym.ID(), memberID, plan.ID(), status, startsAt, endsAt, activatedAt,
		time.Time{}, "", plan.ValidityKind(), plan.DurationValue(), plan.DurationUnit(),
		plan.VisitLimit(), visitsRemaining, plan.Price(), createdAt, createdAt,
	)
}

// NewMembership reconstructs a persisted membership. Optional activation and
// cancellation timestamps use the zero time only when their status permits it.
func NewMembership(id MembershipID, gymID GymID, memberID MemberID, planID MembershipPlanID, status MembershipStatus, startsAt, endsAt, activatedAt, cancelledAt time.Time, cancellationReason string, validityKind MembershipValidityKind, durationValue int, durationUnit MembershipDurationUnit, visitLimit, visitsRemaining int, price Money, createdAt, updatedAt time.Time) (Membership, error) {
	if _, err := ParseMembershipID(string(id)); err != nil {
		return Membership{}, err
	}
	if _, err := ParseGymID(string(gymID)); err != nil {
		return Membership{}, err
	}
	if _, err := ParseMemberID(string(memberID)); err != nil {
		return Membership{}, err
	}
	if _, err := ParseMembershipPlanID(string(planID)); err != nil {
		return Membership{}, err
	}
	status, err := ParseMembershipStatus(string(status))
	if err != nil {
		return Membership{}, err
	}
	startsAt, err = normalizeMembershipTimestamp(startsAt)
	if err != nil {
		return Membership{}, err
	}
	endsAt, err = normalizeMembershipTimestamp(endsAt)
	if err != nil {
		return Membership{}, err
	}
	if !endsAt.After(startsAt) {
		return Membership{}, fmt.Errorf("%w: end must be after start", ErrInvalidMembershipDates)
	}
	createdAt, updatedAt, err = normalizeLifecycleTimestamps(createdAt, updatedAt)
	if err != nil {
		return Membership{}, err
	}
	validityKind, durationUnit, err = validateMembershipTerms(validityKind, durationValue, durationUnit, visitLimit, visitsRemaining, price)
	if err != nil {
		return Membership{}, err
	}
	if status == MembershipStatusExpired && updatedAt.Before(endsAt) && !(validityKind == MembershipValidityVisits && visitsRemaining == 0) {
		return Membership{}, fmt.Errorf("%w: expired status has no elapsed validity bound", ErrInvalidMembershipDates)
	}
	cancellationReason = strings.TrimSpace(cancellationReason)
	if err := validateMembershipLifecycle(status, startsAt, endsAt, activatedAt, cancelledAt, cancellationReason, createdAt); err != nil {
		return Membership{}, err
	}
	if !activatedAt.IsZero() {
		activatedAt = activatedAt.UTC()
	}
	if !cancelledAt.IsZero() {
		cancelledAt = cancelledAt.UTC()
	}
	return Membership{
		id: id, gymID: gymID, memberID: memberID, planID: planID, status: status,
		startsAt: startsAt, endsAt: endsAt, activatedAt: activatedAt, cancelledAt: cancelledAt,
		cancelReason: cancellationReason, validityKind: validityKind, durationValue: durationValue,
		durationUnit: durationUnit, visitLimit: visitLimit, visitsRemain: visitsRemaining,
		price: price, createdAt: createdAt, updatedAt: updatedAt,
	}, nil
}

func (m Membership) ID() MembershipID                     { return m.id }
func (m Membership) GymID() GymID                         { return m.gymID }
func (m Membership) MemberID() MemberID                   { return m.memberID }
func (m Membership) MembershipPlanID() MembershipPlanID   { return m.planID }
func (m Membership) Status() MembershipStatus             { return m.status }
func (m Membership) StartsAt() time.Time                  { return m.startsAt }
func (m Membership) EndsAt() time.Time                    { return m.endsAt }
func (m Membership) ActivatedAt() time.Time               { return m.activatedAt }
func (m Membership) CancelledAt() time.Time               { return m.cancelledAt }
func (m Membership) CancellationReason() string           { return m.cancelReason }
func (m Membership) ValidityKind() MembershipValidityKind { return m.validityKind }
func (m Membership) DurationValue() int                   { return m.durationValue }
func (m Membership) DurationUnit() MembershipDurationUnit { return m.durationUnit }
func (m Membership) VisitLimit() int                      { return m.visitLimit }
func (m Membership) VisitsRemaining() int                 { return m.visitsRemain }
func (m Membership) Price() Money                         { return m.price }
func (m Membership) CreatedAt() time.Time                 { return m.createdAt }
func (m Membership) UpdatedAt() time.Time                 { return m.updatedAt }

func (m Membership) IsValidAt(at time.Time) bool {
	if m.status != MembershipStatusActive || at.IsZero() {
		return false
	}
	at = at.UTC()
	if at.Before(m.startsAt) || !at.Before(m.endsAt) {
		return false
	}
	return m.validityKind != MembershipValidityVisits || m.visitsRemain > 0
}

func (m Membership) Activate(at time.Time) (Membership, error) {
	if m.status != MembershipStatusPending {
		return Membership{}, ErrMembershipNotPending
	}
	at, err := normalizeMembershipTimestamp(at)
	if err != nil {
		return Membership{}, err
	}
	if at.Before(m.startsAt) {
		return Membership{}, ErrMembershipActivationTooEarly
	}
	if !at.Before(m.endsAt) {
		return Membership{}, ErrMembershipAlreadyExpired
	}
	m.status = MembershipStatusActive
	m.activatedAt = at
	m.updatedAt = at
	return m, nil
}

func (m Membership) ConsumeVisit(at time.Time) (Membership, error) {
	if m.validityKind != MembershipValidityVisits {
		return Membership{}, ErrVisitAllowanceNotApplicable
	}
	if !m.IsValidAt(at) {
		return Membership{}, ErrMembershipNotValid
	}
	at = at.UTC()
	if at.Before(m.updatedAt) {
		return Membership{}, fmt.Errorf("%w: visit precedes the latest membership update", ErrInvalidMembershipDates)
	}
	m.visitsRemain--
	m.updatedAt = at
	return m, nil
}

func (m Membership) Cancel(at time.Time, reason string) (Membership, error) {
	if m.status != MembershipStatusActive && m.status != MembershipStatusPending {
		return Membership{}, ErrMembershipCannotCancel
	}
	at, err := normalizeMembershipTimestamp(at)
	if err != nil {
		return Membership{}, err
	}
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return Membership{}, ErrInvalidCancellationReason
	}
	if m.isExpiredAt(at) {
		return Membership{}, ErrMembershipAlreadyExpired
	}
	if at.Before(m.createdAt) {
		return Membership{}, fmt.Errorf("%w: cancellation precedes purchase", ErrInvalidMembershipDates)
	}
	m.status = MembershipStatusCancelled
	m.cancelledAt = at
	m.cancelReason = reason
	m.updatedAt = at
	return m, nil
}

func (m Membership) MarkExpired(at time.Time) (Membership, error) {
	if m.status != MembershipStatusActive {
		return Membership{}, ErrMembershipNotValid
	}
	at, err := normalizeMembershipTimestamp(at)
	if err != nil {
		return Membership{}, err
	}
	if !m.isExpiredAt(at) {
		return Membership{}, ErrMembershipNotValid
	}
	m.status = MembershipStatusExpired
	m.updatedAt = at
	return m, nil
}

func (m Membership) isExpiredAt(at time.Time) bool {
	return !at.Before(m.endsAt) || (m.validityKind == MembershipValidityVisits && m.visitsRemain == 0)
}

func validateMembershipTerms(validityKind MembershipValidityKind, durationValue int, durationUnit MembershipDurationUnit, visitLimit, visitsRemaining int, price Money) (MembershipValidityKind, MembershipDurationUnit, error) {
	validityKind, err := ParseMembershipValidityKind(string(validityKind))
	if err != nil {
		return "", "", err
	}
	if durationValue <= 0 {
		return "", "", fmt.Errorf("%w: duration must be greater than zero", ErrInvalidMembershipTerms)
	}
	durationUnit, err = ParseMembershipDurationUnit(string(durationUnit))
	if err != nil {
		return "", "", err
	}
	if _, err := NewMoney(price.Cents(), price.Currency()); err != nil {
		return "", "", err
	}
	switch validityKind {
	case MembershipValidityTime:
		if visitLimit != 0 || visitsRemaining != 0 {
			return "", "", fmt.Errorf("%w: time memberships cannot retain visits", ErrInvalidMembershipTerms)
		}
	case MembershipValidityVisits:
		if visitLimit <= 0 || visitsRemaining < 0 || visitsRemaining > visitLimit {
			return "", "", fmt.Errorf("%w: invalid visit allowance", ErrInvalidMembershipTerms)
		}
	}
	return validityKind, durationUnit, nil
}

func validateMembershipLifecycle(status MembershipStatus, startsAt, endsAt, activatedAt, cancelledAt time.Time, cancellationReason string, createdAt time.Time) error {
	switch status {
	case MembershipStatusPending:
		if !activatedAt.IsZero() || !cancelledAt.IsZero() || cancellationReason != "" {
			return fmt.Errorf("%w: pending membership cannot have activation or cancellation data", ErrInvalidMembershipDates)
		}
	case MembershipStatusActive, MembershipStatusExpired:
		if activatedAt.IsZero() || !cancelledAt.IsZero() || cancellationReason != "" {
			return fmt.Errorf("%w: active or expired membership requires activation only", ErrInvalidMembershipDates)
		}
		if activatedAt.UTC().Before(startsAt) || !activatedAt.UTC().Before(endsAt) {
			return fmt.Errorf("%w: activation must be within membership bounds", ErrInvalidMembershipDates)
		}
	case MembershipStatusCancelled:
		if cancelledAt.IsZero() || cancellationReason == "" {
			return fmt.Errorf("%w: cancelled membership requires a cancellation timestamp and reason", ErrInvalidMembershipDates)
		}
		if cancelledAt.UTC().Before(createdAt) || !cancelledAt.UTC().Before(endsAt) {
			return fmt.Errorf("%w: cancellation must be after purchase and before expiry", ErrInvalidMembershipDates)
		}
		if !activatedAt.IsZero() && (activatedAt.UTC().Before(startsAt) || !activatedAt.UTC().Before(endsAt)) {
			return fmt.Errorf("%w: activation must be within membership bounds", ErrInvalidMembershipDates)
		}
	}
	return nil
}

func normalizeMembershipTimestamp(value time.Time) (time.Time, error) {
	if value.IsZero() {
		return time.Time{}, ErrInvalidMembershipDates
	}
	return value.UTC(), nil
}

func membershipEndAt(startsAt time.Time, durationValue int, durationUnit MembershipDurationUnit, timezone string) (time.Time, error) {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("load gym timezone: %w", err)
	}
	localStart := startsAt.In(location)
	switch durationUnit {
	case MembershipDurationDays:
		return localStart.AddDate(0, 0, durationValue).UTC(), nil
	case MembershipDurationWeeks:
		return localStart.AddDate(0, 0, 7*durationValue).UTC(), nil
	case MembershipDurationMonths:
		return addMonthsClamped(localStart, durationValue).UTC(), nil
	case MembershipDurationYears:
		return addMonthsClamped(localStart, 12*durationValue).UTC(), nil
	default:
		return time.Time{}, fmt.Errorf("%w: %q", ErrInvalidMembershipTerms, durationUnit)
	}
}

func addMonthsClamped(value time.Time, months int) time.Time {
	targetMonth := int(value.Month()) + months
	targetYear := value.Year() + (targetMonth-1)/12
	targetMonth = (targetMonth-1)%12 + 1
	lastDay := time.Date(targetYear, time.Month(targetMonth)+1, 0, 0, 0, 0, 0, value.Location()).Day()
	day := min(value.Day(), lastDay)
	return time.Date(targetYear, time.Month(targetMonth), day, value.Hour(), value.Minute(), value.Second(), value.Nanosecond(), value.Location())
}
