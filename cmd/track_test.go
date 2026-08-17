package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twoboots/battery/internal/config"
	"github.com/twoboots/battery/internal/track"
)

func TestTrackInitCmd(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup config with 2 barrels
	cfg := config.BatteryConfig{
		Version:   "1.0.0",
		Structure: config.StructureMultiRepo,
		Barrels: []config.BarrelConfig{
			{Name: "folder-a", Path: "./folder-a"},
			{Name: "folder-b", Path: "./folder-b"},
		},
	}
	_, err := config.SaveConfig(&cfg, tmpDir, false)
	require.NoError(t, err)

	buf := new(bytes.Buffer)
	err = runTrackInit(buf, tmpDir, "my-feature", "My Feature Title", "feature", []string{"folder-a", "folder-b"}, []string{"auth"})
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Initialized track 'my-feature'")

	// Verify track directory exists
	trackDir := filepath.Join(tmpDir, ".cooper", "active", "my-feature")
	assert.FileExists(t, filepath.Join(trackDir, "metadata.json"))
}

func TestTrackDispatchCmd(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup barrels
	barrelADir := filepath.Join(tmpDir, "folder-a")
	require.NoError(t, os.MkdirAll(barrelADir, 0755))

	cfg := config.BatteryConfig{
		Version:   "1.0.0",
		Structure: config.StructureMultiRepo,
		Barrels: []config.BarrelConfig{
			{Name: "folder-a", Path: "./folder-a"},
		},
	}
	_, err := config.SaveConfig(&cfg, tmpDir, false)
	require.NoError(t, err)

	// Init track in battery
	_, err = track.InitTrack(tmpDir, "my-feature", []string{"folder-a"}, track.InitTrackOptions{
		Capabilities: []string{"auth"},
	})
	require.NoError(t, err)

	buf := new(bytes.Buffer)
	err = runTrackDispatch(buf, tmpDir, "my-feature", false)
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Dispatched track 'my-feature'")

	// Verify folder-a contains spec-deltas but NO plan.md
	assert.FileExists(t, filepath.Join(barrelADir, ".cooper", "active", "my-feature", "spec-deltas", "auth", "spec.md"))
	assert.NoFileExists(t, filepath.Join(barrelADir, ".cooper", "active", "my-feature", "plan.md"))
}

func TestTrackStatusCmd(t *testing.T) {
	tmpDir := t.TempDir()

	barrelADir := filepath.Join(tmpDir, "folder-a")
	require.NoError(t, os.MkdirAll(barrelADir, 0755))

	cfg := config.BatteryConfig{
		Version:   "1.0.0",
		Structure: config.StructureMultiRepo,
		Barrels: []config.BarrelConfig{
			{Name: "folder-a", Path: "./folder-a"},
		},
	}
	_, err := config.SaveConfig(&cfg, tmpDir, false)
	require.NoError(t, err)

	_, err = track.InitTrack(tmpDir, "my-feature", []string{"folder-a"}, track.InitTrackOptions{})
	require.NoError(t, err)

	buf := new(bytes.Buffer)
	err = runTrackStatus(buf, tmpDir, "my-feature")
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Multi-Barrel Track Status")
	assert.Contains(t, buf.String(), "my-feature")
}

func TestTrackListCmd(t *testing.T) {
	tmpDir := t.TempDir()

	_, err := track.InitTrack(tmpDir, "track-1", []string{"folder-a"}, track.InitTrackOptions{Name: "Track One"})
	require.NoError(t, err)

	buf := new(bytes.Buffer)
	err = runTrackList(buf, tmpDir)
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "track-1")
}
