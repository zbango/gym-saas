package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/zbango/gym-saas/go/core/domain"
	"github.com/zbango/gym-saas/go/core/ports"
)

// PaymentRepository persists member receipts. Monetary details and references
// are immutable after creation; lifecycle status is the only update surface.
type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(store *Store) *PaymentRepository {
	return &PaymentRepository{db: store.db}
}

func (r *PaymentRepository) Create(ctx context.Context, payment domain.Payment) error {
	if err := insertPayment(ctx, r.db, payment); err != nil {
		return fmt.Errorf("insert payment: %w", err)
	}
	return nil
}

func insertPayment(ctx context.Context, executor sqlExecutor, payment domain.Payment) error {
	_, err := executor.ExecContext(ctx, `
INSERT INTO payments (
  id, gym_id, member_id, membership_id, status, kind, amount_cents, currency,
  payment_method, reference, notes, paid_at, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`, payment.ID(), payment.GymID(), payment.MemberID(), nullMembershipID(payment.MembershipID()), payment.Status(), payment.Kind(), payment.Amount().Cents(), payment.Amount().Currency(), payment.Method(), nullIfEmpty(payment.Reference()), nullIfEmpty(payment.Notes()), nullTimestamp(payment.PaidAt()), domain.FormatTimestamp(payment.CreatedAt()), domain.FormatTimestamp(payment.UpdatedAt()))
	return err
}

func (r *PaymentRepository) Get(ctx context.Context, gymID domain.GymID, paymentID domain.PaymentID) (domain.Payment, error) {
	payment, err := scanPayment(r.db.QueryRowContext(ctx, paymentSelect+` WHERE gym_id = ? AND id = ?`, gymID, paymentID))
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Payment{}, ports.ErrPaymentNotFound
	}
	if err != nil {
		return domain.Payment{}, fmt.Errorf("get payment: %w", err)
	}
	return payment, nil
}

func (r *PaymentRepository) ListForMember(ctx context.Context, gymID domain.GymID, memberID domain.MemberID) ([]domain.Payment, error) {
	rows, err := r.db.QueryContext(ctx, paymentSelect+` WHERE gym_id = ? AND member_id = ? ORDER BY created_at DESC`, gymID, memberID)
	if err != nil {
		return nil, fmt.Errorf("list member payments: %w", err)
	}
	defer rows.Close()

	payments := []domain.Payment{}
	for rows.Next() {
		payment, err := scanPayment(rows)
		if err != nil {
			return nil, fmt.Errorf("scan payment: %w", err)
		}
		payments = append(payments, payment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate member payments: %w", err)
	}
	return payments, nil
}

func (r *PaymentRepository) Update(ctx context.Context, payment domain.Payment) error {
	result, err := r.db.ExecContext(ctx, `
UPDATE payments
SET status = ?, paid_at = ?, updated_at = ?
WHERE gym_id = ? AND id = ?
`, payment.Status(), nullTimestamp(payment.PaidAt()), domain.FormatTimestamp(payment.UpdatedAt()), payment.GymID(), payment.ID())
	if err != nil {
		return fmt.Errorf("update payment: %w", err)
	}
	if affected, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("update payment rows affected: %w", err)
	} else if affected == 0 {
		return ports.ErrPaymentNotFound
	}
	return nil
}

const paymentSelect = `
SELECT id, gym_id, member_id, membership_id, status, kind, amount_cents, currency,
       payment_method, reference, notes, paid_at, created_at, updated_at
FROM payments`

func scanPayment(row rowScanner) (domain.Payment, error) {
	var id, gymID, memberID, status, kind, currency, method, createdAt, updatedAt string
	var membershipID, reference, notes, paidAt sql.NullString
	var amountCents int64
	if err := row.Scan(&id, &gymID, &memberID, &membershipID, &status, &kind, &amountCents, &currency, &method, &reference, &notes, &paidAt, &createdAt, &updatedAt); err != nil {
		return domain.Payment{}, err
	}
	paymentID, err := domain.ParsePaymentID(id)
	if err != nil {
		return domain.Payment{}, fmt.Errorf("parse payment ID: %w", err)
	}
	ownerID, err := domain.ParseGymID(gymID)
	if err != nil {
		return domain.Payment{}, fmt.Errorf("parse gym ID: %w", err)
	}
	member, err := domain.ParseMemberID(memberID)
	if err != nil {
		return domain.Payment{}, fmt.Errorf("parse member ID: %w", err)
	}
	var membership domain.MembershipID
	if membershipID.Valid {
		membership, err = domain.ParseMembershipID(membershipID.String)
		if err != nil {
			return domain.Payment{}, fmt.Errorf("parse membership ID: %w", err)
		}
	}
	amount, err := domain.NewMoney(amountCents, currency)
	if err != nil {
		return domain.Payment{}, err
	}
	paid, err := parseOptionalTimestamp(paidAt)
	if err != nil {
		return domain.Payment{}, err
	}
	created, err := domain.ParseTimestamp(createdAt)
	if err != nil {
		return domain.Payment{}, err
	}
	updated, err := domain.ParseTimestamp(updatedAt)
	if err != nil {
		return domain.Payment{}, err
	}
	return domain.NewPayment(paymentID, ownerID, member, membership, domain.PaymentStatus(status), domain.PaymentKind(kind), amount, domain.PaymentMethod(method), reference.String, notes.String, paid, created, updated)
}

func nullMembershipID(value domain.MembershipID) any {
	if value == "" {
		return nil
	}
	return value
}
