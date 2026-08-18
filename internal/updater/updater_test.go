package updater

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSelfUpdate_CheckOnly(t *testing.T) {
	mockRelease := Release{
		TagName: "v1.4.0",
		Name:    "Battery v1.4.0",
		Assets: []Asset{
			{
				Name:               "battery-darwin-aarch64",
				BrowserDownloadURL: "https://example.com/download/darwin-arm64",
			},
			{
				Name:               "battery-darwin-x86_64",
				BrowserDownloadURL: "https://example.com/download/darwin-amd64",
			},
			{
				Name:               "battery-linux-x86_64",
				BrowserDownloadURL: "https://example.com/download/linux-amd64",
			},
			{
				Name:               "battery-linux-aarch64",
				BrowserDownloadURL: "https://example.com/download/linux-arm64",
			},
			{
				Name:               "battery-windows-x86_64.exe",
				BrowserDownloadURL: "https://example.com/download/windows-amd64",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockRelease)
	}))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)

	// Subtest 1: Update available
	res, err := SelfUpdateWithClient(client, Options{
		CurrentVersion: "1.2.0",
		CheckOnly:      true,
		Repo:           "twoBoots/battery",
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.False(t, res.Updated)
	assert.True(t, res.UpdateAvailable)
	assert.Equal(t, "1.2.0", res.CurrentVersion)
	assert.Equal(t, "v1.4.0", res.LatestVersion)

	// Subtest 2: Up to date
	res, err = SelfUpdateWithClient(client, Options{
		CurrentVersion: "1.4.0",
		CheckOnly:      true,
		Repo:           "twoBoots/battery",
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.False(t, res.Updated)
	assert.False(t, res.UpdateAvailable)
}

func TestSelfUpdate_ApplyUpdate(t *testing.T) {
	tempDir := t.TempDir()
	binaryName := "battery"
	if runtime.GOOS == "windows" {
		binaryName = "battery.exe"
	}
	execPath := filepath.Join(tempDir, binaryName)
	err := os.WriteFile(execPath, []byte("old-binary-content"), 0755)
	require.NoError(t, err)

	platformAsset, err := GetCurrentPlatformBinaryName()
	require.NoError(t, err)

	newContent := []byte("new-binary-content-v2.0.0")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/download/"+platformAsset {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(newContent)
			return
		}
		if r.URL.Path == "/repos/twoBoots/battery/releases/latest" {
			w.Header().Set("Content-Type", "application/json")
			mockRelease := Release{
				TagName: "v2.0.0",
				Name:    "Battery v2.0.0",
				Assets: []Asset{
					{
						Name:               platformAsset,
						BrowserDownloadURL: "http://" + r.Host + "/download/" + platformAsset,
						Size:               int64(len(newContent)),
					},
				},
			}
			_ = json.NewEncoder(w).Encode(mockRelease)
			return
		}
		if r.URL.Path == "/repos/twoBoots/battery/releases/tags/v2.0.0" {
			w.Header().Set("Content-Type", "application/json")
			mockRelease := Release{
				TagName: "v2.0.0",
				Name:    "Battery v2.0.0",
				Assets: []Asset{
					{
						Name:               platformAsset,
						BrowserDownloadURL: "http://" + r.Host + "/download/" + platformAsset,
						Size:               int64(len(newContent)),
					},
				},
			}
			_ = json.NewEncoder(w).Encode(mockRelease)
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)

	// Apply Update
	res, err := SelfUpdateWithClient(client, Options{
		CurrentVersion: "1.0.0",
		ExecutablePath: execPath,
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.True(t, res.Updated)
	assert.Equal(t, "v2.0.0", res.LatestVersion)

	// Verify content was replaced
	replacedContent, err := os.ReadFile(execPath)
	assert.NoError(t, err)
	assert.Equal(t, newContent, replacedContent)

	// Try running again without --force (should not update)
	res, err = SelfUpdateWithClient(client, Options{
		CurrentVersion: "2.0.0",
		ExecutablePath: execPath,
	})
	assert.NoError(t, err)
	assert.False(t, res.Updated)

	// Run with --force and specific TargetVersion (should update)
	res, err = SelfUpdateWithClient(client, Options{
		CurrentVersion: "2.0.0",
		TargetVersion:  "v2.0.0",
		ExecutablePath: execPath,
		Force:          true,
	})
	assert.NoError(t, err)
	assert.True(t, res.Updated)
}

func TestSelfUpdate_Errors(t *testing.T) {
	tempDir := t.TempDir()
	execPath := filepath.Join(tempDir, "battery")
	err := os.WriteFile(execPath, []byte("old-binary-content"), 0755)
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/download/error" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if r.URL.Path == "/repos/twoBoots/battery/releases/latest" {
			w.Header().Set("Content-Type", "application/json")
			mockRelease := Release{
				TagName: "v3.0.0",
				Name:    "Battery v3.0.0",
				Assets: []Asset{
					{
						Name:               "battery-darwin-aarch64",
						BrowserDownloadURL: "http://" + r.Host + "/download/error",
					},
					{
						Name:               "battery-darwin-x86_64",
						BrowserDownloadURL: "http://" + r.Host + "/download/error",
					},
					{
						Name:               "battery-linux-x86_64",
						BrowserDownloadURL: "http://" + r.Host + "/download/error",
					},
					{
						Name:               "battery-linux-aarch64",
						BrowserDownloadURL: "http://" + r.Host + "/download/error",
					},
					{
						Name:               "battery-windows-x86_64.exe",
						BrowserDownloadURL: "http://" + r.Host + "/download/error",
					},
				},
			}
			_ = json.NewEncoder(w).Encode(mockRelease)
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)

	// Download failure error test
	_, err = SelfUpdateWithClient(client, Options{
		CurrentVersion: "1.0.0",
		ExecutablePath: execPath,
	})
	assert.Error(t, err)

	// Release fetch error test
	badClient := NewClientWithBaseURL("http://invalid.localhost:9999")
	_, err = SelfUpdateWithClient(badClient, Options{
		CurrentVersion: "1.0.0",
		ExecutablePath: execPath,
	})
	assert.Error(t, err)
}
