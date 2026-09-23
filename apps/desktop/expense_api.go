package main

import (
	"github.com/zbango/gym-saas/go/core/application"
	"github.com/zbango/gym-saas/go/core/domain"
)

// ExpenseAPI is the Wails delivery adapter for business-expense lifecycle use
// cases.
type ExpenseAPI struct {
	runtime  *DesktopRuntime
	expenses *application.ExpenseService
}

func NewExpenseAPI(runtime *DesktopRuntime, expenses *application.ExpenseService) *ExpenseAPI {
	return &ExpenseAPI{runtime: runtime, expenses: expenses}
}

type ExpenseInput struct {
	AmountCents   int64  `json:"amountCents"`
	Currency      string `json:"currency"`
	PaymentMethod string `json:"paymentMethod"`
	Reference     string `json:"reference"`
	Notes         string `json:"notes"`
}

type Expense struct {
	ID            string `json:"id"`
	Status        string `json:"status"`
	AmountCents   int64  `json:"amountCents"`
	Currency      string `json:"currency"`
	PaymentMethod string `json:"paymentMethod"`
	Reference     string `json:"reference"`
	Notes         string `json:"notes"`
	PaidAt        string `json:"paidAt"`
}

func (a *ExpenseAPI) ListExpenses() ([]Expense, error) {
	expenses, err := a.expenses.List(a.runtime.requestContext())
	if err != nil {
		return nil, err
	}
	result := make([]Expense, 0, len(expenses))
	for _, expense := range expenses {
		result = append(result, expenseView(expense))
	}
	return result, nil
}

func (a *ExpenseAPI) RecordExpense(input ExpenseInput) (Expense, error) {
	expense, err := a.expenses.Record(a.runtime.requestContext(), application.RecordExpenseInput(input))
	if err != nil {
		return Expense{}, err
	}
	return expenseView(expense), nil
}

func (a *ExpenseAPI) CreatePendingExpense(input ExpenseInput) (Expense, error) {
	expense, err := a.expenses.CreatePending(a.runtime.requestContext(), application.RecordExpenseInput(input))
	if err != nil {
		return Expense{}, err
	}
	return expenseView(expense), nil
}

func (a *ExpenseAPI) PostExpense(id string) (Expense, error) {
	expense, err := a.expenses.Post(a.runtime.requestContext(), id)
	if err != nil {
		return Expense{}, err
	}
	return expenseView(expense), nil
}

func (a *ExpenseAPI) VoidExpense(id string) (Expense, error) {
	expense, err := a.expenses.Void(a.runtime.requestContext(), id)
	if err != nil {
		return Expense{}, err
	}
	return expenseView(expense), nil
}

func expenseView(expense domain.Expense) Expense {
	return Expense{
		ID: string(expense.ID()), Status: string(expense.Status()), AmountCents: expense.Amount().Cents(),
		Currency: expense.Amount().Currency(), PaymentMethod: string(expense.Method()), Reference: expense.Reference(),
		Notes: expense.Notes(), PaidAt: formatOptionalTimestamp(expense.PaidAt()),
	}
}
