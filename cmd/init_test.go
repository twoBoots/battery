package cmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twoboots/battery/cmd"
	"github.com/twoboots/battery/internal/config"
)

func TestInitCmd_NonInteractive_MultiRepo(t *testing.T) {
	parentDir := t.TempDir()

	currentBatteryDir := filepath.Join(parentDir, "battery")
	err := os.MkdirAll(currentBatteryDir, 0o755)
	require.NoError(t, err)

	// Create sibling repos
	sibling1 := filepath.Join(parentDir, "repo-1")
	err = os.MkdirAll(filepath.Join(sibling1, ".git"), 0o755)
	require.NoError(t, err)

	sibling2 := filepath.Join(parentDir, "repo-2")
	err = os.MkdirAll(filepath.Join(sibling2, ".cooper"), 0o755)
	require.NoError(t, err)

	origDir, _ := os.Getwd()
	_ = os.Chdir(currentBatteryDir)
	defer func() { _ = os.Chdir(origDir) }()

	buf := new(bytes.Buffer)
	cmd.RootCmd.SetOut(buf)
	cmd.RootCmd.SetArgs([]string{"init", "--yes"})

	err = cmd.Execute()
	require.NoError(t, err)

	out := buf.String()
	assert.Contains(t, out, "Initializing Battery Configuration")
	assert.Contains(t, out, "Topology selected: multi-repo")
	assert.Contains(t, out, "Discovered 2 barrels: repo-1, repo-2")
	assert.Contains(t, out, "Battery initialized successfully in .batteryrc")

	cfg, err := config.LoadConfig(currentBatteryDir)
	require.NoError(t, err)
	assert.Equal(t, config.StructureMultiRepo, cfg.Structure)
	assert.Len(t, cfg.Barrels, 2)
}

func TestInitCmd_NonInteractive_LocalMonorepo(t *testing.T) {
	tempDir := t.TempDir()

	// Monorepo package
	pkgDir := filepath.Join(tempDir, "packages", "core")
	err := os.MkdirAll(pkgDir, 0o755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(pkgDir, "package.json"), []byte("{}"), 0o644)
	require.NoError(t, err)

	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	buf := new(bytes.Buffer)
	cmd.RootCmd.SetOut(buf)
	cmd.RootCmd.SetArgs([]string{"init", "--non-interactive", "--local", "--structure", "monorepo"})

	err = cmd.Execute()
	require.NoError(t, err)

	out := buf.String()
	assert.Contains(t, out, "Battery initialized successfully in .batteryrc.local")

	local, err := config.LoadLocalConfig(tempDir)
	require.NoError(t, err)
	require.NotNil(t, local)
	assert.Equal(t, config.StructureMonorepo, local.Structure)
	assert.Len(t, local.Barrels, 1)
	assert.Equal(t, "core", local.Barrels[0].Name)
}

func TestInitCmd_NonInteractive_NoBarrels(t *testing.T) {
	tempDir := t.TempDir()

	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	buf := new(bytes.Buffer)
	cmd.RootCmd.SetOut(buf)
	cmd.RootCmd.SetArgs([]string{"init", "--non-interactive", "--structure", "custom"})

	err := cmd.Execute()
	require.NoError(t, err)

	out := buf.String()
	assert.Contains(t, out, "No candidate barrels automatically discovered")
	assert.Contains(t, out, "Structure: custom")
}

func TestInitCmd_NonInteractive_PreserveExistingConfig(t *testing.T) {
	parentDir := t.TempDir()
	currentBatteryDir := filepath.Join(parentDir, "battery")
	err := os.MkdirAll(currentBatteryDir, 0o755)
	require.NoError(t, err)

	// Create sibling repo
	sibling1 := filepath.Join(parentDir, "repo-discovered")
	err = os.MkdirAll(filepath.Join(sibling1, ".git"), 0o755)
	require.NoError(t, err)

	// Pre-create existing .batteryrc with a custom barrel
	existingCfg := config.BatteryConfig{
		Version:   "1.0.0",
		Structure: config.StructureCustom,
		Barrels: []config.BarrelConfig{
			{Name: "my-existing-barrel", Path: "../my-existing-barrel"},
		},
	}
	_, err = config.SaveConfig(existingCfg, currentBatteryDir, false)
	require.NoError(t, err)

	origDir, _ := os.Getwd()
	_ = os.Chdir(currentBatteryDir)
	defer func() { _ = os.Chdir(origDir) }()

	buf := new(bytes.Buffer)
	cmd.RootCmd.SetOut(buf)
	cmd.RootCmd.SetArgs([]string{"init", "--non-interactive"})

	err = cmd.Execute()
	require.NoError(t, err)

	out := buf.String()
	assert.Contains(t, out, "Existing configuration detected in .batteryrc; preserving current configuration.")
	assert.Contains(t, out, "Structure: custom")
	assert.Contains(t, out, "Barrels  : 1 registered")

	// Ensure config was NOT overwritten with repo-discovered
	cfg, err := config.LoadConfig(currentBatteryDir)
	require.NoError(t, err)
	assert.Equal(t, config.StructureCustom, cfg.Structure)
	require.Len(t, cfg.Barrels, 1)
	assert.Equal(t, "my-existing-barrel", cfg.Barrels[0].Name)
}

