package main

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/zbango/gym-saas/go/core/platform"
)

type updateAsset struct {
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	URL      string `json:"url"`
	Checksum string `json:"checksum"`
}

type updateManifest struct {
	Version                 string        `json:"version"`
	MinimumSupportedVersion string        `json:"minimumSupportedVersion"`
	Notes                   string        `json:"notes"`
	PublishedAt             string        `json:"publishedAt"`
	Assets                  []updateAsset `json:"assets"`
}

func main() {
	lambda.Start(handle)
}

func handle(_ context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	switch request.RawPath {
	case "/health":
		return jsonResponse(map[string]string{"status": "ok"})
	case "/version":
		return jsonResponse(map[string]string{"version": platform.Version})
	default:
		return jsonResponse(releaseManifest())
	}
}

func releaseManifest() updateManifest {
	return updateManifest{
		Version:                 envOr("RELEASE_VERSION", "0.1.1"),
		MinimumSupportedVersion: envOr("RELEASE_MINIMUM_SUPPORTED_VERSION", "0.1.0"),
		Notes:                   envOr("RELEASE_NOTES", "Initial updater manifest for local-first gym-saas."),
		PublishedAt:             time.Now().UTC().Format(time.RFC3339),
		Assets:                  desktopReleaseAssets(),
	}
}

func desktopReleaseAssets() []updateAsset {
	assets := []updateAsset{
		{
			OS:       "darwin",
			Arch:     "arm64",
			URL:      envOr("RELEASE_DARWIN_ARM64_URL", ""),
			Checksum: envOr("RELEASE_DARWIN_ARM64_CHECKSUM", ""),
		},
		{
			OS:       "windows",
			Arch:     "amd64",
			URL:      envOr("RELEASE_WINDOWS_AMD64_URL", ""),
			Checksum: envOr("RELEASE_WINDOWS_AMD64_CHECKSUM", ""),
		},
		{
			OS:       "linux",
			Arch:     "amd64",
			URL:      envOr("RELEASE_LINUX_AMD64_URL", ""),
			Checksum: envOr("RELEASE_LINUX_AMD64_CHECKSUM", ""),
		},
	}

	filtered := make([]updateAsset, 0, len(assets))
	for _, asset := range assets {
		if strings.TrimSpace(asset.URL) == "" {
			continue
		}
		filtered = append(filtered, asset)
	}

	return filtered
}

func envOr(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	return value
}

func jsonResponse(payload any) (events.LambdaFunctionURLResponse, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return events.LambdaFunctionURLResponse{}, err
	}

	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"content-type": "application/json",
		},
		Body: string(body),
	}, nil
}
