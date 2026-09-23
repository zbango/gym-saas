package ports

import (
	"context"
	"errors"

	"github.com/zbango/gym-saas/go/core/domain"
)

var ErrExpenseNotFound = errors.New("expense not found")

// ExpenseRepository persists outgoing business expenses and lifecycle updates.
type ExpenseRepository interface {
	Create(context.Context, domain.Expense) error
	Get(context.Context, domain.GymID, domain.ExpenseID) (domain.Expense, error)
	List(context.Context, domain.GymID) ([]domain.Expense, error)
	Update(context.Context, domain.Expense) error
}
