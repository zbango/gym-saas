package main

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	dbsqlite "github.com/zbango/gym-saas/apps/desktop/internal/sqlite"
	"github.com/zbango/gym-saas/go/core/application"
	"github.com/zbango/gym-saas/go/core/domain"
)

func TestFinancialAPIsUseFeatureServices(t *testing.T) {
	store, err := dbsqlite.Open(filepath.Join(t.TempDir(), "gym-saas.db"))
	if err != nil {
		t.Fatalf("Open returned error: %v", err)
	}
	defer store.Close()
	now := time.Date(2026, time.August, 20, 10, 0, 0, 0, time.UTC)
	gymID, err := domain.NewGymID()
	if err != nil {
		t.Fatalf("NewGymID returned error: %v", err)
	}
	gym, err := domain.NewGym(gymID, "Test gym", "UTC", now, now)
	if err != nil {
		t.Fatalf("NewGym returned error: %v", err)
	}
	if err := store.EnsureGym(context.Background(), gym); err != nil {
		t.Fatalf("EnsureGym returned error: %v", err)
	}
	members := dbsqlite.NewMemberRepository(store)
	plans := dbsqlite.NewMembershipPlanRepository(store)
	memberships := dbsqlite.NewMembershipRepository(store)
	payments := dbsqlite.NewPaymentRepository(store)
	expenses := dbsqlite.NewExpenseRepository(store)
	memberService, err := application.NewMemberService(members, gymID, func() time.Time { return now })
	if err != nil {
		t.Fatalf("NewMemberService returned error: %v", err)
	}
	planService, err := application.NewMembershipPlanService(plans, gymID, func() time.Time { return now })
	if err != nil {
		t.Fatalf("NewMembershipPlanService returned error: %v", err)
	}
	purchaseService, err := application.NewMembershipPurchaseService(members, plans, dbsqlite.NewMembershipPaymentWriter(store), gym, func() time.Time { return now })
	if err != nil {
		t.Fatalf("NewMembershipPurchaseService returned error: %v", err)
	}
	paymentService, err := application.NewPaymentService(members, memberships, payments, gymID, func() time.Time { return now })
	if err != nil {
		t.Fatalf("NewPaymentService returned error: %v", err)
	}
	expenseService, err := application.NewExpenseService(expenses, gymID, func() time.Time { return now })
	if err != nil {
		t.Fatalf("NewExpenseService returned error: %v", err)
	}
	runtime := &DesktopRuntime{}
	memberAPI := NewMemberAPI(runtime, memberService)
	planAPI := NewMembershipPlanAPI(runtime, planService)
	purchaseAPI := NewMembershipPurchaseAPI(runtime, purchaseService)
	paymentAPI := NewPaymentAPI(runtime, paymentService)
	expenseAPI := NewExpenseAPI(runtime, expenseService)

	member, err := memberAPI.CreateMember(MemberInput{FirstName: "Ada", LastName: "Lovelace", Phone: "555-0100", Status: "active"})
	if err != nil {
		t.Fatalf("CreateMember returned error: %v", err)
	}
	plan, err := planAPI.CreateMembershipPlan(MembershipPlanInput{Name: "Monthly", ValidityKind: "time", DurationValue: 1, DurationUnit: "months", PriceCents: 4500, Currency: "USD", Status: "active"})
	if err != nil {
		t.Fatalf("CreateMembershipPlan returned error: %v", err)
	}
	purchase, err := purchaseAPI.PurchaseMembership(PurchaseMembershipInput{MemberID: member.ID, MembershipPlanID: plan.ID, AmountCents: 4500, Currency: "USD", PaymentMethod: "cash", Reference: "receipt-1"})
	if err != nil {
		t.Fatalf("PurchaseMembership returned error: %v", err)
	}
	if purchase.Membership.ID == "" || purchase.Payment.MembershipID != purchase.Membership.ID || purchase.Payment.Status != "posted" {
		t.Fatalf("PurchaseMembership = %#v", purchase)
	}
	memberPayments, err := paymentAPI.ListPaymentsForMember(member.ID)
	if err != nil {
		t.Fatalf("ListPaymentsForMember returned error: %v", err)
	}
	if len(memberPayments) != 1 || memberPayments[0].ID != purchase.Payment.ID {
		t.Fatalf("ListPaymentsForMember = %#v", memberPayments)
	}
	expense, err := expenseAPI.RecordExpense(ExpenseInput{AmountCents: 1250, Currency: "USD", PaymentMethod: "bank_transfer", Reference: "vendor-1"})
	if err != nil {
		t.Fatalf("RecordExpense returned error: %v", err)
	}
	listedExpenses, err := expenseAPI.ListExpenses()
	if err != nil {
		t.Fatalf("ListExpenses returned error: %v", err)
	}
	if len(listedExpenses) != 1 || listedExpenses[0].ID != expense.ID {
		t.Fatalf("ListExpenses = %#v", listedExpenses)
	}
}
