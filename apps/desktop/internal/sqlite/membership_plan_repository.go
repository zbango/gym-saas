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

type MembershipPlanRepository struct {
	db *sql.DB
}

func NewMembershipPlanRepository(store *Store) *MembershipPlanRepository {
	return &MembershipPlanRepository{db: store.db}
}

func (r *MembershipPlanRepository) Create(ctx context.Context, plan domain.MembershipPlan) error {
	_, err := r.db.ExecContext(ctx, `
INSERT INTO membership_plans (
  id, gym_id, name, validity_kind, duration_value, duration_unit, visit_limit,
  price_cents, currency, status, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`, plan.ID(), plan.GymID(), plan.Name(), plan.ValidityKind(), plan.DurationValue(), plan.DurationUnit(), nullIntIfZero(plan.VisitLimit()), plan.Price().Cents(), plan.Price().Currency(), plan.Status(), domain.FormatTimestamp(plan.CreatedAt()), domain.FormatTimestamp(plan.UpdatedAt()))
	if err != nil {
		return fmt.Errorf("insert membership plan: %w", err)
	}
	return nil
}

func (r *MembershipPlanRepository) Get(ctx context.Context, gymID domain.GymID, planID domain.MembershipPlanID) (domain.MembershipPlan, error) {
	plan, err := scanMembershipPlan(r.db.QueryRowContext(ctx, membershipPlanSelect+` WHERE gym_id = ? AND id = ? AND deleted_at IS NULL`, gymID, planID))
	if errors.Is(err, sql.ErrNoRows) {
		return domain.MembershipPlan{}, ports.ErrMembershipPlanNotFound
	}
	if err != nil {
		return domain.MembershipPlan{}, fmt.Errorf("get membership plan: %w", err)
	}
	return plan, nil
}

func (r *MembershipPlanRepository) List(ctx context.Context, gymID domain.GymID) ([]domain.MembershipPlan, error) {
	rows, err := r.db.QueryContext(ctx, membershipPlanSelect+` WHERE gym_id = ? AND deleted_at IS NULL ORDER BY display_order, created_at DESC`, gymID)
	if err != nil {
		return nil, fmt.Errorf("list membership plans: %w", err)
	}
	defer rows.Close()

	plans := []domain.MembershipPlan{}
	for rows.Next() {
		plan, err := scanMembershipPlan(rows)
		if err != nil {
			return nil, fmt.Errorf("scan membership plan: %w", err)
		}
		plans = append(plans, plan)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate membership plans: %w", err)
	}
	return plans, nil
}

func (r *MembershipPlanRepository) Update(ctx context.Context, plan domain.MembershipPlan) error {
	result, err := r.db.ExecContext(ctx, `
UPDATE membership_plans
SET name = ?, validity_kind = ?, duration_value = ?, duration_unit = ?, visit_limit = ?,
    price_cents = ?, currency = ?, status = ?, updated_at = ?, version = version + 1
WHERE gym_id = ? AND id = ? AND deleted_at IS NULL
`, plan.Name(), plan.ValidityKind(), plan.DurationValue(), plan.DurationUnit(), nullIntIfZero(plan.VisitLimit()), plan.Price().Cents(), plan.Price().Currency(), plan.Status(), domain.FormatTimestamp(plan.UpdatedAt()), plan.GymID(), plan.ID())
	if err != nil {
		return fmt.Errorf("update membership plan: %w", err)
	}
	if affected, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("update membership plan rows affected: %w", err)
	} else if affected == 0 {
		return ports.ErrMembershipPlanNotFound
	}
	return nil
}

func (r *MembershipPlanRepository) Archive(ctx context.Context, gymID domain.GymID, planID domain.MembershipPlanID, archivedAt time.Time) error {
	archivedAt = archivedAt.UTC()
	result, err := r.db.ExecContext(ctx, `
UPDATE membership_plans
SET status = 'archived', deleted_at = ?, updated_at = ?, version = version + 1
WHERE gym_id = ? AND id = ? AND deleted_at IS NULL
`, domain.FormatTimestamp(archivedAt), domain.FormatTimestamp(archivedAt), gymID, planID)
	if err != nil {
		return fmt.Errorf("archive membership plan: %w", err)
	}
	if affected, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("archive membership plan rows affected: %w", err)
	} else if affected == 0 {
		return ports.ErrMembershipPlanNotFound
	}
	return nil
}

const membershipPlanSelect = `
SELECT id, gym_id, name, validity_kind, duration_value, duration_unit, visit_limit,
       price_cents, currency, status, created_at, updated_at
FROM membership_plans`

func scanMembershipPlan(row rowScanner) (domain.MembershipPlan, error) {
	var id, gymID, name, validityKind, durationUnit, currency, status, createdAt, updatedAt string
	var durationValue int
	var priceCents int64
	var visitLimit sql.NullInt64
	if err := row.Scan(&id, &gymID, &name, &validityKind, &durationValue, &durationUnit, &visitLimit, &priceCents, &currency, &status, &createdAt, &updatedAt); err != nil {
		return domain.MembershipPlan{}, err
	}
	planID, err := domain.ParseMembershipPlanID(id)
	if err != nil {
		return domain.MembershipPlan{}, fmt.Errorf("parse membership plan ID: %w", err)
	}
	ownerID, err := domain.ParseGymID(gymID)
	if err != nil {
		return domain.MembershipPlan{}, fmt.Errorf("parse gym ID: %w", err)
	}
	price, err := domain.NewMoney(priceCents, currency)
	if err != nil {
		return domain.MembershipPlan{}, err
	}
	created, err := domain.ParseTimestamp(createdAt)
	if err != nil {
		return domain.MembershipPlan{}, err
	}
	updated, err := domain.ParseTimestamp(updatedAt)
	if err != nil {
		return domain.MembershipPlan{}, err
	}
	return domain.NewMembershipPlan(planID, ownerID, name, domain.MembershipValidityKind(validityKind), durationValue, domain.MembershipDurationUnit(durationUnit), int(visitLimit.Int64), price, domain.MembershipPlanStatus(status), created, updated)
}

func nullIntIfZero(value int) any {
	if value == 0 {
		return nil
	}
	return value
}
