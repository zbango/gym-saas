package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/zbango/gym-saas/go/core/domain"
	"github.com/zbango/gym-saas/go/core/ports"
)

func TestMemberServiceCreatesUpdatesAndArchivesMember(t *testing.T) {
	gymID := mustGymID(t)
	repository := &memoryMemberRepository{members: map[domain.MemberID]domain.Member{}}
	times := []time.Time{
		time.Date(2026, time.August, 19, 10, 0, 0, 0, time.UTC),
		time.Date(2026, time.August, 19, 11, 0, 0, 0, time.UTC),
		time.Date(2026, time.August, 19, 12, 0, 0, 0, time.UTC),
	}
	service, err := NewMemberService(repository, gymID, func() time.Time {
		value := times[0]
		times = times[1:]
		return value
	})
	if err != nil {
		t.Fatalf("NewMemberService returned error: %v", err)
	}

	created, err := service.Create(context.Background(), memberInput("Ada", "active"))
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	updatedInput := memberInput("Augusta", "inactive")
	updated, err := service.Update(context.Background(), string(created.ID()), updatedInput)
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if updated.ID() != created.ID() || updated.CreatedAt() != created.CreatedAt() || updated.UpdatedAt().Equal(created.UpdatedAt()) {
		t.Fatalf("updated lifecycle = %#v, created = %#v", updated, created)
	}
	if updated.FirstName() != "Augusta" || updated.Status() != domain.MemberStatusInactive {
		t.Fatalf("updated member = %#v", updated)
	}
	if err := service.Archive(context.Background(), string(created.ID())); err != nil {
		t.Fatalf("Archive returned error: %v", err)
	}
	members, err := service.List(context.Background())
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(members) != 0 {
		t.Fatalf("List length = %d, want 0 after archive", len(members))
	}
}

func TestMemberServiceRejectsInvalidInputBeforePersistence(t *testing.T) {
	repository := &memoryMemberRepository{members: map[domain.MemberID]domain.Member{}}
	service, err := NewMemberService(repository, mustGymID(t), time.Now)
	if err != nil {
		t.Fatalf("NewMemberService returned error: %v", err)
	}
	input := memberInput("Ada", "active")
	input.Email = "not-an-email"
	if _, err := service.Create(context.Background(), input); !errors.Is(err, domain.ErrInvalidMemberEmail) {
		t.Fatalf("Create error = %v, want %v", err, domain.ErrInvalidMemberEmail)
	}
	if len(repository.members) != 0 {
		t.Fatal("repository persisted an invalid member")
	}
}

func memberInput(firstName, status string) MemberInput {
	return MemberInput{
		FirstName: firstName, LastName: "Lovelace", Email: "ada@example.com", Phone: "555-0100",
		IdentificationNumber: "0102030405", DateOfBirth: "1815-12-10", Address: "12 St. James's Square", Status: status,
	}
}

type memoryMemberRepository struct {
	members  map[domain.MemberID]domain.Member
	archived map[domain.MemberID]bool
}

func (r *memoryMemberRepository) Create(_ context.Context, member domain.Member) error {
	r.members[member.ID()] = member
	return nil
}

func (r *memoryMemberRepository) Get(_ context.Context, gymID domain.GymID, memberID domain.MemberID) (domain.Member, error) {
	member, ok := r.members[memberID]
	if !ok || member.GymID() != gymID || r.archived[memberID] {
		return domain.Member{}, ports.ErrMemberNotFound
	}
	return member, nil
}

func (r *memoryMemberRepository) List(_ context.Context, gymID domain.GymID) ([]domain.Member, error) {
	var result []domain.Member
	for id, member := range r.members {
		if member.GymID() == gymID && !r.archived[id] {
			result = append(result, member)
		}
	}
	return result, nil
}

func (r *memoryMemberRepository) Update(_ context.Context, member domain.Member) error {
	if _, ok := r.members[member.ID()]; !ok || r.archived[member.ID()] {
		return ports.ErrMemberNotFound
	}
	r.members[member.ID()] = member
	return nil
}

func (r *memoryMemberRepository) Archive(_ context.Context, _ domain.GymID, memberID domain.MemberID, _ time.Time) error {
	if _, ok := r.members[memberID]; !ok || r.archived[memberID] {
		return ports.ErrMemberNotFound
	}
	if r.archived == nil {
		r.archived = map[domain.MemberID]bool{}
	}
	r.archived[memberID] = true
	return nil
}

func mustGymID(t *testing.T) domain.GymID {
	t.Helper()
	id, err := domain.NewGymID()
	if err != nil {
		t.Fatalf("NewGymID returned error: %v", err)
	}
	return id
}
