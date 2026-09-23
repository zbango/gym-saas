package ports

import (
	"context"

	"github.com/zbango/gym-saas/go/core/domain"
)

// MembershipPaymentWriter persists a purchased membership and its initial
// receipt atomically. It is purpose-specific because generic transaction APIs
// would leak storage concerns into application code.
type MembershipPaymentWriter interface {
	CreateMembershipAndPayment(context.Context, domain.Membership, domain.Payment) error
}
