package sqlite

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/zbango/gym-saas/go/core/domain"
	"github.com/zbango/gym-saas/go/core/ports"
)

func TestPaymentRepositoryPersistsAndUpdatesReceiptAcrossRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "gym-saas.db")
	store := openTestStore(t, path)
	gymID, memberID := seedPaymentMember(t, store)
	amount := mustMoney(t, 4500)
	now := time.Date(2026, time.August, 20, 10, 0, 0, 0, time.UTC)
	payment, err := domain.CreatePendingPayment(gymID, memberID, "", domain.PaymentKindInitial, amount, domain.PaymentMethodCash, "receipt-1", "", now)
	if err != nil {
		t.Fatalf("CreatePendingPayment returned error: %v", err)
	}
	repository := NewPaymentRepository(store)
	if err := repository.Create(context.Background(), payment); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	posted, err := payment.Post(now.Add(time.Hour))
	if err != nil {
		t.Fatalf("Post returned error: %v", err)
	}
	if err := repository.Update(context.Background(), posted); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}

	store = openTestStore(t, path)
	defer store.Close()
	found, err := NewPaymentRepository(store).Get(context.Background(), gymID, payment.ID())
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	assertPaymentEqual(t, found, posted)
	payments, err := NewPaymentRepository(store).ListForMember(context.Background(), gymID, memberID)
	if err != nil {
		t.Fatalf("ListForMember returned error: %v", err)
	}
	if len(payments) != 1 {
		t.Fatalf("ListForMember length = %d, want 1", len(payments))
	}
}

func TestPaymentRepositoryReturnsTenantScopedNotFound(t *testing.T) {
	store := openTestStore(t, filepath.Join(t.TempDir(), "gym-saas.db"))
	defer store.Close()
	gymID, memberID := seedPaymentMember(t, store)
	payment, err := domain.RecordPayment(gymID, memberID, "", domain.PaymentKindOther, mustMoney(t, 500), domain.PaymentMethodOther, "", "", time.Date(2026, time.August, 20, 10, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("RecordPayment returned error: %v", err)
	}
	repository := NewPaymentRepository(store)
	if err := repository.Create(context.Background(), payment); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	otherGymID := seedGym(t, store)
	if _, err := repository.Get(context.Background(), otherGymID, payment.ID()); !errors.Is(err, ports.ErrPaymentNotFound) {
		t.Fatalf("Get other tenant error = %v, want %v", err, ports.ErrPaymentNotFound)
	}
}

func TestPaymentRepositoryPersistsMembershipReference(t *testing.T) {
	store := openTestStore(t, filepath.Join(t.TempDir(), "gym-saas.db"))
	defer store.Close()
	membership, gymID := seedVisitMembership(t, store)
	if err := NewMembershipRepository(store).Create(context.Background(), membership); err != nil {
		t.Fatalf("create membership: %v", err)
	}
	payment, err := domain.RecordPayment(gymID, membership.MemberID(), membership.ID(), domain.PaymentKindInitial, mustMoney(t, 4500), domain.PaymentMethodCash, "", "", membership.StartsAt())
	if err != nil {
		t.Fatalf("RecordPayment returned error: %v", err)
	}
	if err := NewPaymentRepository(store).Create(context.Background(), payment); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	found, err := NewPaymentRepository(store).Get(context.Background(), gymID, payment.ID())
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if found.MembershipID() != membership.ID() {
		t.Fatalf("payment membership ID = %q, want %q", found.MembershipID(), membership.ID())
	}
}

func TestExpenseRepositoryPersistsAndUpdatesExpenseAcrossRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "gym-saas.db")
	store := openTestStore(t, path)
	gymID := seedGym(t, store)
	now := time.Date(2026, time.August, 20, 10, 0, 0, 0, time.UTC)
	expense, err := domain.CreatePendingExpense(gymID, mustMoney(t, 1250), domain.PaymentMethodBankTransfer, "vendor-10", "cleaning", now)
	if err != nil {
		t.Fatalf("CreatePendingExpense returned error: %v", err)
	}
	repository := NewExpenseRepository(store)
	if err := repository.Create(context.Background(), expense); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	posted, err := expense.Post(now.Add(time.Hour))
	if err != nil {
		t.Fatalf("Post returned error: %v", err)
	}
	if err := repository.Update(context.Background(), posted); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}

	store = openTestStore(t, path)
	defer store.Close()
	found, err := NewExpenseRepository(store).Get(context.Background(), gymID, expense.ID())
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	assertExpenseEqual(t, found, posted)
	expenses, err := NewExpenseRepository(store).List(context.Background(), gymID)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(expenses) != 1 {
		t.Fatalf("List length = %d, want 1", len(expenses))
	}
}

func TestExpenseRepositoryReturnsNotFound(t *testing.T) {
	store := openTestStore(t, filepath.Join(t.TempDir(), "gym-saas.db"))
	defer store.Close()
	gymID := seedGym(t, store)
	expenseID, err := domain.NewExpenseID()
	if err != nil {
		t.Fatalf("NewExpenseID returned error: %v", err)
	}
	if _, err := NewExpenseRepository(store).Get(context.Background(), gymID, expenseID); !errors.Is(err, ports.ErrExpenseNotFound) {
		t.Fatalf("Get error = %v, want %v", err, ports.ErrExpenseNotFound)
	}
}

func seedPaymentMember(t *testing.T, store *Store) (domain.GymID, domain.MemberID) {
	t.Helper()
	gymID := seedGym(t, store)
	member := mustMember(t, gymID, "Ada", "Lovelace")
	if err := NewMemberRepository(store).Create(context.Background(), member); err != nil {
		t.Fatalf("create member: %v", err)
	}
	return gymID, member.ID()
}

func mustMoney(t *testing.T, cents int64) domain.Money {
	t.Helper()
	money, err := domain.NewMoney(cents, "USD")
	if err != nil {
		t.Fatalf("NewMoney returned error: %v", err)
	}
	return money
}

func assertPaymentEqual(t *testing.T, got, want domain.Payment) {
	t.Helper()
	if got.ID() != want.ID() || got.GymID() != want.GymID() || got.MemberID() != want.MemberID() || got.MembershipID() != want.MembershipID() || got.Status() != want.Status() || got.Kind() != want.Kind() || got.Amount().Cents() != want.Amount().Cents() || got.Amount().Currency() != want.Amount().Currency() || got.Method() != want.Method() || got.Reference() != want.Reference() || got.Notes() != want.Notes() || !got.PaidAt().Equal(want.PaidAt()) || !got.CreatedAt().Equal(want.CreatedAt()) || !got.UpdatedAt().Equal(want.UpdatedAt()) {
		t.Fatalf("payment = %#v, want %#v", got, want)
	}
}

func assertExpenseEqual(t *testing.T, got, want domain.Expense) {
	t.Helper()
	if got.ID() != want.ID() || got.GymID() != want.GymID() || got.Status() != want.Status() || got.Amount().Cents() != want.Amount().Cents() || got.Amount().Currency() != want.Amount().Currency() || got.Method() != want.Method() || got.Reference() != want.Reference() || got.Notes() != want.Notes() || !got.PaidAt().Equal(want.PaidAt()) || !got.CreatedAt().Equal(want.CreatedAt()) || !got.UpdatedAt().Equal(want.UpdatedAt()) {
		t.Fatalf("expense = %#v, want %#v", got, want)
	}
}
