package updater

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPlatformBinaryName(t *testing.T) {
	tests := []struct {
		name     string
		goos     string
		goarch   string
		expected string
		wantErr  bool
	}{
		{
			name:     "darwin amd64",
			goos:     "darwin",
			goarch:   "amd64",
			expected: "battery-darwin-x86_64",
		},
		{
			name:     "darwin arm64",
			goos:     "darwin",
			goarch:   "arm64",
			expected: "battery-darwin-aarch64",
		},
		{
			name:     "linux amd64",
			goos:     "linux",
			goarch:   "amd64",
			expected: "battery-linux-x86_64",
		},
		{
			name:     "linux arm64",
			goos:     "linux",
			goarch:   "arm64",
			expected: "battery-linux-aarch64",
		},
		{
			name:     "windows amd64",
			goos:     "windows",
			goarch:   "amd64",
			expected: "battery-windows-x86_64.exe",
		},
		{
			name:    "unsupported os",
			goos:    "freebsd",
			goarch:  "amd64",
			wantErr: true,
		},
		{
			name:    "unsupported arch",
			goos:    "linux",
			goarch:  "mips",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPlatformBinaryName(tt.goos, tt.goarch)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestGetCurrentPlatformBinaryName(t *testing.T) {
	name, err := GetCurrentPlatformBinaryName()
	assert.NoError(t, err)
	assert.NotEmpty(t, name)
	assert.Contains(t, name, "battery-")
}
