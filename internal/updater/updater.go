package updater

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

const DefaultRepo = "twoBoots/battery"

// Options specifies configuration for a self-update operation.
type Options struct {
	Repo           string
	TargetVersion  string
	CurrentVersion string
	ExecutablePath string
	Force          bool
	CheckOnly      bool
}

// Result holds the outcome of a self-update operation.
type Result struct {
	CurrentVersion  string
	LatestVersion   string
	UpdateAvailable bool
	Updated         bool
	ExecutablePath  string
	Message         string
}

// SelfUpdate performs an update check or applies an update using the default client.
func SelfUpdate(opts Options) (*Result, error) {
	return SelfUpdateWithClient(NewClient(), opts)
}

// SelfUpdateWithClient performs self-update with a specific release client.
func SelfUpdateWithClient(client *Client, opts Options) (*Result, error) {
	if opts.Repo == "" {
		opts.Repo = DefaultRepo
	}

	execPath := opts.ExecutablePath
	if execPath == "" {
		p, err := os.Executable()
		if err != nil {
			return nil, fmt.Errorf("failed to locate executable: %w", err)
		}
		execPath, err = filepath.EvalSymlinks(p)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve symlink for executable: %w", err)
		}
	}

	// 1. Fetch release info
	var release *Release
	var err error
	if opts.TargetVersion != "" {
		release, err = client.FetchReleaseByTag(opts.Repo, opts.TargetVersion)
	} else {
		release, err = client.FetchLatestRelease(opts.Repo)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release information: %w", err)
	}

	result := &Result{
		CurrentVersion: opts.CurrentVersion,
		LatestVersion:  release.TagName,
		ExecutablePath: execPath,
	}

	cmp := CompareVersions(opts.CurrentVersion, release.TagName)
	if cmp < 0 {
		result.UpdateAvailable = true
	}

	if opts.CheckOnly {
		if result.UpdateAvailable {
			result.Message = fmt.Sprintf("Update available: %s -> %s", opts.CurrentVersion, release.TagName)
		} else {
			result.Message = fmt.Sprintf("Battery is already up to date (%s)", opts.CurrentVersion)
		}
		return result, nil
	}

	// If no update is available and --force is not set, stop here
	if !result.UpdateAvailable && !opts.Force {
		result.Message = fmt.Sprintf("Battery is already up to date (%s). Use --force to reinstall.", opts.CurrentVersion)
		return result, nil
	}

	// 2. Find asset for current platform
	asset, err := release.FindAssetForPlatform(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return nil, fmt.Errorf("incompatible release: %w", err)
	}

	// 3. Download binary to a temporary file in the target binary directory
	execDir := filepath.Dir(execPath)
	tmpFile, err := os.CreateTemp(execDir, "battery-update-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file in %s: %w (check directory write permissions)", execDir, err)
	}
	tmpPath := tmpFile.Name()
	defer func() {
		_ = os.Remove(tmpPath) // Clean up temp file if still present
	}()

	resp, err := client.httpClient.Get(asset.BrowserDownloadURL)
	if err != nil {
		_ = tmpFile.Close()
		return nil, fmt.Errorf("failed to download release binary: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		_ = tmpFile.Close()
		return nil, fmt.Errorf("binary download failed with status %d from %s", resp.StatusCode, asset.BrowserDownloadURL)
	}

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		_ = tmpFile.Close()
		return nil, fmt.Errorf("failed to save downloaded binary: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return nil, fmt.Errorf("failed to close temporary binary file: %w", err)
	}

	// Set executable permissions
	if err := os.Chmod(tmpPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to set executable permissions: %w", err)
	}

	// 4. Atomic replace
	if err := replaceExecutable(tmpPath, execPath); err != nil {
		return nil, fmt.Errorf("failed to replace executable %s: %w", execPath, err)
	}

	result.Updated = true
	result.Message = fmt.Sprintf("Successfully updated Battery to %s (%s)", release.TagName, execPath)
	return result, nil
}

func replaceExecutable(from, to string) error {
	// On Unix systems, os.Rename replaces the target file atomically even if currently open
	if runtime.GOOS != "windows" {
		return os.Rename(from, to)
	}

	// On Windows, a running binary cannot be overwritten directly.
	// We rename the old binary first, then move the new one in place.
	oldPath := to + ".old"
	_ = os.Remove(oldPath)
	if err := os.Rename(to, oldPath); err != nil {
		return err
	}
	if err := os.Rename(from, to); err != nil {
		_ = os.Rename(oldPath, to) // Rollback
		return err
	}
	_ = os.Remove(oldPath)
	return nil
}
