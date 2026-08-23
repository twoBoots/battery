package discovery_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twoBoots/battery/internal/config"
	"github.com/twoBoots/battery/internal/discovery"
)

func TestHasProjectMarker(t *testing.T) {
	tempDir := t.TempDir()

	assert.False(t, discovery.HasProjectMarker(tempDir))

	err := os.WriteFile(filepath.Join(tempDir, "package.json"), []byte("{}"), 0o644)
	require.NoError(t, err)
	assert.True(t, discovery.HasProjectMarker(tempDir))
}

func TestDiscoverSiblingBarrels(t *testing.T) {
	parentDir := t.TempDir()

	currentBatteryDir := filepath.Join(parentDir, "battery")
	err := os.MkdirAll(currentBatteryDir, 0o755)
	require.NoError(t, err)

	// Sibling A with .git
	siblingA := filepath.Join(parentDir, "auth-service")
	err = os.MkdirAll(filepath.Join(siblingA, ".git"), 0o755)
	require.NoError(t, err)

	// Sibling B with .cooper
	siblingB := filepath.Join(parentDir, "web-client")
	err = os.MkdirAll(filepath.Join(siblingB, ".cooper"), 0o755)
	require.NoError(t, err)

	// Sibling C without project marker
	siblingC := filepath.Join(parentDir, "random-folder")
	err = os.MkdirAll(siblingC, 0o755)
	require.NoError(t, err)

	// Sibling hidden directory with marker (should be skipped)
	hiddenDir := filepath.Join(parentDir, ".hidden-repo")
	err = os.MkdirAll(filepath.Join(hiddenDir, ".git"), 0o755)
	require.NoError(t, err)

	barrels := discovery.DiscoverSiblingBarrels(currentBatteryDir)
	assert.Len(t, barrels, 2)
	assert.Equal(t, "auth-service", barrels[0].Name)
	assert.Equal(t, "../auth-service", barrels[0].Path)
	assert.Equal(t, "web-client", barrels[1].Name)
	assert.Equal(t, "../web-client", barrels[1].Path)
}

func TestDiscoverMonorepoBarrels(t *testing.T) {
	monorepoDir := t.TempDir()

	// apps/web with package.json
	webPkg := filepath.Join(monorepoDir, "apps", "web")
	err := os.MkdirAll(webPkg, 0o755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(webPkg, "package.json"), []byte("{}"), 0o644)
	require.NoError(t, err)

	// packages/ui with Cargo.toml
	uiPkg := filepath.Join(monorepoDir, "packages", "ui")
	err = os.MkdirAll(uiPkg, 0o755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(uiPkg, "Cargo.toml"), []byte(""), 0o644)
	require.NoError(t, err)

	// packages/empty (no marker)
	emptyPkg := filepath.Join(monorepoDir, "packages", "empty")
	err = os.MkdirAll(emptyPkg, 0o755)
	require.NoError(t, err)

	barrels := discovery.DiscoverMonorepoBarrels(monorepoDir)
	assert.Len(t, barrels, 2)
	assert.Equal(t, "./packages/ui", barrels[0].Path)
	assert.Equal(t, "ui", barrels[0].Name)
	assert.Equal(t, "./apps/web", barrels[1].Path)
	assert.Equal(t, "web", barrels[1].Name)
}

func TestDetectProjectStructure_And_DiscoverCandidateBarrels(t *testing.T) {
	// 1. Monorepo
	monorepoDir := t.TempDir()
	appPkg := filepath.Join(monorepoDir, "apps", "api")
	err := os.MkdirAll(appPkg, 0o755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(appPkg, "go.mod"), []byte("module api\n"), 0o644)
	require.NoError(t, err)

	res := discovery.DiscoverCandidateBarrels(monorepoDir)
	assert.Equal(t, config.StructureMonorepo, res.Structure)
	assert.Len(t, res.Barrels, 1)

	// 2. Custom fallback
	emptyDir := t.TempDir()
	customDir := filepath.Join(emptyDir, "isolated")
	err = os.MkdirAll(customDir, 0o755)
	require.NoError(t, err)

	resCustom := discovery.DiscoverCandidateBarrels(customDir)
	assert.Equal(t, config.StructureCustom, resCustom.Structure)
	assert.Empty(t, resCustom.Barrels)
}
