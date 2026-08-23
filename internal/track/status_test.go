package track_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twoBoots/battery/internal/config"
	"github.com/twoBoots/battery/internal/track"
)

func TestParsePlanTasks(t *testing.T) {
	planContent := `# Execution Plan

## Phase 1
- [x] Task 1.1: Complete task A (commit123)
- [X] Task 1.2: Complete task B (commit456)
- [~] Task 1.3: In progress task C
- [ ] Task 1.4: Pending task D
- [ ] Task 1.5: Pending task E

Some notes here.
`

	summary := track.ParsePlanTasks(planContent)
	assert.Equal(t, 5, summary.Total)
	assert.Equal(t, 2, summary.Completed)
	assert.Equal(t, 1, summary.InProgress)
	assert.Equal(t, 2, summary.Pending)
}

func TestGetMultiBarrelTrackStatus_WithDiverseBarrelStates(t *testing.T) {
	tmpDir := t.TempDir()

	// 1. Barrels:
	// - folder-a: Has active worktree with plan.md (2/4 tasks done)
	// - folder-b: Completed and archived in .cooper/archive/
	// - folder-c: Dispatched with specs only (no plan.md -> planning)
	// - folder-d: Missing entirely
	barrelADir := filepath.Join(tmpDir, "folder-a")
	barrelBDir := filepath.Join(tmpDir, "folder-b")
	barrelCDir := filepath.Join(tmpDir, "folder-c")
	barrelDDir := filepath.Join(tmpDir, "folder-d")

	require.NoError(t, os.MkdirAll(barrelADir, 0755))
	require.NoError(t, os.MkdirAll(barrelBDir, 0755))
	require.NoError(t, os.MkdirAll(barrelCDir, 0755))
	require.NoError(t, os.MkdirAll(barrelDDir, 0755))

	// Setup folder-a worktree track with plan.md
	aTrackDir := filepath.Join(barrelADir, ".worktrees", "auth-v2", ".cooper", "active", "auth-v2")
	require.NoError(t, os.MkdirAll(aTrackDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(aTrackDir, "plan.md"), []byte("- [x] T1\n- [x] T2\n- [ ] T3\n- [ ] T4\n"), 0644))

	// Setup folder-b archived track
	bArchiveDir := filepath.Join(barrelBDir, ".cooper", "archive", "auth-v2")
	require.NoError(t, os.MkdirAll(bArchiveDir, 0755))

	// Setup folder-c active track without plan.md (planning phase)
	cTrackDir := filepath.Join(barrelCDir, ".cooper", "active", "auth-v2")
	require.NoError(t, os.MkdirAll(cTrackDir, 0755))

	// Setup battery config
	cfg := config.BatteryConfig{
		Version:   "1.0.0",
		Structure: config.StructureMultiRepo,
		Barrels: []config.BarrelConfig{
			{Name: "folder-a", Path: "./folder-a"},
			{Name: "folder-b", Path: "./folder-b"},
			{Name: "folder-c", Path: "./folder-c"},
			{Name: "folder-d", Path: "./folder-d"},
		},
	}
	_, err := config.SaveConfig(&cfg, tmpDir, false)
	require.NoError(t, err)

	// Init track in battery
	_, err = track.InitTrack(tmpDir, "auth-v2", []string{"folder-a", "folder-b", "folder-c", "folder-d"}, track.InitTrackOptions{
		Name: "Authentication V2",
	})
	require.NoError(t, err)

	// Get Multi-Barrel Track Status
	status, err := track.GetMultiBarrelTrackStatus(tmpDir, "auth-v2")
	require.NoError(t, err)
	require.NotNil(t, status)

	assert.Equal(t, "auth-v2", status.TrackID)
	assert.Equal(t, "Authentication V2", status.Name)
	assert.Len(t, status.Barrels, 4)

	// Check folder-a
	assert.Equal(t, "folder-a", status.Barrels[0].BarrelName)
	assert.Equal(t, track.LocationWorktree, status.Barrels[0].Location)
	assert.Equal(t, track.TrackStatusInProgress, status.Barrels[0].Status)
	assert.Equal(t, 4, status.Barrels[0].ActivePlanTasks)
	assert.Equal(t, 2, status.Barrels[0].CompletedTasks)
	assert.Equal(t, 50, status.Barrels[0].PercentComplete())

	// Check folder-b (archived)
	assert.Equal(t, "folder-b", status.Barrels[1].BarrelName)
	assert.Equal(t, track.LocationArchive, status.Barrels[1].Location)
	assert.Equal(t, track.TrackStatusCompleted, status.Barrels[1].Status)
	assert.Equal(t, 100, status.Barrels[1].PercentComplete())

	// Check folder-c (planning)
	assert.Equal(t, "folder-c", status.Barrels[2].BarrelName)
	assert.Equal(t, track.LocationActive, status.Barrels[2].Location)
	assert.Equal(t, track.TrackStatusPlanning, status.Barrels[2].Status)

	// Check folder-d (missing)
	assert.Equal(t, "folder-d", status.Barrels[3].BarrelName)
	assert.Equal(t, track.LocationMissing, status.Barrels[3].Location)
	assert.Equal(t, track.TrackStatusNotStarted, status.Barrels[3].Status)
}

func TestListTracks(t *testing.T) {
	tmpDir := t.TempDir()

	// Init 2 tracks
	_, err := track.InitTrack(tmpDir, "track-1", []string{"folder-a"}, track.InitTrackOptions{Name: "Track 1"})
	require.NoError(t, err)

	_, err = track.InitTrack(tmpDir, "track-2", []string{"folder-b"}, track.InitTrackOptions{Name: "Track 2"})
	require.NoError(t, err)

	tracks, err := track.ListTracks(tmpDir)
	require.NoError(t, err)
	assert.Len(t, tracks, 2)
}
