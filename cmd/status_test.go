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

func TestStatusCmd_Empty(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	buf := new(bytes.Buffer)
	cmd.RootCmd.SetOut(buf)
	cmd.RootCmd.SetArgs([]string{"status"})

	err := cmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Battery Workspace Status")
	assert.Contains(t, buf.String(), "Total Barrels : 0")
	assert.Contains(t, buf.String(), "No barrels currently configured")
}

func TestStatusCmd_WithBarrels(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	// Create connected barrel folder with Cooper tech stack
	connectedDir := filepath.Join(tempDir, "connected-barrel")
	cooperDefDir := filepath.Join(connectedDir, ".cooper", "definition")
	err := os.MkdirAll(cooperDefDir, 0o755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(cooperDefDir, "tech-stack.md"), []byte("- Go 1.23\n- Cobra"), 0o644)
	require.NoError(t, err)

	// Save configuration with 1 connected and 1 missing barrel
	cfg := config.BatteryConfig{
		Version:   "1.0.0",
		Structure: config.StructureMultiRepo,
		Barrels: []config.BarrelConfig{
			{Name: "connected", Path: "./connected-barrel"},
			{Name: "missing", Path: "./missing-barrel"},
		},
	}
	_, err = config.SaveConfig(cfg, tempDir, false)
	require.NoError(t, err)

	buf := new(bytes.Buffer)
	cmd.RootCmd.SetOut(buf)
	cmd.RootCmd.SetArgs([]string{"status"})

	err = cmd.Execute()
	require.NoError(t, err)

	out := buf.String()
	assert.Contains(t, out, "Total Barrels : 2")
	assert.Contains(t, out, "🟢 connected (./connected-barrel)")
	assert.Contains(t, out, "Tech Stack : Go 1.23 | Cobra")
	assert.Contains(t, out, "🔴 missing (./missing-barrel)")
	assert.Contains(t, out, "Summary: 1/2 barrels connected.")
}
