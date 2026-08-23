package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twoBoots/battery/internal/config"
	"github.com/twoBoots/battery/internal/track"
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

	// Default barrels when none specified
	buf.Reset()
	err = runTrackInit(buf, tmpDir, "track-default-barrels", "", "", nil, nil)
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Initialized track 'track-default-barrels'")
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

	// Dispatch again with force=true
	buf.Reset()
	err = runTrackDispatch(buf, tmpDir, "my-feature", true)
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Dispatched track 'my-feature'")
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

	// Status on non-existent track
	buf.Reset()
	err = runTrackStatus(buf, tmpDir, "non-existent")
	assert.Error(t, err)

	// Status with active tasks in barrel plan
	barrelTrackDir := filepath.Join(tmpDir, "folder-a", ".cooper", "active", "my-feature")
	require.NoError(t, os.MkdirAll(barrelTrackDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(barrelTrackDir, "metadata.json"), []byte(`{"track_id":"my-feature","status":"in-progress"}`), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(barrelTrackDir, "plan.md"), []byte("- [x] Task 1\n- [ ] Task 2\n"), 0644))

	buf.Reset()
	err = runTrackStatus(buf, tmpDir, "my-feature")
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "50% (1/2 tasks completed)")
}

func TestTrackListCmd(t *testing.T) {
	tmpDir := t.TempDir()

	buf := new(bytes.Buffer)
	// Empty list
	err := runTrackList(buf, tmpDir)
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "No active or archived tracks found")

	// With tracks
	_, err = track.InitTrack(tmpDir, "track-1", []string{"folder-a"}, track.InitTrackOptions{Name: "Track One"})
	require.NoError(t, err)

	buf.Reset()
	err = runTrackList(buf, tmpDir)
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "track-1")
}
