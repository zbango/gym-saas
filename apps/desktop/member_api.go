package main

import (
	"github.com/zbango/gym-saas/go/core/application"
	"github.com/zbango/gym-saas/go/core/domain"
)

// MemberAPI is the Wails delivery adapter for Member use cases.
type MemberAPI struct {
	runtime *DesktopRuntime
	members *application.MemberService
}

func NewMemberAPI(runtime *DesktopRuntime, members *application.MemberService) *MemberAPI {
	return &MemberAPI{runtime: runtime, members: members}
}

type MemberInput struct {
	FirstName            string `json:"firstName"`
	LastName             string `json:"lastName"`
	Email                string `json:"email"`
	Phone                string `json:"phone"`
	IdentificationNumber string `json:"identificationNumber"`
	DateOfBirth          string `json:"dateOfBirth"`
	Address              string `json:"address"`
	Status               string `json:"status"`
}

type Member struct {
	ID                   string `json:"id"`
	FirstName            string `json:"firstName"`
	LastName             string `json:"lastName"`
	Email                string `json:"email"`
	Phone                string `json:"phone"`
	IdentificationNumber string `json:"identificationNumber"`
	DateOfBirth          string `json:"dateOfBirth"`
	Address              string `json:"address"`
	Status               string `json:"status"`
}

func (a *MemberAPI) ListMembers() ([]Member, error) {
	members, err := a.members.List(a.runtime.requestContext())
	if err != nil {
		return nil, err
	}
	result := make([]Member, 0, len(members))
	for _, member := range members {
		result = append(result, memberView(member))
	}
	return result, nil
}

func (a *MemberAPI) CreateMember(input MemberInput) (Member, error) {
	member, err := a.members.Create(a.runtime.requestContext(), application.MemberInput(input))
	if err != nil {
		return Member{}, err
	}
	return memberView(member), nil
}

func (a *MemberAPI) UpdateMember(id string, input MemberInput) (Member, error) {
	member, err := a.members.Update(a.runtime.requestContext(), id, application.MemberInput(input))
	if err != nil {
		return Member{}, err
	}
	return memberView(member), nil
}

func (a *MemberAPI) ArchiveMember(id string) error {
	return a.members.Archive(a.runtime.requestContext(), id)
}

func memberView(member domain.Member) Member {
	return Member{
		ID:                   string(member.ID()),
		FirstName:            member.FirstName(),
		LastName:             member.LastName(),
		Email:                member.Email(),
		Phone:                member.Phone(),
		IdentificationNumber: member.IdentificationNumber(),
		DateOfBirth:          member.DateOfBirth(),
		Address:              member.Address(),
		Status:               string(member.Status()),
	}
}
