package updater

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"1.0.0", "1.0.0", 0},
		{"v1.0.0", "1.0.0", 0},
		{"1.0.0", "v1.0.1", -1},
		{"1.2.0", "1.1.9", 1},
		{"2.0.0", "1.99.99", 1},
		{"1.0.0-rc1", "1.0.0", -1},
		{"1.0.0", "1.0.0-rc1", 1},
		{"1.0.0-alpha", "1.0.0-beta", -1},
		{"1.0.0-beta", "1.0.0-alpha", 1},
		{"1.0.0-alpha", "1.0.0-alpha", 0},
		{"dev", "1.0.0", -1},
		{"1.0.0", "dev", 1},
	}

	for _, tt := range tests {
		t.Run(tt.v1+" vs "+tt.v2, func(t *testing.T) {
			assert.Equal(t, tt.expected, CompareVersions(tt.v1, tt.v2))
		})
	}
}

func TestFetchLatestRelease(t *testing.T) {
	mockRelease := Release{
		TagName: "v1.3.0",
		Name:    "Battery v1.3.0",
		Assets: []Asset{
			{
				Name:               "battery-darwin-aarch64",
				BrowserDownloadURL: "https://example.com/download/battery-darwin-aarch64",
				Size:               1024,
			},
			{
				Name:               "battery-linux-x86_64",
				BrowserDownloadURL: "https://example.com/download/battery-linux-x86_64",
				Size:               2048,
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/twoBoots/battery/releases/latest" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(mockRelease)
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	release, err := client.FetchLatestRelease("twoBoots/battery")
	assert.NoError(t, err)
	assert.NotNil(t, release)
	assert.Equal(t, "v1.3.0", release.TagName)
	assert.Len(t, release.Assets, 2)

	// Test Finding Asset
	asset, err := release.FindAssetForPlatform("darwin", "arm64")
	assert.NoError(t, err)
	assert.Equal(t, "battery-darwin-aarch64", asset.Name)

	// Test Missing Asset
	_, err = release.FindAssetForPlatform("windows", "amd64")
	assert.Error(t, err)

	// Test Invalid platform
	_, err = release.FindAssetForPlatform("solaris", "sparc")
	assert.Error(t, err)

	// Test 404 / error response
	_, err = client.FetchLatestRelease("twoBoots/nonexistent")
	assert.Error(t, err)
}

func TestFetchReleaseByTag(t *testing.T) {
	mockRelease := Release{
		TagName: "v1.2.5",
		Name:    "Battery v1.2.5",
		Assets: []Asset{
			{
				Name:               "battery-linux-x86_64",
				BrowserDownloadURL: "https://example.com/download/battery-linux-x86_64",
				Size:               2048,
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/repos/twoBoots/battery/releases/tags/v1.2.5", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockRelease)
	}))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	release, err := client.FetchReleaseByTag("twoBoots/battery", "v1.2.5")
	assert.NoError(t, err)
	assert.NotNil(t, release)
	assert.Equal(t, "v1.2.5", release.TagName)
}

func TestNewClientDefault(t *testing.T) {
	c := NewClient()
	assert.NotNil(t, c)
	assert.Equal(t, DefaultGitHubAPIBaseURL, c.baseURL)
}
