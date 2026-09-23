package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/zbango/gym-saas/go/core/domain"
	"github.com/zbango/gym-saas/go/core/ports"
)

func TestMembershipServiceStartsActivatesCancelsAndListsMembership(t *testing.T) {
	gym, member, plan := membershipServiceFixture(t)
	members := &memoryMemberRepository{members: map[domain.MemberID]domain.Member{member.ID(): member}}
	plans := &memoryMembershipPlanRepository{plans: map[domain.MembershipPlanID]domain.MembershipPlan{plan.ID(): plan}}
	memberships := &memoryMembershipRepository{memberships: map[domain.MembershipID]domain.Membership{}}
	times := []time.Time{
		time.Date(2026, time.August, 1, 9, 0, 0, 0, time.UTC),
		time.Date(2026, time.August, 8, 9, 0, 0, 0, time.UTC),
		time.Date(2026, time.August, 9, 9, 0, 0, 0, time.UTC),
	}
	service, err := NewMembershipService(members, plans, memberships, gym, func() time.Time {
		value := times[0]
		times = times[1:]
		return value
	})
	if err != nil {
		t.Fatalf("NewMembershipService returned error: %v", err)
	}

	created, err := service.Start(context.Background(), StartMembershipInput{
		MemberID:         string(member.ID()),
		MembershipPlanID: string(plan.ID()),
		StartsAt:         time.Date(2026, time.August, 8, 9, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("Start returned error: %v", err)
	}
	if created.Status() != domain.MembershipStatusPending || len(memberships.memberships) != 1 {
		t.Fatalf("created membership = %#v", created)
	}

	activated, err := service.Activate(context.Background(), string(created.ID()))
	if err != nil {
		t.Fatalf("Activate returned error: %v", err)
	}
	if activated.Status() != domain.MembershipStatusActive || !activated.IsValidAt(activated.ActivatedAt()) {
		t.Fatalf("activated membership = %#v", activated)
	}
	cancelled, err := service.Cancel(context.Background(), string(created.ID()), "member request")
	if err != nil {
		t.Fatalf("Cancel returned error: %v", err)
	}
	if cancelled.Status() != domain.MembershipStatusCancelled || cancelled.CancellationReason() != "member request" {
		t.Fatalf("cancelled membership = %#v", cancelled)
	}

	listed, err := service.ListForMember(context.Background(), string(member.ID()))
	if err != nil {
		t.Fatalf("ListForMember returned error: %v", err)
	}
	if len(listed) != 1 || listed[0].ID() != created.ID() || listed[0].Status() != domain.MembershipStatusCancelled {
		t.Fatalf("ListForMember = %#v", listed)
	}
}

func TestMembershipServiceRejectsUnknownMemberBeforePersistence(t *testing.T) {
	gym, _, plan := membershipServiceFixture(t)
	members := &memoryMemberRepository{members: map[domain.MemberID]domain.Member{}}
	plans := &memoryMembershipPlanRepository{plans: map[domain.MembershipPlanID]domain.MembershipPlan{plan.ID(): plan}}
	memberships := &memoryMembershipRepository{memberships: map[domain.MembershipID]domain.Membership{}}
	service, err := NewMembershipService(members, plans, memberships, gym, time.Now)
	if err != nil {
		t.Fatalf("NewMembershipService returned error: %v", err)
	}
	memberID, err := domain.NewMemberID()
	if err != nil {
		t.Fatalf("NewMemberID returned error: %v", err)
	}
	_, err = service.Start(context.Background(), StartMembershipInput{MemberID: string(memberID), MembershipPlanID: string(plan.ID())})
	if !errors.Is(err, ports.ErrMemberNotFound) {
		t.Fatalf("Start error = %v, want %v", err, ports.ErrMemberNotFound)
	}
	if len(memberships.memberships) != 0 {
		t.Fatal("unknown member produced a membership")
	}
}

func membershipServiceFixture(t *testing.T) (domain.Gym, domain.Member, domain.MembershipPlan) {
	t.Helper()
	now := time.Date(2026, time.August, 1, 8, 0, 0, 0, time.UTC)
	gym, err := domain.CreateGym("Zeus", "America/Guayaquil", now)
	if err != nil {
		t.Fatalf("CreateGym returned error: %v", err)
	}
	member, err := domain.CreateMember(gym.ID(), "Ada", "Lovelace", "ada@example.com", "555-0100", "0102030405", "1815-12-10", "", domain.MemberStatusActive, now)
	if err != nil {
		t.Fatalf("CreateMember returned error: %v", err)
	}
	price, err := domain.NewMoney(4500, "USD")
	if err != nil {
		t.Fatalf("NewMoney returned error: %v", err)
	}
	plan, err := domain.CreateMembershipPlan(gym.ID(), "Monthly", domain.MembershipValidityTime, 1, domain.MembershipDurationMonths, 0, price, domain.MembershipPlanStatusActive, now)
	if err != nil {
		t.Fatalf("CreateMembershipPlan returned error: %v", err)
	}
	return gym, member, plan
}

type memoryMembershipRepository struct {
	memberships map[domain.MembershipID]domain.Membership
}

func (r *memoryMembershipRepository) Create(_ context.Context, membership domain.Membership) error {
	r.memberships[membership.ID()] = membership
	return nil
}

func (r *memoryMembershipRepository) Get(_ context.Context, gymID domain.GymID, membershipID domain.MembershipID) (domain.Membership, error) {
	membership, ok := r.memberships[membershipID]
	if !ok || membership.GymID() != gymID {
		return domain.Membership{}, ports.ErrMembershipNotFound
	}
	return membership, nil
}

func (r *memoryMembershipRepository) ListForMember(_ context.Context, gymID domain.GymID, memberID domain.MemberID) ([]domain.Membership, error) {
	var result []domain.Membership
	for _, membership := range r.memberships {
		if membership.GymID() == gymID && membership.MemberID() == memberID {
			result = append(result, membership)
		}
	}
	return result, nil
}

func (r *memoryMembershipRepository) Update(_ context.Context, membership domain.Membership) error {
	if _, ok := r.memberships[membership.ID()]; !ok {
		return ports.ErrMembershipNotFound
	}
	r.memberships[membership.ID()] = membership
	return nil
}
