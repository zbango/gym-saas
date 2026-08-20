package application

import (
	"context"
	"fmt"
	"time"

	"github.com/zbango/gym-saas/go/core/domain"
	"github.com/zbango/gym-saas/go/core/ports"
)

// MemberInput is the mutable profile data accepted by the member use cases.
// Validation remains in the domain when the input becomes a Member.
type MemberInput struct {
	FirstName            string
	LastName             string
	Email                string
	Phone                string
	IdentificationNumber string
	DateOfBirth          string
	Address              string
	Status               string
}

// MemberService coordinates member use cases for one gym tenant.
type MemberService struct {
	repository ports.MemberRepository
	gymID      domain.GymID
	now        func() time.Time
}

func NewMemberService(repository ports.MemberRepository, gymID domain.GymID, now func() time.Time) (*MemberService, error) {
	if repository == nil {
		return nil, fmt.Errorf("member repository is required")
	}
	if _, err := domain.ParseGymID(string(gymID)); err != nil {
		return nil, err
	}
	if now == nil {
		now = time.Now
	}
	return &MemberService{repository: repository, gymID: gymID, now: now}, nil
}

func (s *MemberService) Create(ctx context.Context, input MemberInput) (domain.Member, error) {
	member, err := memberFromInput(s.gymID, "", time.Time{}, input, s.now())
	if err != nil {
		return domain.Member{}, err
	}
	if err := s.repository.Create(ctx, member); err != nil {
		return domain.Member{}, fmt.Errorf("save member: %w", err)
	}
	return member, nil
}

func (s *MemberService) List(ctx context.Context) ([]domain.Member, error) {
	members, err := s.repository.List(ctx, s.gymID)
	if err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}
	return members, nil
}

func (s *MemberService) Update(ctx context.Context, id string, input MemberInput) (domain.Member, error) {
	memberID, err := domain.ParseMemberID(id)
	if err != nil {
		return domain.Member{}, err
	}
	existing, err := s.repository.Get(ctx, s.gymID, memberID)
	if err != nil {
		return domain.Member{}, fmt.Errorf("find member: %w", err)
	}
	member, err := memberFromInput(s.gymID, existing.ID(), existing.CreatedAt(), input, s.now())
	if err != nil {
		return domain.Member{}, err
	}
	if err := s.repository.Update(ctx, member); err != nil {
		return domain.Member{}, fmt.Errorf("update member: %w", err)
	}
	return member, nil
}

func (s *MemberService) Archive(ctx context.Context, id string) error {
	memberID, err := domain.ParseMemberID(id)
	if err != nil {
		return err
	}
	if err := s.repository.Archive(ctx, s.gymID, memberID, s.now()); err != nil {
		return fmt.Errorf("archive member: %w", err)
	}
	return nil
}

func memberFromInput(gymID domain.GymID, existingID domain.MemberID, createdAt time.Time, input MemberInput, updatedAt time.Time) (domain.Member, error) {
	status, err := domain.ParseMemberStatus(input.Status)
	if err != nil {
		return domain.Member{}, err
	}
	if existingID == "" {
		return domain.CreateMember(gymID, input.FirstName, input.LastName, input.Email, input.Phone, input.IdentificationNumber, input.DateOfBirth, input.Address, status, updatedAt)
	}
	return domain.NewMember(existingID, gymID, input.FirstName, input.LastName, input.Email, input.Phone, input.IdentificationNumber, input.DateOfBirth, input.Address, status, createdAt, updatedAt)
}
