package domain

import (
	"errors"
	"testing"
	"time"
)

func TestCreateMember(t *testing.T) {
	createdAt := time.Date(2026, time.August, 19, 10, 0, 0, 0, time.FixedZone("ECT", -5*60*60))
	member, err := CreateMember(
		mustGymID(t),
		"  Ada  ",
		"  Lovelace  ",
		"  ada@example.com ",
		"  +593 99 123 4567 ",
		"  0102030405 ",
		"  1815-12-10 ",
		"  12 St. James's Square ",
		MemberStatusActive,
		createdAt,
	)
	if err != nil {
		t.Fatalf("CreateMember returned error: %v", err)
	}
	if err := ValidateUUID(string(member.ID())); err != nil {
		t.Fatalf("member ID is invalid: %v", err)
	}
	if member.FirstName() != "Ada" || member.LastName() != "Lovelace" {
		t.Fatalf("member names = %q %q, want Ada Lovelace", member.FirstName(), member.LastName())
	}
	if member.Email() != "ada@example.com" || member.Phone() != "+593 99 123 4567" || member.IdentificationNumber() != "0102030405" || member.DateOfBirth() != "1815-12-10" || member.Address() != "12 St. James's Square" {
		t.Fatalf("member fields were not normalized: %#v", member)
	}
	if member.Status() != MemberStatusActive {
		t.Fatalf("status = %q, want active", member.Status())
	}
	if member.CreatedAt().Location() != time.UTC || !member.UpdatedAt().Equal(member.CreatedAt()) {
		t.Fatalf("timestamps = %s / %s, want equal UTC values", member.CreatedAt(), member.UpdatedAt())
	}
}

func TestCreateMemberAllowsOptionalEmailAndIdentificationNumber(t *testing.T) {
	member, err := CreateMember(
		mustGymID(t), "Ada", "Lovelace", "", "555-0100", "", "", "", MemberStatusInactive,
		time.Date(2026, time.August, 19, 10, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("CreateMember returned error: %v", err)
	}
	if member.Email() != "" || member.IdentificationNumber() != "" || member.DateOfBirth() != "" || member.Address() != "" || member.Status() != MemberStatusInactive {
		t.Fatalf("optional fields/status = %#v, want empty optional fields and inactive", member)
	}
}

func TestNewMemberRejectsInvalidInput(t *testing.T) {
	id := mustMemberID(t)
	gymID := mustGymID(t)
	createdAt := time.Date(2026, time.August, 19, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name                 string
		id                   MemberID
		gymID                GymID
		firstName            string
		lastName             string
		email                string
		phone                string
		identificationNumber string
		dateOfBirth          string
		address              string
		status               MemberStatus
		createdAt            time.Time
		updatedAt            time.Time
		wantErr              error
	}{
		{name: "invalid member ID", id: "not-a-uuid", gymID: gymID, firstName: "Ada", lastName: "Lovelace", phone: "555-0100", status: MemberStatusActive, createdAt: createdAt, updatedAt: createdAt, wantErr: ErrInvalidMemberID},
		{name: "invalid gym ID", id: id, gymID: "not-a-uuid", firstName: "Ada", lastName: "Lovelace", phone: "555-0100", status: MemberStatusActive, createdAt: createdAt, updatedAt: createdAt, wantErr: ErrInvalidGymID},
		{name: "blank first name", id: id, gymID: gymID, firstName: " ", lastName: "Lovelace", phone: "555-0100", status: MemberStatusActive, createdAt: createdAt, updatedAt: createdAt, wantErr: ErrInvalidMemberName},
		{name: "blank last name", id: id, gymID: gymID, firstName: "Ada", lastName: " ", phone: "555-0100", status: MemberStatusActive, createdAt: createdAt, updatedAt: createdAt, wantErr: ErrInvalidMemberName},
		{name: "blank phone", id: id, gymID: gymID, firstName: "Ada", lastName: "Lovelace", phone: " ", status: MemberStatusActive, createdAt: createdAt, updatedAt: createdAt, wantErr: ErrInvalidMemberPhone},
		{name: "invalid email", id: id, gymID: gymID, firstName: "Ada", lastName: "Lovelace", email: "not-an-email", phone: "555-0100", status: MemberStatusActive, createdAt: createdAt, updatedAt: createdAt, wantErr: ErrInvalidMemberEmail},
		{name: "invalid date of birth", id: id, gymID: gymID, firstName: "Ada", lastName: "Lovelace", phone: "555-0100", dateOfBirth: "10/12/1815", status: MemberStatusActive, createdAt: createdAt, updatedAt: createdAt, wantErr: ErrInvalidMemberDOB},
		{name: "invalid status", id: id, gymID: gymID, firstName: "Ada", lastName: "Lovelace", phone: "555-0100", status: "deleted", createdAt: createdAt, updatedAt: createdAt, wantErr: ErrInvalidMemberStatus},
		{name: "zero timestamp", id: id, gymID: gymID, firstName: "Ada", lastName: "Lovelace", phone: "555-0100", status: MemberStatusActive, updatedAt: createdAt, wantErr: ErrInvalidTimestamp},
		{name: "reversed timestamps", id: id, gymID: gymID, firstName: "Ada", lastName: "Lovelace", phone: "555-0100", status: MemberStatusActive, createdAt: createdAt, updatedAt: createdAt.Add(-time.Nanosecond), wantErr: ErrInvalidTimeOrder},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewMember(
				test.id, test.gymID, test.firstName, test.lastName, test.email, test.phone,
				test.identificationNumber, test.dateOfBirth, test.address, test.status, test.createdAt, test.updatedAt,
			)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("NewMember() error = %v, want %v", err, test.wantErr)
			}
		})
	}
}

func TestParseMemberStatus(t *testing.T) {
	for _, status := range []MemberStatus{MemberStatusActive, MemberStatusInactive, MemberStatusBlocked} {
		parsed, err := ParseMemberStatus(string(status))
		if err != nil || parsed != status {
			t.Fatalf("ParseMemberStatus(%q) = %q, %v", status, parsed, err)
		}
	}
}

func mustMemberID(t *testing.T) MemberID {
	t.Helper()
	id, err := NewMemberID()
	if err != nil {
		t.Fatalf("NewMemberID returned error: %v", err)
	}
	return id
}
