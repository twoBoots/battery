package track

import (
	"time"
)

// TrackStatus represents the lifecycle state of a track.
type TrackStatus string

const (
	TrackStatusNotStarted TrackStatus = "not-started"
	TrackStatusPlanning   TrackStatus = "planning"
	TrackStatusInProgress TrackStatus = "in-progress"
	TrackStatusCompleted  TrackStatus = "completed"
	TrackStatusArchived   TrackStatus = "archived"
)

// TrackType classifies the intent of the track.
type TrackType string

const (
	TrackTypeFeature  TrackType = "feature"
	TrackTypeFix      TrackType = "fix"
	TrackTypeRefactor TrackType = "refactor"
	TrackTypeChore    TrackType = "chore"
)

// TrackLocation indicates where a track was located within a barrel or battery.
type TrackLocation string

const (
	LocationWorktree TrackLocation = "worktree"
	LocationActive   TrackLocation = "active"
	LocationArchive  TrackLocation = "archive"
	LocationMissing  TrackLocation = "missing"
)

// TrackMetadata represents metadata stored in metadata.json.
type TrackMetadata struct {
	TrackID      string      `json:"track_id"`
	Name         string      `json:"name"`
	Status       TrackStatus `json:"status"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	Type         TrackType   `json:"type,omitempty"`
	Barrels      []string    `json:"barrels,omitempty"`
	Capabilities []string    `json:"capabilities,omitempty"`
}

// BarrelTrackSummary represents the aggregate state of a track within a single barrel.
type BarrelTrackSummary struct {
	BarrelName      string        `json:"barrel_name"`
	BarrelPath      string        `json:"barrel_path"`
	Location        TrackLocation `json:"location"`
	Status          TrackStatus   `json:"status"`
	ActivePlanTasks int           `json:"active_plan_tasks"`
	CompletedTasks  int           `json:"completed_tasks"`
	SpecDeltas      []string      `json:"spec_deltas,omitempty"`
	Error           string        `json:"error,omitempty"`
}

// PercentComplete calculates the completion percentage (0-100).
func (b BarrelTrackSummary) PercentComplete() int {
	if b.Location == LocationArchive || b.Status == TrackStatusCompleted || b.Status == TrackStatusArchived {
		return 100
	}
	if b.ActivePlanTasks <= 0 {
		return 0
	}
	if b.CompletedTasks >= b.ActivePlanTasks {
		return 100
	}
	return int(float64(b.CompletedTasks) / float64(b.ActivePlanTasks) * 100.0)
}
