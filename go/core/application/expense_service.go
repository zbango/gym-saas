package application

import (
	"context"
	"fmt"
	"time"

	"github.com/zbango/gym-saas/go/core/domain"
	"github.com/zbango/gym-saas/go/core/ports"
)

type RecordExpenseInput struct {
	AmountCents   int64
	Currency      string
	PaymentMethod string
	Reference     string
	Notes         string
}

// ExpenseService coordinates standalone outgoing-expense lifecycle operations
// for one gym.
type ExpenseService struct {
	expenses ports.ExpenseRepository
	gymID    domain.GymID
	now      func() time.Time
}

func NewExpenseService(expenses ports.ExpenseRepository, gymID domain.GymID, now func() time.Time) (*ExpenseService, error) {
	if expenses == nil {
		return nil, fmt.Errorf("expense repository is required")
	}
	if _, err := domain.ParseGymID(string(gymID)); err != nil {
		return nil, err
	}
	if now == nil {
		now = time.Now
	}
	return &ExpenseService{expenses: expenses, gymID: gymID, now: now}, nil
}

func (s *ExpenseService) Record(ctx context.Context, input RecordExpenseInput) (domain.Expense, error) {
	expense, err := s.expenseFromInput(input, false)
	if err != nil {
		return domain.Expense{}, err
	}
	if err := s.expenses.Create(ctx, expense); err != nil {
		return domain.Expense{}, fmt.Errorf("save expense: %w", err)
	}
	return expense, nil
}

func (s *ExpenseService) CreatePending(ctx context.Context, input RecordExpenseInput) (domain.Expense, error) {
	expense, err := s.expenseFromInput(input, true)
	if err != nil {
		return domain.Expense{}, err
	}
	if err := s.expenses.Create(ctx, expense); err != nil {
		return domain.Expense{}, fmt.Errorf("save pending expense: %w", err)
	}
	return expense, nil
}

func (s *ExpenseService) Void(ctx context.Context, id string) (domain.Expense, error) {
	expense, err := s.find(ctx, id)
	if err != nil {
		return domain.Expense{}, err
	}
	updated, err := expense.Void(s.now())
	if err != nil {
		return domain.Expense{}, err
	}
	if err := s.expenses.Update(ctx, updated); err != nil {
		return domain.Expense{}, fmt.Errorf("void expense: %w", err)
	}
	return updated, nil
}

func (s *ExpenseService) Post(ctx context.Context, id string) (domain.Expense, error) {
	expense, err := s.find(ctx, id)
	if err != nil {
		return domain.Expense{}, err
	}
	updated, err := expense.Post(s.now())
	if err != nil {
		return domain.Expense{}, err
	}
	if err := s.expenses.Update(ctx, updated); err != nil {
		return domain.Expense{}, fmt.Errorf("post expense: %w", err)
	}
	return updated, nil
}

func (s *ExpenseService) Get(ctx context.Context, id string) (domain.Expense, error) {
	return s.find(ctx, id)
}

func (s *ExpenseService) List(ctx context.Context) ([]domain.Expense, error) {
	expenses, err := s.expenses.List(ctx, s.gymID)
	if err != nil {
		return nil, fmt.Errorf("list expenses: %w", err)
	}
	return expenses, nil
}

func (s *ExpenseService) expenseFromInput(input RecordExpenseInput, pending bool) (domain.Expense, error) {
	amount, err := domain.NewMoney(input.AmountCents, input.Currency)
	if err != nil {
		return domain.Expense{}, err
	}
	method, err := domain.ParsePaymentMethod(input.PaymentMethod)
	if err != nil {
		return domain.Expense{}, err
	}
	if pending {
		return domain.CreatePendingExpense(s.gymID, amount, method, input.Reference, input.Notes, s.now())
	}
	return domain.RecordExpense(s.gymID, amount, method, input.Reference, input.Notes, s.now())
}

func (s *ExpenseService) find(ctx context.Context, id string) (domain.Expense, error) {
	expenseID, err := domain.ParseExpenseID(id)
	if err != nil {
		return domain.Expense{}, err
	}
	expense, err := s.expenses.Get(ctx, s.gymID, expenseID)
	if err != nil {
		return domain.Expense{}, fmt.Errorf("find expense: %w", err)
	}
	return expense, nil
}
