package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/zbango/gym-saas/apps/desktop/internal/updater"
)

type App struct {
	ctx     context.Context
	updates *updater.Client
	version string
}

func NewApp(updates *updater.Client, version string) *App {
	return &App{updates: updates, version: version}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s from gym-saas %s", name, a.version)
}

func (a *App) GetDesktopVersion() string {
	return a.version
}

func (a *App) CheckForUpdates(manifestURL string) (*updater.Manifest, error) {
	manifest, err := a.updates.Check(manifestURL)
	if err != nil {
		return nil, err
	}

	return manifest, nil
}

func (a *App) DownloadUpdatePackage(url, checksum string) (*updater.DownloadedPackage, error) {
	result, err := a.updates.Download(updater.Asset{
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		URL:      url,
		Checksum: checksum,
	})
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (a *App) InstallDownloadedPackage(path string) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("path is required")
	}

	switch runtime.GOOS {
	case "darwin":
		return a.installDownloadedPackageDarwin(path)
	default:
		return a.OpenPath(path)
	}
}

func (a *App) OpenExternalURL(url string) error {
	if strings.TrimSpace(url) == "" {
		return fmt.Errorf("url is required")
	}

	wailsruntime.BrowserOpenURL(a.ctx, url)
	return nil
}

func (a *App) OpenPath(path string) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("path is required")
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", path)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("open path: %w", err)
	}

	return nil
}

func defaultDesktopDBPath() (string, error) {
	configDir, err := updater.AppConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "gym-saas.db"), nil
}

func (a *App) installDownloadedPackageDarwin(packagePath string) error {
	configDir, err := updater.AppConfigDir()
	if err != nil {
		return err
	}

	stageDir := filepath.Join(configDir, "updates", fmt.Sprintf("stage-%d", time.Now().UnixNano()))
	if err := os.MkdirAll(stageDir, 0o755); err != nil {
		return fmt.Errorf("create stage dir: %w", err)
	}

	if strings.HasSuffix(strings.ToLower(packagePath), ".zip") {
		unzip := exec.Command("ditto", "-x", "-k", packagePath, stageDir)
		if output, err := unzip.CombinedOutput(); err != nil {
			return fmt.Errorf("extract package: %w: %s", err, strings.TrimSpace(string(output)))
		}
	} else {
		return a.OpenPath(packagePath)
	}

	nextAppPath, err := firstAppBundle(stageDir)
	if err != nil {
		return err
	}

	currentAppPath, replaceable, err := currentAppBundlePath()
	if err != nil {
		return err
	}

	if !replaceable {
		if err := a.OpenPath(nextAppPath); err != nil {
			return err
		}
		wailsruntime.Quit(a.ctx)
		return nil
	}

	scriptPath := filepath.Join(stageDir, "install-update.sh")
	backupPath := currentAppPath + ".previous"
	script := fmt.Sprintf(`#!/bin/bash
set -e
sleep 2
rm -rf %q
if [ -d %q ]; then mv %q %q; fi
rm -rf %q
mv %q %q
open %q
`, backupPath, currentAppPath, currentAppPath, backupPath, currentAppPath, nextAppPath, currentAppPath, currentAppPath)

	if err := os.WriteFile(scriptPath, []byte(script), 0o755); err != nil {
		return fmt.Errorf("write install script: %w", err)
	}

	cmd := exec.Command("/bin/bash", scriptPath)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start install script: %w", err)
	}

	wailsruntime.Quit(a.ctx)
	return nil
}

func currentAppBundlePath() (string, bool, error) {
	executable, err := os.Executable()
	if err != nil {
		return "", false, fmt.Errorf("resolve executable: %w", err)
	}

	macosDir := filepath.Dir(executable)
	contentsDir := filepath.Dir(macosDir)
	appDir := filepath.Dir(contentsDir)

	if filepath.Base(macosDir) != "MacOS" || filepath.Base(contentsDir) != "Contents" || !strings.HasSuffix(appDir, ".app") {
		return "", false, nil
	}

	return appDir, true, nil
}

func firstAppBundle(root string) (string, error) {
	var found string

	err := filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() && strings.HasSuffix(info.Name(), ".app") {
			found = path
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("scan extracted package: %w", err)
	}
	if found == "" {
		return "", fmt.Errorf("no .app bundle found in extracted package")
	}

	return found, nil
}
