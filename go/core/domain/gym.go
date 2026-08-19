package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidGymID      = errors.New("invalid gym ID")
	ErrInvalidEntityName = errors.New("entity name cannot be blank")
	ErrInvalidTimezone   = errors.New("invalid IANA timezone")
	ErrInvalidTimestamp  = errors.New("timestamp cannot be zero")
	ErrInvalidTimeOrder  = errors.New("updated timestamp cannot be before created timestamp")
)

type GymID string

func NewGymID() (GymID, error) {
	value, err := NewUUID()
	if err != nil {
		return "", err
	}
	return GymID(value), nil
}

func ParseGymID(value string) (GymID, error) {
	if err := ValidateUUID(value); err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidGymID, err)
	}
	return GymID(value), nil
}

// Gym identifies the tenant that owns operational records.
type Gym struct {
	id        GymID
	name      string
	timezone  string
	createdAt time.Time
	updatedAt time.Time
}

func CreateGym(name, timezone string, createdAt time.Time) (Gym, error) {
	id, err := NewGymID()
	if err != nil {
		return Gym{}, fmt.Errorf("create gym ID: %w", err)
	}
	return NewGym(id, name, timezone, createdAt, createdAt)
}

func NewGym(id GymID, name, timezone string, createdAt, updatedAt time.Time) (Gym, error) {
	if _, err := ParseGymID(string(id)); err != nil {
		return Gym{}, err
	}
	normalizedName, err := normalizeEntityName(name)
	if err != nil {
		return Gym{}, err
	}
	normalizedTimezone, err := normalizeTimezone(timezone)
	if err != nil {
		return Gym{}, err
	}
	createdAt, updatedAt, err = normalizeLifecycleTimestamps(createdAt, updatedAt)
	if err != nil {
		return Gym{}, err
	}

	return Gym{
		id:        id,
		name:      normalizedName,
		timezone:  normalizedTimezone,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

func (g Gym) ID() GymID {
	return g.id
}

func (g Gym) Name() string {
	return g.name
}

func (g Gym) Timezone() string {
	return g.timezone
}

func (g Gym) CreatedAt() time.Time {
	return g.createdAt
}

func (g Gym) UpdatedAt() time.Time {
	return g.updatedAt
}

func normalizeEntityName(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", ErrInvalidEntityName
	}
	return value, nil
}

func normalizeTimezone(value string) (string, error) {
	if _, err := time.LoadLocation(value); err != nil {
		return "", fmt.Errorf("%w: %q", ErrInvalidTimezone, value)
	}
	return value, nil
}

func normalizeLifecycleTimestamps(createdAt, updatedAt time.Time) (time.Time, time.Time, error) {
	if createdAt.IsZero() || updatedAt.IsZero() {
		return time.Time{}, time.Time{}, ErrInvalidTimestamp
	}
	createdAt = createdAt.UTC()
	updatedAt = updatedAt.UTC()
	if updatedAt.Before(createdAt) {
		return time.Time{}, time.Time{}, ErrInvalidTimeOrder
	}
	return createdAt, updatedAt, nil
}
