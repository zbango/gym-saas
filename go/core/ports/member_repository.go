package ports

import (
	"context"
	"errors"
	"time"

	"github.com/zbango/gym-saas/go/core/domain"
)

var ErrMemberNotFound = errors.New("member not found")

// MemberRepository persists and retrieves Members without exposing storage
// concerns to domain rules.
type MemberRepository interface {
	Create(context.Context, domain.Member) error
	Get(context.Context, domain.GymID, domain.MemberID) (domain.Member, error)
	List(context.Context, domain.GymID) ([]domain.Member, error)
	Update(context.Context, domain.Member) error
	Archive(context.Context, domain.GymID, domain.MemberID, time.Time) error
}
