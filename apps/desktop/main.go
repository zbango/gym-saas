package main

import (
	"context"
	"embed"
	"log"
	"os"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/zbango/gym-saas/go/core/application"
	"github.com/zbango/gym-saas/go/core/domain"
	shared "github.com/zbango/gym-saas/go/core/platform"

	dbsqlite "github.com/zbango/gym-saas/apps/desktop/internal/sqlite"
	"github.com/zbango/gym-saas/apps/desktop/internal/updater"
)

const localGymID = "5e78f083-1eb9-4c51-a248-b3424973d3ff"

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	dbPath, err := defaultDesktopDBPath()
	if err != nil {
		log.Fatal(err)
	}

	store, err := dbsqlite.Open(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	defaultGymID, err := domain.ParseGymID(localGymID)
	if err != nil {
		log.Fatal(err)
	}
	now := time.Now().UTC()
	localGym, err := domain.NewGym(defaultGymID, "Local Gym", "UTC", now, now)
	if err != nil {
		log.Fatal(err)
	}
	if err := store.EnsureGym(context.Background(), localGym); err != nil {
		log.Fatal(err)
	}
	memberRepository := dbsqlite.NewMemberRepository(store)
	planRepository := dbsqlite.NewMembershipPlanRepository(store)
	membershipRepository := dbsqlite.NewMembershipRepository(store)
	paymentRepository := dbsqlite.NewPaymentRepository(store)
	expenseRepository := dbsqlite.NewExpenseRepository(store)
	memberService, err := application.NewMemberService(memberRepository, defaultGymID, time.Now)
	if err != nil {
		log.Fatal(err)
	}
	planService, err := application.NewMembershipPlanService(planRepository, defaultGymID, time.Now)
	if err != nil {
		log.Fatal(err)
	}
	purchaseService, err := application.NewMembershipPurchaseService(memberRepository, planRepository, dbsqlite.NewMembershipPaymentWriter(store), localGym, time.Now)
	if err != nil {
		log.Fatal(err)
	}
	paymentService, err := application.NewPaymentService(memberRepository, membershipRepository, paymentRepository, defaultGymID, time.Now)
	if err != nil {
		log.Fatal(err)
	}
	expenseService, err := application.NewExpenseService(expenseRepository, defaultGymID, time.Now)
	if err != nil {
		log.Fatal(err)
	}

	desktopVersion := shared.Version
	if override := strings.TrimSpace(os.Getenv("GYM_SAAS_DESKTOP_VERSION_OVERRIDE")); override != "" {
		desktopVersion = override
	}

	desktopRuntime := &DesktopRuntime{}
	memberAPI := NewMemberAPI(desktopRuntime, memberService)
	planAPI := NewMembershipPlanAPI(desktopRuntime, planService)
	purchaseAPI := NewMembershipPurchaseAPI(desktopRuntime, purchaseService)
	paymentAPI := NewPaymentAPI(desktopRuntime, paymentService)
	expenseAPI := NewExpenseAPI(desktopRuntime, expenseService)
	desktopAPI := NewDesktopAPI(desktopRuntime, updater.NewClient(desktopVersion), desktopVersion)

	err = wails.Run(&options.App{
		Title:  "gym-saas desktop",
		Width:  1200,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: desktopRuntime.startup,
		Bind: []any{
			memberAPI,
			planAPI,
			purchaseAPI,
			paymentAPI,
			expenseAPI,
			desktopAPI,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
