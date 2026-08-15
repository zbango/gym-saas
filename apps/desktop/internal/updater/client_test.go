package updater

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestClientCheckReturnsManifestForNewerVersion(t *testing.T) {
	client := NewClient("0.1.0")
	client.httpClient = &http.Client{
		Transport: roundTripper(func(*http.Request) (*http.Response, error) {
			return jsonResponse(`{
      "version":"0.1.1",
      "minimumSupportedVersion":"0.1.0",
      "notes":"Update available",
      "publishedAt":"2026-08-14T00:00:00Z",
      "assets":[
        {"os":"` + runtime.GOOS + `","arch":"` + runtime.GOARCH + `","url":"https://example.com/app","checksum":"abc"},
        {"os":"windows","arch":"amd64","url":"https://example.com/windows","checksum":"def"}
      ]
    }`)
		}),
	}

	manifest, err := client.Check("https://updates.example.com/manifest")
	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}
	if manifest == nil {
		t.Fatal("Check returned nil manifest, want update")
	}
	if manifest.Version != "0.1.1" {
		t.Fatalf("manifest.Version = %q, want %q", manifest.Version, "0.1.1")
	}
	if len(manifest.Assets) != 1 {
		t.Fatalf("manifest.Assets length = %d, want 1", len(manifest.Assets))
	}
}

func TestClientCheckReturnsNilForSameVersion(t *testing.T) {
	client := NewClient("0.1.0")
	client.httpClient = &http.Client{
		Transport: roundTripper(func(*http.Request) (*http.Response, error) {
			return jsonResponse(`{
      "version":"0.1.0",
      "minimumSupportedVersion":"0.1.0",
      "notes":"No change",
      "publishedAt":"2026-08-14T00:00:00Z",
      "assets":[
        {"os":"` + runtime.GOOS + `","arch":"` + runtime.GOARCH + `","url":"https://example.com/app","checksum":"abc"}
      ]
    }`)
		}),
	}

	manifest, err := client.Check("https://updates.example.com/manifest")
	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}
	if manifest != nil {
		t.Fatalf("Check returned manifest %+v, want nil", manifest)
	}
}

func TestClientDownloadWritesFile(t *testing.T) {
	tempDir := t.TempDir()
	client := NewClient("0.1.0")
	client.downloadDir = func() (string, error) {
		return tempDir, nil
	}
	client.httpClient = &http.Client{
		Transport: roundTripper(func(*http.Request) (*http.Response, error) {
			return jsonResponse("hello update package")
		}),
	}

	result, err := client.Download(Asset{
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		URL:      "https://downloads.example.com/gym-saas-desktop-darwin-arm64.zip",
		Checksum: "replace-me",
	})
	if err != nil {
		t.Fatalf("Download returned error: %v", err)
	}

	if result.FileName != "gym-saas-desktop-darwin-arm64.zip" {
		t.Fatalf("FileName = %q", result.FileName)
	}

	content, err := os.ReadFile(result.Path)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}

	if got := string(content); got != "hello update package" {
		t.Fatalf("downloaded content = %q", got)
	}

	if filepath.Dir(result.Path) != filepath.Join(tempDir, "updates") {
		t.Fatalf("unexpected download directory: %q", filepath.Dir(result.Path))
	}
}

type roundTripper func(*http.Request) (*http.Response, error)

func (fn roundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	return fn(request)
}

func jsonResponse(body string) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}, nil
}
