package updater

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// DefaultGitHubAPIBaseURL is the base URL for GitHub API requests.
const DefaultGitHubAPIBaseURL = "https://api.github.com"

// Release represents a GitHub Release object.
type Release struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	Body        string    `json:"body"`
	Prerelease  bool      `json:"prerelease"`
	Draft       bool      `json:"draft"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []Asset   `json:"assets"`
}

// Asset represents a release asset binary.
type Asset struct {
	Name               string `json:"name"`
	Size               int64  `json:"size"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// FindAssetForPlatform finds the matching release asset for the given OS and Architecture.
func (r *Release) FindAssetForPlatform(goos, goarch string) (*Asset, error) {
	expectedName, err := GetPlatformBinaryName(goos, goarch)
	if err != nil {
		return nil, err
	}

	for _, asset := range r.Assets {
		if asset.Name == expectedName {
			return &asset, nil
		}
	}

	return nil, fmt.Errorf("no release asset found for platform %s/%s (expected %s)", goos, goarch, expectedName)
}

// Client handles communication with GitHub Releases.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new GitHub release client.
func NewClient() *Client {
	return NewClientWithBaseURL(DefaultGitHubAPIBaseURL)
}

// NewClientWithBaseURL creates a new GitHub release client with a custom base URL.
func NewClientWithBaseURL(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FetchLatestRelease retrieves the latest release for the given repository (e.g. "twoBoots/battery").
func (c *Client) FetchLatestRelease(repo string) (*Release, error) {
	url := fmt.Sprintf("%s/repos/%s/releases/latest", c.baseURL, repo)
	return c.fetchRelease(url)
}

// FetchReleaseByTag retrieves a specific release by tag.
func (c *Client) FetchReleaseByTag(repo, tag string) (*Release, error) {
	url := fmt.Sprintf("%s/repos/%s/releases/tags/%s", c.baseURL, repo, tag)
	return c.fetchRelease(url)
}

func (c *Client) fetchRelease(url string) (*Release, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "Battery-CLI-Updater")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api returned status %d for %s", resp.StatusCode, url)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to decode release response: %w", err)
	}

	return &release, nil
}
