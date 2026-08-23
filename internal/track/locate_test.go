package track_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twoBoots/battery/internal/track"
)

func TestLocateBarrelTrack(t *testing.T) {
	tmpDir := t.TempDir()

	// Scenario 1: Worktree active track exists
	worktreeTrackDir := filepath.Join(tmpDir, "barrel-wt", ".worktrees", "my-feature", ".cooper", "active", "my-feature")
	require.NoError(t, os.MkdirAll(worktreeTrackDir, 0755))

	loc, path, exists := track.LocateBarrelTrack(filepath.Join(tmpDir, "barrel-wt"), "my-feature")
	assert.True(t, exists)
	assert.Equal(t, track.LocationWorktree, loc)
	assert.Equal(t, worktreeTrackDir, path)

	// Scenario 2: Main trunk active track exists (when worktree does not)
	activeTrackDir := filepath.Join(tmpDir, "barrel-active", ".cooper", "active", "my-feature")
	require.NoError(t, os.MkdirAll(activeTrackDir, 0755))

	loc, path, exists = track.LocateBarrelTrack(filepath.Join(tmpDir, "barrel-active"), "my-feature")
	assert.True(t, exists)
	assert.Equal(t, track.LocationActive, loc)
	assert.Equal(t, activeTrackDir, path)

	// Scenario 3: Completed & archived track exists in .cooper/archive/
	archiveTrackDir := filepath.Join(tmpDir, "barrel-archived", ".cooper", "archive", "my-feature")
	require.NoError(t, os.MkdirAll(archiveTrackDir, 0755))

	loc, path, exists = track.LocateBarrelTrack(filepath.Join(tmpDir, "barrel-archived"), "my-feature")
	assert.True(t, exists)
	assert.Equal(t, track.LocationArchive, loc)
	assert.Equal(t, archiveTrackDir, path)

	// Scenario 4: Missing track
	missingBarrelDir := filepath.Join(tmpDir, "barrel-missing")
	require.NoError(t, os.MkdirAll(missingBarrelDir, 0755))

	loc, path, exists = track.LocateBarrelTrack(missingBarrelDir, "my-feature")
	assert.False(t, exists)
	assert.Equal(t, track.LocationMissing, loc)
	assert.Empty(t, path)
}
