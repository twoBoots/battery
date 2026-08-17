package track_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twoboots/battery/internal/config"
	"github.com/twoboots/battery/internal/track"
)

func TestInitTrack_Success(t *testing.T) {
	tmpDir := t.TempDir()

	opts := track.InitTrackOptions{
		Name:         "User Profile V2",
		Type:         track.TrackTypeFeature,
		Capabilities: []string{"user-profile", "avatar-upload"},
	}

	meta, err := track.InitTrack(tmpDir, "user-profile-v2", []string{"folder-a", "folder-b"}, opts)
	require.NoError(t, err)
	require.NotNil(t, meta)

	assert.Equal(t, "user-profile-v2", meta.TrackID)
	assert.Equal(t, "User Profile V2", meta.Name)
	assert.Equal(t, track.TrackStatusInProgress, meta.Status)
	assert.Equal(t, track.TrackTypeFeature, meta.Type)
	assert.Equal(t, []string{"folder-a", "folder-b"}, meta.Barrels)
	assert.Equal(t, []string{"user-profile", "avatar-upload"}, meta.Capabilities)

	// Check files created
	trackDir := filepath.Join(tmpDir, ".cooper", "active", "user-profile-v2")
	assert.FileExists(t, filepath.Join(trackDir, "metadata.json"))
	assert.FileExists(t, filepath.Join(trackDir, "proposal.md"))
	assert.FileExists(t, filepath.Join(trackDir, "design.md"))
	assert.FileExists(t, filepath.Join(trackDir, "plan.md"))

	// Check metadata content
	data, err := os.ReadFile(filepath.Join(trackDir, "metadata.json"))
	require.NoError(t, err)
	var loadedMeta track.TrackMetadata
	require.NoError(t, json.Unmarshal(data, &loadedMeta))
	assert.Equal(t, "user-profile-v2", loadedMeta.TrackID)
}

func TestInitTrack_InvalidInputs(t *testing.T) {
	tmpDir := t.TempDir()

	// Empty track ID
	_, err := track.InitTrack(tmpDir, "", []string{"folder-a"}, track.InitTrackOptions{})
	assert.Error(t, err)

	// Already existing track without overwrite
	_, err = track.InitTrack(tmpDir, "existing-track", []string{"folder-a"}, track.InitTrackOptions{})
	require.NoError(t, err)

	_, err = track.InitTrack(tmpDir, "existing-track", []string{"folder-a"}, track.InitTrackOptions{})
	assert.Error(t, err)
}

func TestDispatchTrack_SeedsBarrelsWithoutPlan(t *testing.T) {
	tmpDir := t.TempDir()

	// 1. Setup battery workspace with .batteryrc and 2 barrels
	barrelADir := filepath.Join(tmpDir, "folder-a")
	barrelBDir := filepath.Join(tmpDir, "folder-b")
	require.NoError(t, os.MkdirAll(barrelADir, 0755))
	require.NoError(t, os.MkdirAll(barrelBDir, 0755))

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

	// 2. Initialize a track in battery
	_, err = track.InitTrack(tmpDir, "user-profile-v2", []string{"folder-a", "folder-b"}, track.InitTrackOptions{
		Name:         "User Profile V2",
		Capabilities: []string{"user-profile"},
	})
	require.NoError(t, err)

	// 3. Dispatch track to barrels
	results, err := track.DispatchTrack(tmpDir, "user-profile-v2", track.DispatchTrackOptions{})
	require.NoError(t, err)
	assert.Len(t, results, 2)

	// 4. Verify folder-a created artifacts
	barrelATrackDir := filepath.Join(barrelADir, ".cooper", "active", "user-profile-v2")
	assert.FileExists(t, filepath.Join(barrelATrackDir, "metadata.json"))
	assert.FileExists(t, filepath.Join(barrelATrackDir, "proposal.md"))
	assert.FileExists(t, filepath.Join(barrelATrackDir, "spec-deltas", "user-profile", "spec.md"))

	// CRITICAL REQUIREMENT: plan.md MUST NOT be created in target barrel
	assert.NoFileExists(t, filepath.Join(barrelATrackDir, "plan.md"), "plan.md must be omitted in target barrel to allow local agent planning")

	// 5. Verify folder-b created artifacts
	barrelBTrackDir := filepath.Join(barrelBDir, ".cooper", "active", "user-profile-v2")
	assert.FileExists(t, filepath.Join(barrelBTrackDir, "metadata.json"))
	assert.FileExists(t, filepath.Join(barrelBTrackDir, "proposal.md"))
	assert.FileExists(t, filepath.Join(barrelBTrackDir, "spec-deltas", "user-profile", "spec.md"))
	assert.NoFileExists(t, filepath.Join(barrelBTrackDir, "plan.md"))
}

func TestDispatchTrack_TargetsWorktreeWhenPresent(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup barrel with an existing Troop worktree
	barrelDir := filepath.Join(tmpDir, "folder-a")
	worktreeRoot := filepath.Join(barrelDir, ".worktrees", "user-profile-v2")
	require.NoError(t, os.MkdirAll(worktreeRoot, 0755))

	cfg := config.BatteryConfig{
		Version:   "1.0.0",
		Structure: config.StructureMultiRepo,
		Barrels: []config.BarrelConfig{
			{Name: "folder-a", Path: "./folder-a"},
		},
	}
	_, err := config.SaveConfig(&cfg, tmpDir, false)
	require.NoError(t, err)

	_, err = track.InitTrack(tmpDir, "user-profile-v2", []string{"folder-a"}, track.InitTrackOptions{
		Name:         "User Profile V2",
		Capabilities: []string{"user-profile"},
	})
	require.NoError(t, err)

	results, err := track.DispatchTrack(tmpDir, "user-profile-v2", track.DispatchTrackOptions{})
	require.NoError(t, err)
	require.Len(t, results, 1)

	// Track files should be written inside .worktrees/user-profile-v2/.cooper/active/user-profile-v2
	worktreeTrackDir := filepath.Join(worktreeRoot, ".cooper", "active", "user-profile-v2")
	assert.FileExists(t, filepath.Join(worktreeTrackDir, "metadata.json"))
	assert.FileExists(t, filepath.Join(worktreeTrackDir, "proposal.md"))
	assert.NoFileExists(t, filepath.Join(worktreeTrackDir, "plan.md"))
}
