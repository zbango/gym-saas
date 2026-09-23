package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/zbango/gym-saas/go/core/domain"
	"github.com/zbango/gym-saas/go/core/ports"
)

// ExpenseRepository persists outgoing business expenses. The original amount
// and settlement metadata are immutable once created.
type ExpenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(store *Store) *ExpenseRepository {
	return &ExpenseRepository{db: store.db}
}

func (r *ExpenseRepository) Create(ctx context.Context, expense domain.Expense) error {
	_, err := r.db.ExecContext(ctx, `
INSERT INTO expenses (
  id, gym_id, status, amount_cents, currency, payment_method, reference, notes,
  paid_at, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`, expense.ID(), expense.GymID(), expense.Status(), expense.Amount().Cents(), expense.Amount().Currency(), expense.Method(), nullIfEmpty(expense.Reference()), nullIfEmpty(expense.Notes()), nullTimestamp(expense.PaidAt()), domain.FormatTimestamp(expense.CreatedAt()), domain.FormatTimestamp(expense.UpdatedAt()))
	if err != nil {
		return fmt.Errorf("insert expense: %w", err)
	}
	return nil
}

func (r *ExpenseRepository) Get(ctx context.Context, gymID domain.GymID, expenseID domain.ExpenseID) (domain.Expense, error) {
	expense, err := scanExpense(r.db.QueryRowContext(ctx, expenseSelect+` WHERE gym_id = ? AND id = ?`, gymID, expenseID))
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Expense{}, ports.ErrExpenseNotFound
	}
	if err != nil {
		return domain.Expense{}, fmt.Errorf("get expense: %w", err)
	}
	return expense, nil
}

func (r *ExpenseRepository) List(ctx context.Context, gymID domain.GymID) ([]domain.Expense, error) {
	rows, err := r.db.QueryContext(ctx, expenseSelect+` WHERE gym_id = ? ORDER BY created_at DESC`, gymID)
	if err != nil {
		return nil, fmt.Errorf("list expenses: %w", err)
	}
	defer rows.Close()

	expenses := []domain.Expense{}
	for rows.Next() {
		expense, err := scanExpense(rows)
		if err != nil {
			return nil, fmt.Errorf("scan expense: %w", err)
		}
		expenses = append(expenses, expense)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate expenses: %w", err)
	}
	return expenses, nil
}

func (r *ExpenseRepository) Update(ctx context.Context, expense domain.Expense) error {
	result, err := r.db.ExecContext(ctx, `
UPDATE expenses
SET status = ?, paid_at = ?, updated_at = ?, version = version + 1
WHERE gym_id = ? AND id = ?
`, expense.Status(), nullTimestamp(expense.PaidAt()), domain.FormatTimestamp(expense.UpdatedAt()), expense.GymID(), expense.ID())
	if err != nil {
		return fmt.Errorf("update expense: %w", err)
	}
	if affected, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("update expense rows affected: %w", err)
	} else if affected == 0 {
		return ports.ErrExpenseNotFound
	}
	return nil
}

const expenseSelect = `
SELECT id, gym_id, status, amount_cents, currency, payment_method, reference, notes,
       paid_at, created_at, updated_at
FROM expenses`

func scanExpense(row rowScanner) (domain.Expense, error) {
	var id, gymID, status, currency, method, createdAt, updatedAt string
	var reference, notes, paidAt sql.NullString
	var amountCents int64
	if err := row.Scan(&id, &gymID, &status, &amountCents, &currency, &method, &reference, &notes, &paidAt, &createdAt, &updatedAt); err != nil {
		return domain.Expense{}, err
	}
	expenseID, err := domain.ParseExpenseID(id)
	if err != nil {
		return domain.Expense{}, fmt.Errorf("parse expense ID: %w", err)
	}
	ownerID, err := domain.ParseGymID(gymID)
	if err != nil {
		return domain.Expense{}, fmt.Errorf("parse gym ID: %w", err)
	}
	amount, err := domain.NewMoney(amountCents, currency)
	if err != nil {
		return domain.Expense{}, err
	}
	paid, err := parseOptionalTimestamp(paidAt)
	if err != nil {
		return domain.Expense{}, err
	}
	created, err := domain.ParseTimestamp(createdAt)
	if err != nil {
		return domain.Expense{}, err
	}
	updated, err := domain.ParseTimestamp(updatedAt)
	if err != nil {
		return domain.Expense{}, err
	}
	return domain.NewExpense(expenseID, ownerID, domain.ExpenseStatus(status), amount, domain.PaymentMethod(method), reference.String, notes.String, paid, created, updated)
}
