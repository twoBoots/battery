package cmd_test

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twoboots/battery/internal/config"
)

func getInstallScriptPath(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok)
	repoRoot := filepath.Dir(filepath.Dir(filename))
	installScript := filepath.Join(repoRoot, "install.sh")
	require.FileExists(t, installScript)
	return installScript
}

func TestInstallScriptE2EPreserveExistingConfig(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Initialize git in temp repo
	gitCmd := exec.Command("git", "init")
	gitCmd.Dir = tempDir
	err := gitCmd.Run()
	require.NoError(t, err)

	// 2. Pre-create custom .batteryrc
	initialCfg := config.BatteryConfig{
		Version:   "1.0.0",
		Structure: config.StructureCustom,
		Barrels: []config.BarrelConfig{
			{Name: "my-custom-barrel", Path: "../custom-path"},
		},
	}
	_, err = config.SaveConfig(initialCfg, tempDir, false)
	require.NoError(t, err)

	// 3. Run install.sh in non-interactive mode
	installScript := getInstallScriptPath(t)
	installCmd := exec.Command("bash", installScript, tempDir, "--non-interactive")
	output, err := installCmd.CombinedOutput()
	require.NoError(t, err, "install.sh failed: %s", string(output))

	outStr := string(output)
	assert.Contains(t, outStr, "Existing configuration detected in .batteryrc; preserving current configuration.")
	assert.Contains(t, outStr, "Structure: custom")
	assert.Contains(t, outStr, "Barrels  : 1 registered")

	// 4. Verify that .batteryrc was preserved
	preservedCfg, err := config.LoadConfig(tempDir)
	require.NoError(t, err)
	assert.Equal(t, config.StructureCustom, preservedCfg.Structure)
	require.Len(t, preservedCfg.Barrels, 1)
	assert.Equal(t, "my-custom-barrel", preservedCfg.Barrels[0].Name)
	assert.Equal(t, "../custom-path", preservedCfg.Barrels[0].Path)
}

func TestInstallScriptE2EForceOverwrite(t *testing.T) {
	parentDir := t.TempDir()
	tempDir := filepath.Join(parentDir, "my-battery-repo")
	err := exec.Command("mkdir", "-p", tempDir).Run()
	require.NoError(t, err)

	// Create sibling repo so discovery finds something
	sibling := filepath.Join(parentDir, "sibling-repo")
	err = exec.Command("mkdir", "-p", filepath.Join(sibling, ".git")).Run()
	require.NoError(t, err)

	gitCmd := exec.Command("git", "init")
	gitCmd.Dir = tempDir
	err = gitCmd.Run()
	require.NoError(t, err)

	initialCfg := config.BatteryConfig{
		Version:   "1.0.0",
		Structure: config.StructureCustom,
		Barrels: []config.BarrelConfig{
			{Name: "old-barrel", Path: "../old-barrel"},
		},
	}
	_, err = config.SaveConfig(initialCfg, tempDir, false)
	require.NoError(t, err)

	installScript := getInstallScriptPath(t)
	installCmd := exec.Command("bash", installScript, tempDir, "--non-interactive", "--force")
	output, err := installCmd.CombinedOutput()
	require.NoError(t, err, "install.sh --force failed: %s", string(output))

	outStr := string(output)
	assert.Contains(t, outStr, "Overwriting existing configuration in .batteryrc...")
	assert.Contains(t, outStr, "sibling-repo")

	overwrittenCfg, err := config.LoadConfig(tempDir)
	require.NoError(t, err)
	assert.Equal(t, config.StructureMultiRepo, overwrittenCfg.Structure)
	require.Len(t, overwrittenCfg.Barrels, 1)
	assert.Equal(t, "sibling-repo", overwrittenCfg.Barrels[0].Name)
}
