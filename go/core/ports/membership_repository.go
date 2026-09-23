package ports

import (
	"context"
	"errors"

	"github.com/zbango/gym-saas/go/core/domain"
)

var ErrMembershipNotFound = errors.New("membership not found")

// MembershipRepository persists purchased membership snapshots and their
// lifecycle updates without exposing SQLite details to application services.
type MembershipRepository interface {
	Create(context.Context, domain.Membership) error
	Get(context.Context, domain.GymID, domain.MembershipID) (domain.Membership, error)
	ListForMember(context.Context, domain.GymID, domain.MemberID) ([]domain.Membership, error)
	Update(context.Context, domain.Membership) error
}