func TestInitCmd_NonInteractive_ForceOverwrite(t *testing.T) {
	parentDir := t.TempDir()
	currentBatteryDir := filepath.Join(parentDir, "battery")
	err := os.MkdirAll(currentBatteryDir, 0o755)
	require.NoError(t, err)

	// Create sibling repo
	sibling1 := filepath.Join(parentDir, "repo-discovered")
	err = os.MkdirAll(filepath.Join(sibling1, ".git"), 0o755)
	require.NoError(t, err)

	// Pre-create existing .batteryrc
	existingCfg := config.BatteryConfig{
		Version:   "1.0.0",
		Structure: config.StructureCustom,
		Barrels: []config.BarrelConfig{
			{Name: "old-barrel", Path: "../old-barrel"},
		},
	}
	_, err = config.SaveConfig(existingCfg, currentBatteryDir, false)
	require.NoError(t, err)

	origDir, _ := os.Getwd()
	_ = os.Chdir(currentBatteryDir)
	defer func() { _ = os.Chdir(origDir) }()

	buf := new(bytes.Buffer)
	cmd.RootCmd.SetOut(buf)
	cmd.RootCmd.SetArgs([]string{"init", "--non-interactive", "--force"})

	err = cmd.Execute()
	require.NoError(t, err)

	out := buf.String()
	assert.Contains(t, out, "Overwriting existing configuration")
	assert.Contains(t, out, "repo-discovered")

	// Verify config was overwritten with discovered barrel
	cfg, err := config.LoadConfig(currentBatteryDir)
	require.NoError(t, err)
	assert.Equal(t, config.StructureMultiRepo, cfg.Structure)
	require.Len(t, cfg.Barrels, 1)
	assert.Equal(t, "repo-discovered", cfg.Barrels[0].Name)
}

func TestInitCmd_NonInteractive_Local_PreserveExistingConfig(t *testing.T) {
	parentDir := t.TempDir()
	currentBatteryDir := filepath.Join(parentDir, "battery")
	err := os.MkdirAll(currentBatteryDir, 0o755)
	require.NoError(t, err)

	// Pre-create existing .batteryrc.local
	existingLocalCfg := config.LocalBatteryConfig{
		Structure: config.StructureMonorepo,
		Barrels: []config.BarrelConfig{
			{Name: "local-barrel", Path: "./local-barrel"},
		},
	}
	_, err = config.SaveConfig(existingLocalCfg, currentBatteryDir, true)
	require.NoError(t, err)

	origDir, _ := os.Getwd()
	_ = os.Chdir(currentBatteryDir)
	defer func() { _ = os.Chdir(origDir) }()

	buf := new(bytes.Buffer)
	cmd.RootCmd.SetOut(buf)
	cmd.RootCmd.SetArgs([]string{"init", "--non-interactive", "--local"})

	err = cmd.Execute()
	require.NoError(t, err)

	out := buf.String()
	assert.Contains(t, out, "Existing configuration detected in .batteryrc.local; preserving current configuration.")
	assert.Contains(t, out, "Structure: monorepo")
	assert.Contains(t, out, "Barrels  : 1 registered")

	local, err := config.LoadLocalConfig(currentBatteryDir)
	require.NoError(t, err)
	require.NotNil(t, local)
	assert.Equal(t, config.StructureMonorepo, local.Structure)
	require.Len(t, local.Barrels, 1)
	assert.Equal(t, "local-barrel", local.Barrels[0].Name)
}

func TestInitCmd_NonInteractive_OverwriteAlias(t *testing.T) {
	parentDir := t.TempDir()
	currentBatteryDir := filepath.Join(parentDir, "battery")
	err := os.MkdirAll(currentBatteryDir, 0o755)
	require.NoError(t, err)

	sibling1 := filepath.Join(parentDir, "repo-discovered")
	err = os.MkdirAll(filepath.Join(sibling1, ".git"), 0o755)
	require.NoError(t, err)

	existingCfg := config.BatteryConfig{
		Version:   "1.0.0",
		Structure: config.StructureCustom,
		Barrels: []config.BarrelConfig{
			{Name: "old-barrel", Path: "../old-barrel"},
		},
	}
	_, err = config.SaveConfig(existingCfg, currentBatteryDir, false)
	require.NoError(t, err)

	origDir, _ := os.Getwd()
	_ = os.Chdir(currentBatteryDir)
	defer func() { _ = os.Chdir(origDir) }()

	buf := new(bytes.Buffer)
	cmd.RootCmd.SetOut(buf)
	cmd.RootCmd.SetArgs([]string{"init", "--non-interactive", "--overwrite"})

	err = cmd.Execute()
	require.NoError(t, err)

	out := buf.String()
	assert.Contains(t, out, "Overwriting existing configuration")
	assert.Contains(t, out, "repo-discovered")

	cfg, err := config.LoadConfig(currentBatteryDir)
	require.NoError(t, err)
	assert.Equal(t, config.StructureMultiRepo, cfg.Structure)
	require.Len(t, cfg.Barrels, 1)
	assert.Equal(t, "repo-discovered", cfg.Barrels[0].Name)
}
