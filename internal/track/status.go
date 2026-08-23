package track

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/twoBoots/battery/internal/config"
)

// PlanTaskSummary contains counts of tasks in a plan.md file.
type PlanTaskSummary struct {
	Total      int `json:"total"`
	Completed  int `json:"completed"`
	InProgress int `json:"in_progress"`
	Pending    int `json:"pending"`
}

var (
	taskDoneRegex       = regexp.MustCompile(`^\s*-\s*\[[xX]\]`)
	taskInProgressRegex = regexp.MustCompile(`^\s*-\s*\[~\]`)
	taskPendingRegex    = regexp.MustCompile(`^\s*-\s*\[\s*\]`)
)

// ParsePlanTasks parses markdown content and extracts task counts.
func ParsePlanTasks(content string) PlanTaskSummary {
	var summary PlanTaskSummary
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if taskDoneRegex.MatchString(line) {
			summary.Completed++
			summary.Total++
		} else if taskInProgressRegex.MatchString(line) {
			summary.InProgress++
			summary.Total++
		} else if taskPendingRegex.MatchString(line) {
			summary.Pending++
			summary.Total++
		}
	}
	return summary
}

// MultiBarrelTrackStatus summarizes a track across Battery and all participating barrels.
type MultiBarrelTrackStatus struct {
	TrackID      string               `json:"track_id"`
	Name         string               `json:"name"`
	Status       TrackStatus          `json:"status"`
	Barrels      []BarrelTrackSummary `json:"barrels"`
	Capabilities []string             `json:"capabilities,omitempty"`
}

// GetMultiBarrelTrackStatus inspects the track in battery and in each registered barrel.
func GetMultiBarrelTrackStatus(cwd, trackID string) (*MultiBarrelTrackStatus, error) {
	// Find track in battery (.cooper/active or .cooper/archive or worktree)
	loc, trackDir, exists := LocateBarrelTrack(cwd, trackID)
	if !exists {
		return nil, fmt.Errorf("track %q not found in battery workspace", trackID)
	}

	metaBytes, err := os.ReadFile(filepath.Join(trackDir, "metadata.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to read track metadata: %w", err)
	}

	var meta TrackMetadata
	if err := json.Unmarshal(metaBytes, &meta); err != nil {
		return nil, fmt.Errorf("failed to parse track metadata: %w", err)
	}

	effCfg, err := config.GetEffectiveConfig(cwd)
	if err != nil {
		return nil, fmt.Errorf("failed to load battery configuration: %w", err)
	}

	barrelMap := make(map[string]config.EffectiveBarrel)
	for _, b := range effCfg.Barrels {
		barrelMap[b.Name] = b
	}

	status := MultiBarrelTrackStatus{
		TrackID:      meta.TrackID,
		Name:         meta.Name,
		Status:       meta.Status,
		Capabilities: meta.Capabilities,
		Barrels:      make([]BarrelTrackSummary, 0, len(meta.Barrels)),
	}

	if loc == LocationArchive {
		status.Status = TrackStatusArchived
	}

	for _, barrelName := range meta.Barrels {
		barrel, barrelExists := barrelMap[barrelName]
		if !barrelExists {
			status.Barrels = append(status.Barrels, BarrelTrackSummary{
				BarrelName: barrelName,
				Location:   LocationMissing,
				Status:     TrackStatusNotStarted,
				Error:      "not configured in .batteryrc",
			})
			continue
		}

		bSummary := inspectBarrelTrack(cwd, barrel, trackID)
		status.Barrels = append(status.Barrels, bSummary)
	}

	return &status, nil
}

func inspectBarrelTrack(cwd string, barrel config.EffectiveBarrel, trackID string) BarrelTrackSummary {
	barrelAbs := barrel.AbsolutePath
	if barrelAbs == "" {
		if filepath.IsAbs(barrel.Path) {
			barrelAbs = barrel.Path
		} else {
			barrelAbs = filepath.Join(cwd, barrel.Path)
		}
	}

	bLoc, bTrackDir, bExists := LocateBarrelTrack(barrelAbs, trackID)
	if !bExists {
		return BarrelTrackSummary{
			BarrelName: barrel.Name,
			BarrelPath: barrel.Path,
			Location:   LocationMissing,
			Status:     TrackStatusNotStarted,
		}
	}

	summary := BarrelTrackSummary{
		BarrelName: barrel.Name,
		BarrelPath: barrel.Path,
		Location:   bLoc,
	}

	if bLoc == LocationArchive {
		summary.Status = TrackStatusCompleted
		return summary
	}

	// Check if plan.md exists
	planPath := filepath.Join(bTrackDir, "plan.md")
	if planData, err := os.ReadFile(planPath); err == nil {
		tasks := ParsePlanTasks(string(planData))
		summary.ActivePlanTasks = tasks.Total
		summary.CompletedTasks = tasks.Completed

		if tasks.Total > 0 && tasks.Completed >= tasks.Total {
			summary.Status = TrackStatusCompleted
		} else if tasks.Completed > 0 || tasks.InProgress > 0 {
			summary.Status = TrackStatusInProgress
		} else {
			summary.Status = TrackStatusPlanning
		}
	} else {
		// No plan.md yet -> barrel is in planning phase
		summary.Status = TrackStatusPlanning
	}

	// Scan spec-deltas
	specDeltasDir := filepath.Join(bTrackDir, "spec-deltas")
	if entries, err := os.ReadDir(specDeltasDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				summary.SpecDeltas = append(summary.SpecDeltas, entry.Name())
			}
		}
	}

	return summary
}

// ListTracks scans .cooper/active and .cooper/archive in cwd and returns track metadata.
func ListTracks(cwd string) ([]TrackMetadata, error) {
	trackMap := make(map[string]TrackMetadata)

	scanDirs := []struct {
		path   string
		status TrackStatus
	}{
		{filepath.Join(cwd, ".cooper", "active"), TrackStatusInProgress},
		{filepath.Join(cwd, ".cooper", "archive"), TrackStatusArchived},
	}

	for _, sd := range scanDirs {
		entries, err := os.ReadDir(sd.path)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			trackID := entry.Name()
			metaPath := filepath.Join(sd.path, trackID, "metadata.json")
			metaBytes, err := os.ReadFile(metaPath)
			if err != nil {
				// Fallback minimal metadata
				trackMap[trackID] = TrackMetadata{
					TrackID: trackID,
					Name:    trackID,
					Status:  sd.status,
				}
				continue
			}

			var meta TrackMetadata
			if err := json.Unmarshal(metaBytes, &meta); err == nil {
				if meta.Status == "" {
					meta.Status = sd.status
				}
				trackMap[trackID] = meta
			}
		}
	}

	tracks := make([]TrackMetadata, 0, len(trackMap))
	for _, t := range trackMap {
		tracks = append(tracks, t)
	}

	sort.Slice(tracks, func(i, j int) bool {
		return tracks[i].TrackID < tracks[j].TrackID
	})

	return tracks, nil
}
