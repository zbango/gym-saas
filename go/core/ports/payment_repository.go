package ports

import (
	"context"
	"errors"

	"github.com/zbango/gym-saas/go/core/domain"
)

var ErrPaymentNotFound = errors.New("payment not found")

// PaymentRepository persists member-side receipts and their lifecycle updates.
type PaymentRepository interface {
	Create(context.Context, domain.Payment) error
	Get(context.Context, domain.GymID, domain.PaymentID) (domain.Payment, error)
	ListForMember(context.Context, domain.GymID, domain.MemberID) ([]domain.Payment, error)
	Update(context.Context, domain.Payment) error
}
