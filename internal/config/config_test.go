package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twoBoots/battery/internal/config"
)

func TestLoadConfig_ReturnsDefaultWhenNoFile(t *testing.T) {
	tempDir := t.TempDir()

	cfg, err := config.LoadConfig(tempDir)
	require.NoError(t, err)
	assert.Equal(t, config.CurrentConfigVersion, cfg.Version)
	assert.Equal(t, config.StructureMultiRepo, cfg.Structure)
	assert.Empty(t, cfg.Barrels)
}

func TestSaveAndLoadConfig_Roundtrip(t *testing.T) {
	tempDir := t.TempDir()

	initial := config.BatteryConfig{
		Version:   "1.0.0",
		Structure: config.StructureMonorepo,
		Barrels: []config.BarrelConfig{
			{Name: "web", Path: "./apps/web"},
			{Name: "api", Path: "./apps/api", Type: config.BarrelTypeBarrel},
		},
	}

	savedPath, err := config.SaveConfig(initial, tempDir, false)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(tempDir, config.ConfigFilename), savedPath)

	loaded, err := config.LoadConfig(tempDir)
	require.NoError(t, err)
	assert.Equal(t, initial.Version, loaded.Version)
	assert.Equal(t, initial.Structure, loaded.Structure)
	assert.Equal(t, initial.Barrels, loaded.Barrels)
}

func TestLoadLocalConfig_ReturnsNilWhenNoFile(t *testing.T) {
	tempDir := t.TempDir()

	local, err := config.LoadLocalConfig(tempDir)
	require.NoError(t, err)
	assert.Nil(t, local)
}

func TestGetEffectiveConfig_MergesCanonicalAndLocal(t *testing.T) {
	tempDir := t.TempDir()

	canonical := config.BatteryConfig{
		Version:   "1.0.0",
		Structure: config.StructureMultiRepo,
		Barrels: []config.BarrelConfig{
			{Name: "auth", Path: "../auth-service"},
			{Name: "billing", Path: "../billing-service"},
		},
	}

	local := config.LocalBatteryConfig{
		Structure: config.StructureCustom,
		Barrels: []config.BarrelConfig{
			{Name: "auth", Path: "/custom/path/to/auth"},
			{Name: "dev-stub", Path: "../stubs/dev-stub", Type: config.BarrelTypeBarrel},
		},
	}

	_, err := config.SaveConfig(canonical, tempDir, false)
	require.NoError(t, err)

	_, err = config.SaveConfig(local, tempDir, true)
	require.NoError(t, err)

	effective, err := config.GetEffectiveConfig(tempDir)
	require.NoError(t, err)

	assert.Equal(t, config.StructureCustom, effective.Structure)
	assert.Len(t, effective.Barrels, 3)

	// Auth should be overridden by local
	var authBarrel *config.EffectiveBarrel
	var billingBarrel *config.EffectiveBarrel
	var stubBarrel *config.EffectiveBarrel

	for i := range effective.Barrels {
		b := &effective.Barrels[i]
		if b.Name == "auth" {
			authBarrel = b
		} else if b.Name == "billing" {
			billingBarrel = b
		} else if b.Name == "dev-stub" {
			stubBarrel = b
		}
	}

	require.NotNil(t, authBarrel)
	assert.Equal(t, "/custom/path/to/auth", authBarrel.Path)
	assert.Equal(t, "local", authBarrel.Source)

	require.NotNil(t, billingBarrel)
	assert.Equal(t, "../billing-service", billingBarrel.Path)
	assert.Equal(t, "canonical", billingBarrel.Source)

	require.NotNil(t, stubBarrel)
	assert.Equal(t, "../stubs/dev-stub", stubBarrel.Path)
	assert.Equal(t, "local", stubBarrel.Source)
}

func TestAddBarrel_CanonicalAndLocal(t *testing.T) {
	tempDir := t.TempDir()

	// Canonical add
	effective, err := config.AddBarrel(config.BarrelConfig{Path: "../core-repo", Name: "core"}, tempDir, false)
	require.NoError(t, err)
	assert.Len(t, effective.Barrels, 1)
	assert.Equal(t, "core", effective.Barrels[0].Name)

	// Duplicate canonical add rejection by name
	_, err = config.AddBarrel(config.BarrelConfig{Path: "../other-path", Name: "core"}, tempDir, false)
	assert.Error(t, err)

	// Duplicate canonical add rejection by path
	_, err = config.AddBarrel(config.BarrelConfig{Path: "../core-repo", Name: "core-2"}, tempDir, false)
	assert.Error(t, err)

	// Local add
	effectiveLocal, err := config.AddBarrel(config.BarrelConfig{Path: "../local-stub", Name: "stub", Type: config.BarrelTypeBarrel}, tempDir, true)
	require.NoError(t, err)
	assert.Len(t, effectiveLocal.Barrels, 2)

	// Duplicate local add rejection
	_, err = config.AddBarrel(config.BarrelConfig{Path: "../local-stub", Name: "stub2"}, tempDir, true)
	assert.Error(t, err)

	// Empty name and path rejection
	_, err = config.AddBarrel(config.BarrelConfig{Path: "   ", Name: "   "}, tempDir, false)
	assert.Error(t, err)
}

