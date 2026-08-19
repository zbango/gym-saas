package main

import (
	"embed"
	"log"
	"os"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	shared "github.com/zbango/gym-saas/go/core/platform"

	dbsqlite "github.com/zbango/gym-saas/apps/desktop/internal/sqlite"
	"github.com/zbango/gym-saas/apps/desktop/internal/updater"
)

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

	desktopVersion := shared.Version
	if override := strings.TrimSpace(os.Getenv("GYM_SAAS_DESKTOP_VERSION_OVERRIDE")); override != "" {
		desktopVersion = override
	}

	app := NewApp(updater.NewClient(desktopVersion), desktopVersion)

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
