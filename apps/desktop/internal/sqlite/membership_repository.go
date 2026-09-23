package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/zbango/gym-saas/go/core/domain"
	"github.com/zbango/gym-saas/go/core/ports"
)

// MembershipRepository stores purchased plan snapshots and the lifecycle
// fields that change after purchase. Plan terms are immutable once inserted.
type MembershipRepository struct {
	db *sql.DB
}

func NewMembershipRepository(store *Store) *MembershipRepository {
	return &MembershipRepository{db: store.db}
}

func (r *MembershipRepository) Create(ctx context.Context, membership domain.Membership) error {
	if err := insertMembership(ctx, r.db, membership); err != nil {
		return fmt.Errorf("insert membership: %w", err)
	}
	return nil
}

func insertMembership(ctx context.Context, executor sqlExecutor, membership domain.Membership) error {
	_, err := executor.ExecContext(ctx, `
INSERT INTO memberships (
  id, gym_id, member_id, membership_plan_id, status, starts_at, ends_at,
  activated_at, cancelled_at, cancellation_reason, validity_kind_snapshot,
  duration_value_snapshot, duration_unit_snapshot, visit_limit_snapshot,
  visits_remaining, price_cents_snapshot, currency_snapshot, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`, membership.ID(), membership.GymID(), membership.MemberID(), membership.MembershipPlanID(), membership.Status(), domain.FormatTimestamp(membership.StartsAt()), domain.FormatTimestamp(membership.EndsAt()), nullTimestamp(membership.ActivatedAt()), nullTimestamp(membership.CancelledAt()), nullIfEmpty(membership.CancellationReason()), membership.ValidityKind(), membership.DurationValue(), membership.DurationUnit(), nullableVisitLimit(membership), nullableVisitsRemaining(membership), membership.Price().Cents(), membership.Price().Currency(), domain.FormatTimestamp(membership.CreatedAt()), domain.FormatTimestamp(membership.UpdatedAt()))
	return err
}

func (r *MembershipRepository) Get(ctx context.Context, gymID domain.GymID, membershipID domain.MembershipID) (domain.Membership, error) {
	membership, err := scanMembership(r.db.QueryRowContext(ctx, membershipSelect+` WHERE gym_id = ? AND id = ?`, gymID, membershipID))
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Membership{}, ports.ErrMembershipNotFound
	}
	if err != nil {
		return domain.Membership{}, fmt.Errorf("get membership: %w", err)
	}
	return membership, nil
}

func (r *MembershipRepository) ListForMember(ctx context.Context, gymID domain.GymID, memberID domain.MemberID) ([]domain.Membership, error) {
	rows, err := r.db.QueryContext(ctx, membershipSelect+` WHERE gym_id = ? AND member_id = ? ORDER BY starts_at DESC, created_at DESC`, gymID, memberID)
	if err != nil {
		return nil, fmt.Errorf("list member memberships: %w", err)
	}
	defer rows.Close()

	memberships := []domain.Membership{}
	for rows.Next() {
		membership, err := scanMembership(rows)
		if err != nil {
			return nil, fmt.Errorf("scan membership: %w", err)
		}
		memberships = append(memberships, membership)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate member memberships: %w", err)
	}
	return memberships, nil
}

func (r *MembershipRepository) Update(ctx context.Context, membership domain.Membership) error {
	result, err := r.db.ExecContext(ctx, `
UPDATE memberships
SET status = ?, activated_at = ?, cancelled_at = ?, cancellation_reason = ?,
    visits_remaining = ?, updated_at = ?, version = version + 1
WHERE gym_id = ? AND id = ?
`, membership.Status(), nullTimestamp(membership.ActivatedAt()), nullTimestamp(membership.CancelledAt()), nullIfEmpty(membership.CancellationReason()), nullableVisitsRemaining(membership), domain.FormatTimestamp(membership.UpdatedAt()), membership.GymID(), membership.ID())
	if err != nil {
		return fmt.Errorf("update membership: %w", err)
	}
	if affected, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("update membership rows affected: %w", err)
	} else if affected == 0 {
		return ports.ErrMembershipNotFound
	}
	return nil
}