func TestRemoveBarrel_CanonicalAndLocal(t *testing.T) {
	tempDir := t.TempDir()

	_, err := config.AddBarrel(config.BarrelConfig{Path: "../repo-a", Name: "repo-a"}, tempDir, false)
	require.NoError(t, err)
	_, err = config.AddBarrel(config.BarrelConfig{Path: "../repo-b", Name: "repo-b"}, tempDir, false)
	require.NoError(t, err)
	_, err = config.AddBarrel(config.BarrelConfig{Path: "../local-a", Name: "local-a"}, tempDir, true)
	require.NoError(t, err)

	// Remove from canonical
	effective, err := config.RemoveBarrel("repo-a", tempDir, false)
	require.NoError(t, err)
	assert.Len(t, effective.Barrels, 2) // repo-b (canonical) + local-a (local)

	// Verify canonical state
	canonical, err := config.LoadConfig(tempDir)
	require.NoError(t, err)
	assert.Len(t, canonical.Barrels, 1)
	assert.Equal(t, "repo-b", canonical.Barrels[0].Name)

	// Remove from local
	effectiveAfterLocal, err := config.RemoveBarrel("local-a", tempDir, true)
	require.NoError(t, err)
	assert.Len(t, effectiveAfterLocal.Barrels, 1)

	// Removing non-existent barrel errors
	_, err = config.RemoveBarrel("non-existent", tempDir, false)
	assert.Error(t, err)

	_, err = config.RemoveBarrel("non-existent", tempDir, true)
	assert.Error(t, err)

	_, err = config.RemoveBarrel("   ", tempDir, false)
	assert.Error(t, err)
}

func TestInferBarrelName(t *testing.T) {
	assert.Equal(t, "service-auth", config.InferBarrelName("../service-auth"))
	assert.Equal(t, "ui", config.InferBarrelName("./packages/ui/"))
	assert.Equal(t, "billing", config.InferBarrelName("/var/repos/billing"))
	assert.Equal(t, "barrel", config.InferBarrelName(""))
	assert.Equal(t, "barrel", config.InferBarrelName("/"))
}

func TestInvalidConfigJSON(t *testing.T) {
	tempDir := t.TempDir()
	invalidJSON := []byte(`{ invalid json }`)

	err := os.WriteFile(filepath.Join(tempDir, config.ConfigFilename), invalidJSON, 0o644)
	require.NoError(t, err)

	_, err = config.LoadConfig(tempDir)
	assert.Error(t, err)

	err = os.WriteFile(filepath.Join(tempDir, config.LocalConfigFilename), invalidJSON, 0o644)
	require.NoError(t, err)

	_, err = config.LoadLocalConfig(tempDir)
	assert.Error(t, err)
}

func TestBarrelConfig_MetadataRoundtripAndDynamicFields(t *testing.T) {
	tempDir := t.TempDir()

	rawJSON := `{
  "version": "1.0.0",
  "structure": "multi-repo",
  "barrels": [
    {
      "name": "ros2-rust",
      "path": "../ros2-rust",
      "type": "barrel",
      "role": "ROS 2 robotics node bindings & hardware interfacing",
      "tech": "Rust / ROS 2 Humble / Cargo",
      "docs": "docs/barrels/ros2-rust.md",
      "jira": "ROBOT-123",
      "team": "robotics",
      "ci_badge": "https://ci.example.com/badge"
    }
  ]
}`

	err := os.WriteFile(filepath.Join(tempDir, config.ConfigFilename), []byte(rawJSON), 0o644)
	require.NoError(t, err)

	cfg, err := config.LoadConfig(tempDir)
	require.NoError(t, err)
	require.Len(t, cfg.Barrels, 1)

	b := cfg.Barrels[0]
	assert.Equal(t, "ros2-rust", b.Name)
	assert.Equal(t, "../ros2-rust", b.Path)
	assert.Equal(t, "ROS 2 robotics node bindings & hardware interfacing", b.Role)
	assert.Equal(t, "Rust / ROS 2 Humble / Cargo", b.Tech)
	assert.Equal(t, "docs/barrels/ros2-rust.md", b.Docs)
	assert.Equal(t, "ROBOT-123", b.Jira)

	// Verify custom dynamic fields are captured
	require.NotNil(t, b.Extra)
	assert.Contains(t, string(b.Extra["team"]), "robotics")
	assert.Contains(t, string(b.Extra["ci_badge"]), "https://ci.example.com/badge")

	// Save back and verify the JSON output still has team and ci_badge
	savedPath, err := config.SaveConfig(cfg, tempDir, false)
	require.NoError(t, err)

	savedData, err := os.ReadFile(savedPath)
	require.NoError(t, err)
	assert.Contains(t, string(savedData), `"role": "ROS 2 robotics node bindings & hardware interfacing"`)
	assert.Contains(t, string(savedData), `"tech": "Rust / ROS 2 Humble / Cargo"`)
	assert.Contains(t, string(savedData), `"docs": "docs/barrels/ros2-rust.md"`)
	assert.Contains(t, string(savedData), `"jira": "ROBOT-123"`)
	assert.Contains(t, string(savedData), `"team": "robotics"`)
	assert.Contains(t, string(savedData), `"ci_badge": "https://ci.example.com/badge"`)
}

