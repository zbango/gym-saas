package sqlite

import (
	"context"
	"fmt"

	"github.com/zbango/gym-saas/go/core/domain"
)

// MembershipPaymentWriter is the SQLite implementation of the purchase
// transaction boundary. It intentionally exposes one use-case operation
// rather than a generic transaction to application code.
type MembershipPaymentWriter struct {
	store *Store
}

func NewMembershipPaymentWriter(store *Store) *MembershipPaymentWriter {
	return &MembershipPaymentWriter{store: store}
}

func (w *MembershipPaymentWriter) CreateMembershipAndPayment(ctx context.Context, membership domain.Membership, payment domain.Payment) error {
	if payment.GymID() != membership.GymID() || payment.MemberID() != membership.MemberID() || payment.MembershipID() != membership.ID() {
		return fmt.Errorf("membership payment references do not match")
	}
	tx, err := w.store.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin membership payment transaction: %w", err)
	}
	defer tx.Rollback()

	if err := insertMembership(ctx, tx, membership); err != nil {
		return fmt.Errorf("insert membership: %w", err)
	}
	if err := insertPayment(ctx, tx, payment); err != nil {
		return fmt.Errorf("insert payment: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit membership payment transaction: %w", err)
	}
	return nil
}