const membershipSelect = `
SELECT id, gym_id, member_id, membership_plan_id, status, starts_at, ends_at,
       activated_at, cancelled_at, cancellation_reason, validity_kind_snapshot,
       duration_value_snapshot, duration_unit_snapshot, visit_limit_snapshot,
       visits_remaining, price_cents_snapshot, currency_snapshot, created_at, updated_at
FROM memberships`

func scanMembership(row rowScanner) (domain.Membership, error) {
	var id, gymID, memberID, planID, status, startsAt, endsAt string
	var validityKind, durationUnit, currency, createdAt, updatedAt string
	var activatedAt, cancelledAt, cancellationReason sql.NullString
	var durationValue int
	var visitLimit, visitsRemaining sql.NullInt64
	var priceCents int64
	if err := row.Scan(
		&id, &gymID, &memberID, &planID, &status, &startsAt, &endsAt,
		&activatedAt, &cancelledAt, &cancellationReason, &validityKind,
		&durationValue, &durationUnit, &visitLimit, &visitsRemaining,
		&priceCents, &currency, &createdAt, &updatedAt,
	); err != nil {
		return domain.Membership{}, err
	}
	membershipID, err := domain.ParseMembershipID(id)
	if err != nil {
		return domain.Membership{}, fmt.Errorf("parse membership ID: %w", err)
	}
	ownerID, err := domain.ParseGymID(gymID)
	if err != nil {
		return domain.Membership{}, fmt.Errorf("parse gym ID: %w", err)
	}
	member, err := domain.ParseMemberID(memberID)
	if err != nil {
		return domain.Membership{}, fmt.Errorf("parse member ID: %w", err)
	}
	plan, err := domain.ParseMembershipPlanID(planID)
	if err != nil {
		return domain.Membership{}, fmt.Errorf("parse membership plan ID: %w", err)
	}
	starts, err := domain.ParseTimestamp(startsAt)
	if err != nil {
		return domain.Membership{}, err
	}
	ends, err := domain.ParseTimestamp(endsAt)
	if err != nil {
		return domain.Membership{}, err
	}
	activated, err := parseOptionalTimestamp(activatedAt)
	if err != nil {
		return domain.Membership{}, err
	}
	cancelled, err := parseOptionalTimestamp(cancelledAt)
	if err != nil {
		return domain.Membership{}, err
	}
	price, err := domain.NewMoney(priceCents, currency)
	if err != nil {
		return domain.Membership{}, err
	}
	created, err := domain.ParseTimestamp(createdAt)
	if err != nil {
		return domain.Membership{}, err
	}
	updated, err := domain.ParseTimestamp(updatedAt)
	if err != nil {
		return domain.Membership{}, err
	}
	return domain.NewMembership(
		membershipID, ownerID, member, plan, domain.MembershipStatus(status), starts, ends,
		activated, cancelled, cancellationReason.String, domain.MembershipValidityKind(validityKind),
		durationValue, domain.MembershipDurationUnit(durationUnit), int(visitLimit.Int64), int(visitsRemaining.Int64),
		price, created, updated,
	)
}

func parseOptionalTimestamp(value sql.NullString) (time.Time, error) {
	if !value.Valid {
		return time.Time{}, nil
	}
	parsed, err := domain.ParseTimestamp(value.String)
	if err != nil {
		return time.Time{}, err
	}
	return parsed, nil
}

func nullTimestamp(value time.Time) any {
	if value.IsZero() {
		return nil
	}
	return domain.FormatTimestamp(value)
}

func nullableVisitLimit(membership domain.Membership) any {
	if membership.ValidityKind() == domain.MembershipValidityTime {
		return nil
	}
	return membership.VisitLimit()
}

func nullableVisitsRemaining(membership domain.Membership) any {
	if membership.ValidityKind() == domain.MembershipValidityTime {
		return nil
	}
	return membership.VisitsRemaining()
}
