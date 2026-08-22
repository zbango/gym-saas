package ports

import (
	"context"
	"errors"
	"time"

	"github.com/zbango/gym-saas/go/core/domain"
)

var ErrMembershipPlanNotFound = errors.New("membership plan not found")

type MembershipPlanRepository interface {
	Create(context.Context, domain.MembershipPlan) error
	Get(context.Context, domain.GymID, domain.MembershipPlanID) (domain.MembershipPlan, error)
	List(context.Context, domain.GymID) ([]domain.MembershipPlan, error)
	Update(context.Context, domain.MembershipPlan) error
	Archive(context.Context, domain.GymID, domain.MembershipPlanID, time.Time) error
}
