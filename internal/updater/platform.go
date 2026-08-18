package updater

import (
	"fmt"
	"runtime"
)

// GetCurrentPlatformBinaryName returns the release binary name for the current runtime.
func GetCurrentPlatformBinaryName() (string, error) {
	return GetPlatformBinaryName(runtime.GOOS, runtime.GOARCH)
}

// GetPlatformBinaryName maps GOOS and GOARCH to the release asset binary name convention.
func GetPlatformBinaryName(goos, goarch string) (string, error) {
	var normalizedArch string
	switch goarch {
	case "amd64", "x86_64":
		normalizedArch = "x86_64"
	case "arm64", "aarch64":
		normalizedArch = "aarch64"
	default:
		return "", fmt.Errorf("unsupported architecture: %s", goarch)
	}

	switch goos {
	case "darwin", "linux":
		return fmt.Sprintf("battery-%s-%s", goos, normalizedArch), nil
	case "windows":
		if normalizedArch != "x86_64" {
			return "", fmt.Errorf("unsupported architecture for windows: %s", goarch)
		}
		return fmt.Sprintf("battery-windows-%s.exe", normalizedArch), nil
	default:
		return "", fmt.Errorf("unsupported operating system: %s", goos)
	}
}
