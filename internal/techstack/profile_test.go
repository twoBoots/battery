package techstack_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twoBoots/battery/internal/config"
	"github.com/twoBoots/battery/internal/techstack"
)

func TestResolveBarrelProfile_DocsBarrels(t *testing.T) {
	tempDir := t.TempDir()
	docsBarrelsDir := filepath.Join(tempDir, "docs", "barrels")
	err := os.MkdirAll(docsBarrelsDir, 0o755)
	require.NoError(t, err)

	profileContent := `# Barrel Profile: ros2-rust

## Role & Responsibilities
ROS 2 robotics node bindings & hardware interfacing

## Tech Stack & Runtime
- Rust 1.78+
- ROS 2 Humble
- Cargo
`
	profilePath := filepath.Join(docsBarrelsDir, "ros2-rust.md")
	err = os.WriteFile(profilePath, []byte(profileContent), 0o644)
	require.NoError(t, err)

	barrel := config.EffectiveBarrel{
		Name: "ros2-rust",
		Path: "../ros2-rust",
	}

	info := techstack.ResolveBarrelProfile(tempDir, barrel)
	assert.True(t, info.Exists)
	assert.Equal(t, profilePath, info.FilePath)
	assert.Equal(t, profileContent, info.Content)
	assert.Equal(t, "Rust 1.78+ | ROS 2 Humble | Cargo", info.Summary)
}

func TestResolveBarrelProfile_CooperBarrelsFallback(t *testing.T) {
	tempDir := t.TempDir()
	cooperBarrelsDir := filepath.Join(tempDir, ".cooper", "barrels")
	err := os.MkdirAll(cooperBarrelsDir, 0o755)
	require.NoError(t, err)

	profileContent := `# Profile: embedded-firmware
* C / ARM GCC
* FreeRTOS
`
	profilePath := filepath.Join(cooperBarrelsDir, "embedded-firmware.md")
	err = os.WriteFile(profilePath, []byte(profileContent), 0o644)
	require.NoError(t, err)

	barrel := config.EffectiveBarrel{
		Name: "embedded-firmware",
		Path: "../firmware",
	}

	info := techstack.ResolveBarrelProfile(tempDir, barrel)
	assert.True(t, info.Exists)
	assert.Equal(t, profilePath, info.FilePath)
	assert.Equal(t, "C / ARM GCC | FreeRTOS", info.Summary)
}

func TestResolveBarrelProfile_CustomDocsPath(t *testing.T) {
	tempDir := t.TempDir()
	customDir := filepath.Join(tempDir, "architecture", "profiles")
	err := os.MkdirAll(customDir, 0o755)
	require.NoError(t, err)

	profileContent := "Custom Profile Summary Content"
	profilePath := filepath.Join(customDir, "custom.md")
	err = os.WriteFile(profilePath, []byte(profileContent), 0o644)
	require.NoError(t, err)

	barrel := config.EffectiveBarrel{
		Name: "custom-barrel",
		Path: "../custom-barrel",
		Docs: "architecture/profiles/custom.md",
	}

	info := techstack.ResolveBarrelProfile(tempDir, barrel)
	assert.True(t, info.Exists)
	assert.Equal(t, profilePath, info.FilePath)
	assert.Equal(t, "Custom Profile Summary Content", info.Summary)
}

func TestResolveBarrelProfile_NotFound(t *testing.T) {
	tempDir := t.TempDir()
	barrel := config.EffectiveBarrel{
		Name: "missing-barrel",
		Path: "../missing",
	}

	info := techstack.ResolveBarrelProfile(tempDir, barrel)
	assert.False(t, info.Exists)
	assert.Empty(t, info.FilePath)
	assert.Empty(t, info.Content)
	assert.Empty(t, info.Summary)
}

func TestScaffoldBarrelProfile_SuccessAndOverwrite(t *testing.T) {
	tempDir := t.TempDir()

	path, err := techstack.ScaffoldBarrelProfile(tempDir, "driver", false)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(tempDir, "docs", "barrels", "driver.md"), path)

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(data), "# Barrel Profile: driver")
	assert.Contains(t, string(data), "## Role & Responsibilities")
	assert.Contains(t, string(data), "## Tech Stack & Runtime")
	assert.Contains(t, string(data), "## Development & Build Commands")
	assert.Contains(t, string(data), "## AI Agent Guidelines")

	// Attempting without force should fail
	_, err = techstack.ScaffoldBarrelProfile(tempDir, "driver", false)
	assert.Error(t, err)

	// Attempting with force should succeed
	path2, err := techstack.ScaffoldBarrelProfile(tempDir, "driver", true)
	require.NoError(t, err)
	assert.Equal(t, path, path2)
}
