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
	memberService, err := application.NewMemberService(dbsqlite.NewMemberRepository(store), defaultGymID, time.Now)
	if err != nil {
		log.Fatal(err)
	}

	desktopVersion := shared.Version
	if override := strings.TrimSpace(os.Getenv("GYM_SAAS_DESKTOP_VERSION_OVERRIDE")); override != "" {
		desktopVersion = override
	}

	app := NewApp(updater.NewClient(desktopVersion), desktopVersion, memberService)

	err = wails.Run(&options.App{
		Title:  "gym-saas desktop",
		Width:  1200,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		Bind: []any{
			app,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
