package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidMembershipPlanID     = errors.New("invalid membership plan ID")
	ErrInvalidMembershipPlanName   = errors.New("membership plan name cannot be blank")
	ErrInvalidValidityKind         = errors.New("invalid membership plan validity kind")
	ErrInvalidDurationValue        = errors.New("membership plan duration must be greater than zero")
	ErrInvalidDurationUnit         = errors.New("invalid membership plan duration unit")
	ErrInvalidVisitLimit           = errors.New("membership plan visit limit must be greater than zero")
	ErrInvalidMembershipPlanStatus = errors.New("invalid membership plan status")
)

type MembershipPlanID string

func NewMembershipPlanID() (MembershipPlanID, error) {
	value, err := NewUUID()
	if err != nil {
		return "", err
	}
	return MembershipPlanID(value), nil
}

func ParseMembershipPlanID(value string) (MembershipPlanID, error) {
	if err := ValidateUUID(value); err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidMembershipPlanID, err)
	}
	return MembershipPlanID(value), nil
}

type MembershipValidityKind string

const (
	MembershipValidityTime   MembershipValidityKind = "time"
	MembershipValidityVisits MembershipValidityKind = "visits"
)

func ParseMembershipValidityKind(value string) (MembershipValidityKind, error) {
	switch kind := MembershipValidityKind(value); kind {
	case MembershipValidityTime, MembershipValidityVisits:
		return kind, nil
	default:
		return "", fmt.Errorf("%w: %q", ErrInvalidValidityKind, value)
	}
}

type MembershipDurationUnit string

const (
	MembershipDurationDays   MembershipDurationUnit = "days"
	MembershipDurationWeeks  MembershipDurationUnit = "weeks"
	MembershipDurationMonths MembershipDurationUnit = "months"
	MembershipDurationYears  MembershipDurationUnit = "years"
)

func ParseMembershipDurationUnit(value string) (MembershipDurationUnit, error) {
	switch unit := MembershipDurationUnit(value); unit {
	case MembershipDurationDays, MembershipDurationWeeks, MembershipDurationMonths, MembershipDurationYears:
		return unit, nil
	default:
		return "", fmt.Errorf("%w: %q", ErrInvalidDurationUnit, value)
	}
}

type MembershipPlanStatus string

const (
	MembershipPlanStatusActive   MembershipPlanStatus = "active"
	MembershipPlanStatusInactive MembershipPlanStatus = "inactive"
)

func ParseMembershipPlanStatus(value string) (MembershipPlanStatus, error) {
	switch status := MembershipPlanStatus(value); status {
	case MembershipPlanStatusActive, MembershipPlanStatusInactive:
		return status, nil
	default:
		return "", fmt.Errorf("%w: %q", ErrInvalidMembershipPlanStatus, value)
	}
}

// MembershipPlan is an administrator-configured offer. A visit plan always
// has both a visit limit and a calendar validity window; a purchase expires
// when either condition is reached.
type MembershipPlan struct {
	id            MembershipPlanID
	gymID         GymID
	name          string
	validityKind  MembershipValidityKind
	durationValue int
	durationUnit  MembershipDurationUnit
	visitLimit    int
	price         Money
	status        MembershipPlanStatus
	createdAt     time.Time
	updatedAt     time.Time
}

func CreateMembershipPlan(gymID GymID, name string, validityKind MembershipValidityKind, durationValue int, durationUnit MembershipDurationUnit, visitLimit int, price Money, status MembershipPlanStatus, createdAt time.Time) (MembershipPlan, error) {
	id, err := NewMembershipPlanID()
	if err != nil {
		return MembershipPlan{}, fmt.Errorf("create membership plan ID: %w", err)
	}
	return NewMembershipPlan(id, gymID, name, validityKind, durationValue, durationUnit, visitLimit, price, status, createdAt, createdAt)
}

func NewMembershipPlan(id MembershipPlanID, gymID GymID, name string, validityKind MembershipValidityKind, durationValue int, durationUnit MembershipDurationUnit, visitLimit int, price Money, status MembershipPlanStatus, createdAt, updatedAt time.Time) (MembershipPlan, error) {
	if _, err := ParseMembershipPlanID(string(id)); err != nil {
		return MembershipPlan{}, err
	}
	if _, err := ParseGymID(string(gymID)); err != nil {
		return MembershipPlan{}, err
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return MembershipPlan{}, ErrInvalidMembershipPlanName
	}
	validityKind, err := ParseMembershipValidityKind(string(validityKind))
	if err != nil {
		return MembershipPlan{}, err
	}
	if durationValue <= 0 {
		return MembershipPlan{}, ErrInvalidDurationValue
	}
	durationUnit, err = ParseMembershipDurationUnit(string(durationUnit))
	if err != nil {
		return MembershipPlan{}, err
	}
	if validityKind == MembershipValidityTime && visitLimit != 0 {
		return MembershipPlan{}, fmt.Errorf("%w: time plans cannot have visits", ErrInvalidVisitLimit)
	}
	if validityKind == MembershipValidityVisits && visitLimit <= 0 {
		return MembershipPlan{}, ErrInvalidVisitLimit
	}
	if _, err := NewMoney(price.Cents(), price.Currency()); err != nil {
		return MembershipPlan{}, err
	}
	status, err = ParseMembershipPlanStatus(string(status))
	if err != nil {
		return MembershipPlan{}, err
	}
	createdAt, updatedAt, err = normalizeLifecycleTimestamps(createdAt, updatedAt)
	if err != nil {
		return MembershipPlan{}, err
	}
	return MembershipPlan{id: id, gymID: gymID, name: name, validityKind: validityKind, durationValue: durationValue, durationUnit: durationUnit, visitLimit: visitLimit, price: price, status: status, createdAt: createdAt, updatedAt: updatedAt}, nil
}

func (p MembershipPlan) ID() MembershipPlanID                 { return p.id }
func (p MembershipPlan) GymID() GymID                         { return p.gymID }
func (p MembershipPlan) Name() string                         { return p.name }
func (p MembershipPlan) ValidityKind() MembershipValidityKind { return p.validityKind }
func (p MembershipPlan) DurationValue() int                   { return p.durationValue }
func (p MembershipPlan) DurationUnit() MembershipDurationUnit { return p.durationUnit }
func (p MembershipPlan) VisitLimit() int                      { return p.visitLimit }
func (p MembershipPlan) Price() Money                         { return p.price }
func (p MembershipPlan) Status() MembershipPlanStatus         { return p.status }
func (p MembershipPlan) CreatedAt() time.Time                 { return p.createdAt }
func (p MembershipPlan) UpdatedAt() time.Time                 { return p.updatedAt }
