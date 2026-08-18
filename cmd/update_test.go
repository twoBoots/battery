package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twoboots/battery/internal/updater"
)

func TestUpdateCmd_Help(t *testing.T) {
	buf := new(bytes.Buffer)
	UpdateCmd.SetOut(buf)
	err := UpdateCmd.Help()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Checks GitHub Releases for new versions of Battery")
	assert.Contains(t, buf.String(), "--check")
	assert.Contains(t, buf.String(), "--force")
	assert.Contains(t, buf.String(), "--target-version")
}

func TestUpdateCmd_CheckOnly(t *testing.T) {
	mockRelease := updater.Release{
		TagName: "v9.9.9",
		Name:    "Battery v9.9.9",
		Assets: []updater.Asset{
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

	client := updater.NewClientWithBaseURL(server.URL)

	// Update Available
	buf := new(bytes.Buffer)
	err := RunUpdate(buf, buf, updater.Options{
		CheckOnly:      true,
		CurrentVersion: "1.0.0",
		Repo:           "twoBoots/battery",
	}, client)
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Update available: 1.0.0 -> v9.9.9")
	assert.Contains(t, buf.String(), "Run 'battery update' to install the new version.")

	// Already Up To Date
	buf.Reset()
	err = RunUpdate(buf, buf, updater.Options{
		CheckOnly:      true,
		CurrentVersion: "9.9.9",
		Repo:           "twoBoots/battery",
	}, client)
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Battery is already up to date")
}

func TestUpdateCmd_ApplyUpdate(t *testing.T) {
	tempDir := t.TempDir()
	binaryName := "battery"
	if runtime.GOOS == "windows" {
		binaryName = "battery.exe"
	}
	execPath := filepath.Join(tempDir, binaryName)
	err := os.WriteFile(execPath, []byte("legacy-binary"), 0755)
	require.NoError(t, err)

	platformAsset, err := updater.GetCurrentPlatformBinaryName()
	require.NoError(t, err)

	newBinary := []byte("new-upgraded-binary")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/download/"+platformAsset {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(newBinary)
			return
		}
		if r.URL.Path == "/repos/twoBoots/battery/releases/latest" {
			w.Header().Set("Content-Type", "application/json")
			mockRelease := updater.Release{
				TagName: "v9.9.9",
				Name:    "Battery v9.9.9",
				Assets: []updater.Asset{
					{
						Name:               platformAsset,
						BrowserDownloadURL: "http://" + r.Host + "/download/" + platformAsset,
						Size:               int64(len(newBinary)),
					},
				},
			}
			_ = json.NewEncoder(w).Encode(mockRelease)
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	client := updater.NewClientWithBaseURL(server.URL)

	buf := new(bytes.Buffer)
	err = RunUpdate(buf, buf, updater.Options{
		CurrentVersion: "1.0.0",
		ExecutablePath: execPath,
		Force:          true,
		Repo:           "twoBoots/battery",
	}, client)

	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Successfully updated Battery to v9.9.9")

	// Verify binary was actually swapped
	content, err := os.ReadFile(execPath)
	assert.NoError(t, err)
	assert.Equal(t, newBinary, content)

	// Up to date without force
	buf.Reset()
	err = RunUpdate(buf, buf, updater.Options{
		CurrentVersion: "9.9.9",
		ExecutablePath: execPath,
		Force:          false,
		Repo:           "twoBoots/battery",
	}, client)
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Battery is already up to date")
}

func TestUpdateCmd_Errors(t *testing.T) {
	badClient := updater.NewClientWithBaseURL("http://invalid.localhost:9999")
	buf := new(bytes.Buffer)
	err := RunUpdate(buf, buf, updater.Options{
		CurrentVersion: "1.0.0",
		Repo:           "twoBoots/battery",
	}, badClient)
	assert.Error(t, err)
}

func TestUpdateCmd_Aliases(t *testing.T) {
	assert.Contains(t, UpdateCmd.Aliases, "self-update")
}

func TestUpdateCmd_RootCmdExecution(t *testing.T) {
	mockRelease := updater.Release{
		TagName: "v1.2.0",
		Name:    "Battery v1.2.0",
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockRelease)
	}))
	defer server.Close()

	SetUpdaterClient(updater.NewClientWithBaseURL(server.URL))
	defer SetUpdaterClient(nil)

	buf := new(bytes.Buffer)
	RootCmd.SetOut(buf)
	UpdateCmd.SetOut(buf)
	RootCmd.SetArgs([]string{"update", "--check"})
	err := RootCmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Checking for updates")
}