func TestGetEffectiveConfig_PreservesMetadataAndCustomFields(t *testing.T) {
	tempDir := t.TempDir()

	canonical := config.BatteryConfig{
		Version:   "1.0.0",
		Structure: config.StructureMultiRepo,
		Barrels: []config.BarrelConfig{
			{
				Name: "firmware",
				Path: "../firmware-repo",
				Role: "Microcontroller C firmware",
				Tech: "C / ARM GCC",
				Docs: "docs/barrels/firmware.md",
				Jira: "EMB-45",
			},
		},
	}

	_, err := config.SaveConfig(canonical, tempDir, false)
	require.NoError(t, err)

	effective, err := config.GetEffectiveConfig(tempDir)
	require.NoError(t, err)
	require.Len(t, effective.Barrels, 1)

	eb := effective.Barrels[0]
	assert.Equal(t, "firmware", eb.Name)
	assert.Equal(t, "Microcontroller C firmware", eb.Role)
	assert.Equal(t, "C / ARM GCC", eb.Tech)
	assert.Equal(t, "docs/barrels/firmware.md", eb.Docs)
	assert.Equal(t, "EMB-45", eb.Jira)
}

func TestAddBarrel_WithMetadataAndPreserveOtherBarrels(t *testing.T) {
	tempDir := t.TempDir()

	// Initial barrel with custom attributes
	rawJSON := `{
  "version": "1.0.0",
  "structure": "multi-repo",
  "barrels": [
    {
      "name": "existing-barrel",
      "path": "../existing",
      "team": "core-infra"
    }
  ]
}`
	err := os.WriteFile(filepath.Join(tempDir, config.ConfigFilename), []byte(rawJSON), 0o644)
	require.NoError(t, err)

	// Add new barrel with metadata
	newBarrel := config.BarrelConfig{
		Name: "robot-driver",
		Path: "../driver",
		Role: "Hardware driver",
		Tech: "C++20",
		Docs: "docs/barrels/robot-driver.md",
		Jira: "DRV-1",
	}

	effective, err := config.AddBarrel(newBarrel, tempDir, false)
	require.NoError(t, err)
	require.Len(t, effective.Barrels, 2)

	// Verify new barrel has metadata
	var addedBarrel *config.EffectiveBarrel
	var existingBarrel *config.EffectiveBarrel
	for i := range effective.Barrels {
		if effective.Barrels[i].Name == "robot-driver" {
			addedBarrel = &effective.Barrels[i]
		}
		if effective.Barrels[i].Name == "existing-barrel" {
			existingBarrel = &effective.Barrels[i]
		}
	}

	require.NotNil(t, addedBarrel)
	assert.Equal(t, "Hardware driver", addedBarrel.Role)
	assert.Equal(t, "C++20", addedBarrel.Tech)
	assert.Equal(t, "docs/barrels/robot-driver.md", addedBarrel.Docs)
	assert.Equal(t, "DRV-1", addedBarrel.Jira)

	require.NotNil(t, existingBarrel)
	assert.Contains(t, string(existingBarrel.Extra["team"]), "core-infra")

	// Verify canonical file persisted properly
	loaded, err := config.LoadConfig(tempDir)
	require.NoError(t, err)
	require.Len(t, loaded.Barrels, 2)
	assert.Equal(t, "Hardware driver", loaded.Barrels[1].Role)
	assert.Equal(t, "C++20", loaded.Barrels[1].Tech)
	assert.Equal(t, "docs/barrels/robot-driver.md", loaded.Barrels[1].Docs)
	assert.Equal(t, "DRV-1", loaded.Barrels[1].Jira)
	assert.Contains(t, string(loaded.Barrels[0].Extra["team"]), "core-infra")
}
