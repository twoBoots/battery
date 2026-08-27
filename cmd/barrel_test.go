package cmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twoBoots/battery/cmd"
	"github.com/twoBoots/battery/internal/config"
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

func TestBarrelInitCmd_Scaffolding(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	barrelDir := filepath.Join(tempDir, "services", "payment")
	err := os.MkdirAll(barrelDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(barrelDir, "go.mod"), []byte("module payment\n\ngo 1.23.0\n"), 0644)
	require.NoError(t, err)

	buf := new(bytes.Buffer)
	cmd.RootCmd.SetOut(buf)
	cmd.RootCmd.SetErr(buf)
	cmd.RootCmd.SetArgs([]string{"barrel", "init", "./services/payment", "--framework", "Chi"})

	err = cmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Created Cooper tech stack")

	techFile := filepath.Join(barrelDir, ".cooper", "definition", "tech-stack.md")
	assert.FileExists(t, techFile)

	content, err := os.ReadFile(techFile)
	require.NoError(t, err)
	assert.Contains(t, string(content), "Chi")
	assert.Contains(t, string(content), "go test")
}

func TestBarrelAdd_WithMetadataFlags(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	buf := new(bytes.Buffer)
	cmd.RootCmd.SetOut(buf)
	cmd.RootCmd.SetArgs([]string{
		"barrel", "add", "../ros2-node",
		"--name", "ros2-node",
		"--role", "ROS 2 robotics node bindings",
		"--tech", "Rust / ROS 2 Humble",
		"--docs", "docs/barrels/ros2-node.md",
		"--jira", "ROBOT-42",
	})
	err := cmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Added barrel 'ros2-node'")

	cfg, err := config.LoadConfig(tempDir)
	require.NoError(t, err)
	require.Len(t, cfg.Barrels, 1)

	b := cfg.Barrels[0]
	assert.Equal(t, "ros2-node", b.Name)
	assert.Equal(t, "ROS 2 robotics node bindings", b.Role)
	assert.Equal(t, "Rust / ROS 2 Humble", b.Tech)
	assert.Equal(t, "docs/barrels/ros2-node.md", b.Docs)
	assert.Equal(t, "ROBOT-42", b.Jira)
}

func TestBarrelDocInitCmd_Scaffolding(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	// Scaffold profile using 'barrel doc init'
	buf := new(bytes.Buffer)
	cmd.RootCmd.SetOut(buf)
	cmd.RootCmd.SetErr(buf)
	cmd.RootCmd.SetArgs([]string{"barrel", "doc", "init", "robotics-firmware"})
	err := cmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Created barrel profile")

	profileFile := filepath.Join(tempDir, "docs", "barrels", "robotics-firmware.md")
	assert.FileExists(t, profileFile)

	content, err := os.ReadFile(profileFile)
	require.NoError(t, err)
	assert.Contains(t, string(content), "# Barrel Profile: robotics-firmware")

	// Attempting without --force should fail
	buf.Reset()
	cmd.RootCmd.SetArgs([]string{"barrel", "profile", "init", "robotics-firmware"})
	err = cmd.Execute()
	assert.Error(t, err)

	// Overwriting with --force
	buf.Reset()
	cmd.RootCmd.SetArgs([]string{"barrel", "doc", "init", "robotics-firmware", "--force"})
	err = cmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Updated barrel profile")
}
