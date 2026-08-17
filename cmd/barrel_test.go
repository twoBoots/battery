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

func TestBarrelListCmd_Empty(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	buf := new(bytes.Buffer)
	cmd.RootCmd.SetOut(buf)
	cmd.RootCmd.SetErr(buf)
	cmd.RootCmd.SetArgs([]string{"barrel", "list"})

	err := cmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Registered Barrels (0)")
	assert.Contains(t, buf.String(), "No barrels registered yet")
}

func TestBarrelAddAndRemoveCmd(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	// Add barrel
	buf := new(bytes.Buffer)
	cmd.RootCmd.SetOut(buf)
	cmd.RootCmd.SetArgs([]string{"barrel", "add", "../auth-service", "--name", "auth"})
	err := cmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Added barrel 'auth' (../auth-service) to .batteryrc")

	// Verify in config
	cfg, err := config.LoadConfig(tempDir)
	require.NoError(t, err)
	assert.Len(t, cfg.Barrels, 1)
	assert.Equal(t, "auth", cfg.Barrels[0].Name)

	// List barrel
	buf.Reset()
	cmd.RootCmd.SetArgs([]string{"barrel", "list"})
	err = cmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Registered Barrels (1)")
	assert.Contains(t, buf.String(), "auth [canonical]")

	// Remove barrel
	buf.Reset()
	cmd.RootCmd.SetArgs([]string{"barrel", "remove", "auth"})
	err = cmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Removed barrel 'auth' from .batteryrc")

	// Verify empty
	cfg, err = config.LoadConfig(tempDir)
	require.NoError(t, err)
	assert.Empty(t, cfg.Barrels)
}

func TestBarrelAddLocalCmd(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	targetStub := filepath.Join(tempDir, "stub-repo")
	err := os.MkdirAll(targetStub, 0o755)
	require.NoError(t, err)

	buf := new(bytes.Buffer)
	cmd.RootCmd.SetOut(buf)
	cmd.RootCmd.SetArgs([]string{"barrel", "add", "./stub-repo", "--local", "--type", "battery"})
	err = cmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "to .batteryrc.local")

	local, err := config.LoadLocalConfig(tempDir)
	require.NoError(t, err)
	require.NotNil(t, local)
	assert.Len(t, local.Barrels, 1)
	assert.Equal(t, config.BarrelTypeBattery, local.Barrels[0].Type)
}

func TestBarrelCmd_DefaultList(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	buf := new(bytes.Buffer)
	cmd.RootCmd.SetOut(buf)
	cmd.RootCmd.SetArgs([]string{"barrel"})
	err := cmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Registered Barrels (0)")
}

func TestBarrelAdd_MissingPath(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	cmd.RootCmd.SetArgs([]string{"barrel", "add"})
	err := cmd.Execute()
	assert.Error(t, err)
}

func TestBarrelRemove_NotFound(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	cmd.RootCmd.SetArgs([]string{"barrel", "remove", "unknown"})
	err := cmd.Execute()
	assert.Error(t, err)
}
