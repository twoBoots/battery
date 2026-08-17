package track

import (
	"os"
	"path/filepath"
)

// LocateBarrelTrack searches for a track inside a barrel in the following priority order:
// 1. .worktrees/<track_id>/.cooper/active/<track_id> (isolated Troop worktree active track)
// 2. .cooper/active/<track_id> (main trunk active track)
// 3. .cooper/archive/<track_id> (completed & archived track)
// Returns the location type, absolute/normalized directory path, and whether it exists.
func LocateBarrelTrack(barrelPath, trackID string) (TrackLocation, string, bool) {
	// 1. Check worktree location
	worktreePath := filepath.Join(barrelPath, ".worktrees", trackID, ".cooper", "active", trackID)
	if isDir(worktreePath) {
		return LocationWorktree, worktreePath, true
	}

	// 2. Check trunk active location
	activePath := filepath.Join(barrelPath, ".cooper", "active", trackID)
	if isDir(activePath) {
		return LocationActive, activePath, true
	}

	// 3. Check archive location
	archivePath := filepath.Join(barrelPath, ".cooper", "archive", trackID)
	if isDir(archivePath) {
		return LocationArchive, archivePath, true
	}

	return LocationMissing, "", false
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
