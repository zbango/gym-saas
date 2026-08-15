package updater

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Asset struct {
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	URL      string `json:"url"`
	Checksum string `json:"checksum"`
}

type Manifest struct {
	Version                 string  `json:"version"`
	MinimumSupportedVersion string  `json:"minimumSupportedVersion"`
	Notes                   string  `json:"notes"`
	PublishedAt             string  `json:"publishedAt"`
	Assets                  []Asset `json:"assets"`
}

type DownloadedPackage struct {
	FileName string `json:"fileName"`
	Path     string `json:"path"`
}

type Client struct {
	currentVersion string
	httpClient     *http.Client
	downloadDir    func() (string, error)
}

func NewClient(currentVersion string) *Client {
	return &Client{
		currentVersion: currentVersion,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		downloadDir: AppConfigDir,
	}
}

func (c *Client) Check(manifestURL string) (*Manifest, error) {
	response, err := c.httpClient.Get(manifestURL)
	if err != nil {
		return nil, fmt.Errorf("fetch manifest: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("manifest status: %s", response.Status)
	}

	var manifest Manifest
	if err := json.NewDecoder(response.Body).Decode(&manifest); err != nil {
		return nil, fmt.Errorf("decode manifest: %w", err)
	}

	filtered := make([]Asset, 0, len(manifest.Assets))
	for _, asset := range manifest.Assets {
		if strings.EqualFold(asset.OS, runtime.GOOS) && strings.EqualFold(asset.Arch, runtime.GOARCH) {
			filtered = append(filtered, asset)
		}
	}
	manifest.Assets = filtered

	if manifest.Version == c.currentVersion {
		return nil, nil
	}

	return &manifest, nil
}

func (c *Client) Download(asset Asset) (DownloadedPackage, error) {
	if strings.TrimSpace(asset.URL) == "" {
		return DownloadedPackage{}, fmt.Errorf("asset url is required")
	}

	response, err := c.httpClient.Get(asset.URL)
	if err != nil {
		return DownloadedPackage{}, fmt.Errorf("download asset: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return DownloadedPackage{}, fmt.Errorf("asset status: %s", response.Status)
	}

	configDir, err := c.downloadDir()
	if err != nil {
		return DownloadedPackage{}, err
	}

	updatesDir := filepath.Join(configDir, "updates")
	if err := os.MkdirAll(updatesDir, 0o755); err != nil {
		return DownloadedPackage{}, fmt.Errorf("create updates dir: %w", err)
	}

	fileName := fileNameFromURL(asset.URL)
	targetPath := filepath.Join(updatesDir, fileName)
	tempPath := targetPath + ".download"

	file, err := os.Create(tempPath)
	if err != nil {
		return DownloadedPackage{}, fmt.Errorf("create temp file: %w", err)
	}

	hasher := sha256.New()
	writer := io.MultiWriter(file, hasher)
	_, copyErr := io.Copy(writer, response.Body)
	closeErr := file.Close()
	if copyErr != nil {
		_ = os.Remove(tempPath)
		return DownloadedPackage{}, fmt.Errorf("copy asset: %w", copyErr)
	}
	if closeErr != nil {
		_ = os.Remove(tempPath)
		return DownloadedPackage{}, fmt.Errorf("close temp file: %w", closeErr)
	}

	if err := verifyChecksum(asset.Checksum, fmt.Sprintf("%x", hasher.Sum(nil))); err != nil {
		_ = os.Remove(tempPath)
		return DownloadedPackage{}, err
	}

	if err := os.Rename(tempPath, targetPath); err != nil {
		_ = os.Remove(tempPath)
		return DownloadedPackage{}, fmt.Errorf("rename downloaded file: %w", err)
	}

	return DownloadedPackage{
		FileName: fileName,
		Path:     targetPath,
	}, nil
}

func AppConfigDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve user config dir: %w", err)
	}

	path := filepath.Join(dir, "gym-saas")
	if err := os.MkdirAll(path, 0o755); err != nil {
		return "", fmt.Errorf("create config dir: %w", err)
	}

	return path, nil
}

func fileNameFromURL(raw string) string {
	parsed, err := url.Parse(raw)
	if err != nil {
		return "gym-saas-update"
	}

	name := filepath.Base(parsed.Path)
	if name == "." || name == "/" || name == "" {
		return "gym-saas-update"
	}

	return name
}

func verifyChecksum(expected, actual string) error {
	normalized := strings.TrimSpace(expected)
	if normalized == "" || strings.HasPrefix(normalized, "replace-me") {
		return nil
	}

	normalized = strings.TrimPrefix(normalized, "sha256:")
	if !strings.EqualFold(normalized, actual) {
		return fmt.Errorf("checksum mismatch")
	}

	return nil
}
