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

func TestResolveBarrelContext_CooperPriority(t *testing.T) {
	tempDir := t.TempDir()
	barrelDir := filepath.Join(tempDir, "cooper-repo")
	cooperDir := filepath.Join(barrelDir, ".cooper", "definition")
	err := os.MkdirAll(cooperDir, 0o755)
	require.NoError(t, err)

	content := "# Tech Stack\n- Go 1.27\n- Cobra\n"
	err = os.WriteFile(filepath.Join(cooperDir, "tech-stack.md"), []byte(content), 0o644)
	require.NoError(t, err)

	barrel := config.EffectiveBarrel{
		Name: "cooper-repo",
		Path: "./cooper-repo",
		Tech: "Metadata Tech",
		Role: "Metadata Role",
	}

	ctx := techstack.ResolveBarrelContext(tempDir, barrelDir, barrel)
	assert.Equal(t, "cooper", ctx.Source)
	assert.True(t, ctx.HasCooperSpec)
	assert.Equal(t, "Go 1.27 | Cobra", ctx.Summary)
}

func TestResolveBarrelContext_MetadataFallback(t *testing.T) {
	tempDir := t.TempDir()
	barrelDir := filepath.Join(tempDir, "meta-repo")
	err := os.MkdirAll(barrelDir, 0o755)
	require.NoError(t, err)

	// Both Tech and Role
	barrel1 := config.EffectiveBarrel{
		Name: "meta-repo",
		Path: "./meta-repo",
		Tech: "Rust / ROS 2",
		Role: "Hardware Controller",
	}
	ctx1 := techstack.ResolveBarrelContext(tempDir, barrelDir, barrel1)
	assert.Equal(t, "metadata", ctx1.Source)
	assert.False(t, ctx1.HasCooperSpec)
	assert.Equal(t, "Rust / ROS 2 (Hardware Controller)", ctx1.Summary)

	// Only Tech
	barrel2 := config.EffectiveBarrel{
		Name: "meta-repo",
		Path: "./meta-repo",
		Tech: "Rust / ROS 2",
	}
	ctx2 := techstack.ResolveBarrelContext(tempDir, barrelDir, barrel2)
	assert.Equal(t, "metadata", ctx2.Source)
	assert.Equal(t, "Rust / ROS 2", ctx2.Summary)

	// Only Role
	barrel3 := config.EffectiveBarrel{
		Name: "meta-repo",
		Path: "./meta-repo",
		Role: "Hardware Controller",
	}
	ctx3 := techstack.ResolveBarrelContext(tempDir, barrelDir, barrel3)
	assert.Equal(t, "metadata", ctx3.Source)
	assert.Equal(t, "Hardware Controller", ctx3.Summary)
}

func TestResolveBarrelContext_ProfileDocFallback(t *testing.T) {
	tempDir := t.TempDir()
	barrelDir := filepath.Join(tempDir, "doc-repo")
	err := os.MkdirAll(barrelDir, 0o755)
	require.NoError(t, err)

	docsDir := filepath.Join(tempDir, "docs", "barrels")
	err = os.MkdirAll(docsDir, 0o755)
	require.NoError(t, err)

	profileContent := "# Barrel Profile: doc-repo\n- Python 3.12\n- PyTorch\n"
	profilePath := filepath.Join(docsDir, "doc-repo.md")
	err = os.WriteFile(profilePath, []byte(profileContent), 0o644)
	require.NoError(t, err)

	barrel := config.EffectiveBarrel{
		Name: "doc-repo",
		Path: "./doc-repo",
	}

	ctx := techstack.ResolveBarrelContext(tempDir, barrelDir, barrel)
	assert.Equal(t, "profile", ctx.Source)
	assert.False(t, ctx.HasCooperSpec)
	assert.True(t, ctx.HasProfile)
	assert.Equal(t, profilePath, ctx.ProfilePath)
	assert.Equal(t, "Python 3.12 | PyTorch", ctx.Summary)
}

func TestResolveBarrelContext_DefaultFallback(t *testing.T) {
	tempDir := t.TempDir()
	barrelDir := filepath.Join(tempDir, "empty-repo")
	err := os.MkdirAll(barrelDir, 0o755)
	require.NoError(t, err)

	barrel := config.EffectiveBarrel{
		Name: "empty-repo",
		Path: "./empty-repo",
	}

	ctx := techstack.ResolveBarrelContext(tempDir, barrelDir, barrel)
	assert.Equal(t, "none", ctx.Source)
	assert.False(t, ctx.HasCooperSpec)
	assert.False(t, ctx.HasProfile)
	assert.Equal(t, "No Cooper tech-stack or profile defined", ctx.Summary)
}
