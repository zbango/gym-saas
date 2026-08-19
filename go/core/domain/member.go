package domain

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"
)

var (
	ErrInvalidMemberID     = errors.New("invalid member ID")
	ErrInvalidMemberName   = errors.New("member name cannot be blank")
	ErrInvalidMemberPhone  = errors.New("member phone cannot be blank")
	ErrInvalidMemberEmail  = errors.New("invalid member email")
	ErrInvalidMemberDOB    = errors.New("invalid member date of birth")
	ErrInvalidMemberStatus = errors.New("invalid member status")
)

type MemberID string

func NewMemberID() (MemberID, error) {
	value, err := NewUUID()
	if err != nil {
		return "", err
	}
	return MemberID(value), nil
}

func ParseMemberID(value string) (MemberID, error) {
	if err := ValidateUUID(value); err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidMemberID, err)
	}
	return MemberID(value), nil
}

type MemberStatus string

const (
	MemberStatusActive   MemberStatus = "active"
	MemberStatusInactive MemberStatus = "inactive"
	MemberStatusBlocked  MemberStatus = "blocked"
)

func ParseMemberStatus(value string) (MemberStatus, error) {
	status := MemberStatus(value)
	switch status {
	case MemberStatusActive, MemberStatusInactive, MemberStatusBlocked:
		return status, nil
	default:
		return "", fmt.Errorf("%w: %q", ErrInvalidMemberStatus, value)
	}
}

// Member owns member identity and profile data. Membership and access rules
// are deliberately separate domain concerns.
type Member struct {
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
}

func CreateMember(gymID GymID, firstName, lastName, email, phone, identificationNumber, dateOfBirth, address string, status MemberStatus, createdAt time.Time) (Member, error) {
	id, err := NewMemberID()
	if err != nil {
		return Member{}, fmt.Errorf("create member ID: %w", err)
	}
	return NewMember(id, gymID, firstName, lastName, email, phone, identificationNumber, dateOfBirth, address, status, createdAt, createdAt)
}

func NewMember(id MemberID, gymID GymID, firstName, lastName, email, phone, identificationNumber, dateOfBirth, address string, status MemberStatus, createdAt, updatedAt time.Time) (Member, error) {
	if _, err := ParseMemberID(string(id)); err != nil {
		return Member{}, err
	}
	if _, err := ParseGymID(string(gymID)); err != nil {
		return Member{}, err
	}
	firstName, err := normalizeMemberName(firstName)
	if err != nil {
		return Member{}, err
	}
	lastName, err = normalizeMemberName(lastName)
	if err != nil {
		return Member{}, err
	}
	email, err = normalizeMemberEmail(email)
	if err != nil {
		return Member{}, err
	}
	phone, err = normalizeMemberPhone(phone)
	if err != nil {
		return Member{}, err
	}
	dateOfBirth, err = normalizeMemberDateOfBirth(dateOfBirth)
	if err != nil {
		return Member{}, err
	}
	status, err = ParseMemberStatus(string(status))
	if err != nil {
		return Member{}, err
	}
	createdAt, updatedAt, err = normalizeLifecycleTimestamps(createdAt, updatedAt)
	if err != nil {
		return Member{}, err
	}

	return Member{
		id:                   id,
		gymID:                gymID,
		firstName:            firstName,
		lastName:             lastName,
		email:                email,
		phone:                phone,
		identificationNumber: strings.TrimSpace(identificationNumber),
		dateOfBirth:          dateOfBirth,
		address:              strings.TrimSpace(address),
		status:               status,
		createdAt:            createdAt,
		updatedAt:            updatedAt,
	}, nil
}

func (m Member) ID() MemberID {
	return m.id
}

func (m Member) GymID() GymID {
	return m.gymID
}

func (m Member) FirstName() string {
	return m.firstName
}

func (m Member) LastName() string {
	return m.lastName
}

func (m Member) Email() string {
	return m.email
}

func (m Member) Phone() string {
	return m.phone
}

func (m Member) IdentificationNumber() string {
	return m.identificationNumber
}

func (m Member) DateOfBirth() string {
	return m.dateOfBirth
}

func (m Member) Address() string {
	return m.address
}

func (m Member) Status() MemberStatus {
	return m.status
}

func (m Member) CreatedAt() time.Time {
	return m.createdAt
}

func (m Member) UpdatedAt() time.Time {
	return m.updatedAt
}

func normalizeMemberName(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", ErrInvalidMemberName
	}
	return value, nil
}

func normalizeMemberPhone(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", ErrInvalidMemberPhone
	}
	return value, nil
}

func normalizeMemberEmail(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", nil
	}
	address, err := mail.ParseAddress(value)
	if err != nil || address.Address != value {
		return "", fmt.Errorf("%w: %q", ErrInvalidMemberEmail, value)
	}
	return value, nil
}

func normalizeMemberDateOfBirth(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", nil
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil || parsed.Format("2006-01-02") != value {
		return "", fmt.Errorf("%w: %q", ErrInvalidMemberDOB, value)
	}
	return value, nil
}
