package domain

import (
	"errors"
	"testing"
	"time"
)

func TestCreateGym(t *testing.T) {
	createdAt := time.Date(2026, time.August, 19, 10, 0, 0, 0, time.FixedZone("ECT", -5*60*60))
	gym, err := CreateGym("  Zeus Gym  ", "America/Guayaquil", createdAt)
	if err != nil {
		t.Fatalf("CreateGym returned error: %v", err)
	}
	if err := ValidateUUID(string(gym.ID())); err != nil {
		t.Fatalf("gym ID is invalid: %v", err)
	}
	if gym.Name() != "Zeus Gym" {
		t.Fatalf("name = %q, want Zeus Gym", gym.Name())
	}
	if gym.Timezone() != "America/Guayaquil" {
		t.Fatalf("timezone = %q, want America/Guayaquil", gym.Timezone())
	}
	if gym.CreatedAt().Location() != time.UTC || !gym.UpdatedAt().Equal(gym.CreatedAt()) {
		t.Fatalf("timestamps = %s / %s, want equal UTC values", gym.CreatedAt(), gym.UpdatedAt())
	}
}

func TestNewGymRejectsInvalidInput(t *testing.T) {
	id := mustGymID(t)
	createdAt := time.Date(2026, time.August, 19, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		id        GymID
		gymName   string
		timezone  string
		createdAt time.Time
		updatedAt time.Time
		wantErr   error
	}{
		{name: "invalid ID", id: "not-a-uuid", gymName: "Zeus", timezone: "UTC", createdAt: createdAt, updatedAt: createdAt, wantErr: ErrInvalidGymID},
		{name: "blank name", id: id, gymName: " \t", timezone: "UTC", createdAt: createdAt, updatedAt: createdAt, wantErr: ErrInvalidEntityName},
		{name: "invalid timezone", id: id, gymName: "Zeus", timezone: "Mars/Olympus", createdAt: createdAt, updatedAt: createdAt, wantErr: ErrInvalidTimezone},
		{name: "zero created timestamp", id: id, gymName: "Zeus", timezone: "UTC", updatedAt: createdAt, wantErr: ErrInvalidTimestamp},
		{name: "updated before created", id: id, gymName: "Zeus", timezone: "UTC", createdAt: createdAt, updatedAt: createdAt.Add(-time.Nanosecond), wantErr: ErrInvalidTimeOrder},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewGym(test.id, test.gymName, test.timezone, test.createdAt, test.updatedAt)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("NewGym() error = %v, want %v", err, test.wantErr)
			}
		})
	}
}

func mustGymID(t *testing.T) GymID {
	t.Helper()
	id, err := NewGymID()
	if err != nil {
		t.Fatalf("NewGymID returned error: %v", err)
	}
	return id
}
