package track_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twoboots/battery/internal/track"
)

func TestTrackMetadata_JSONSerialization(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	meta := track.TrackMetadata{
		TrackID:      "test-track",
		Name:         "Test Track Name",
		Status:       track.TrackStatusInProgress,
		CreatedAt:    now,
		UpdatedAt:    now,
		Type:         track.TrackTypeFeature,
		Barrels:      []string{"barrel-a", "barrel-b"},
		Capabilities: []string{"auth", "billing"},
	}

	data, err := json.Marshal(meta)
	require.NoError(t, err)

	var decoded track.TrackMetadata
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, meta.TrackID, decoded.TrackID)
	assert.Equal(t, meta.Name, decoded.Name)
	assert.Equal(t, meta.Status, decoded.Status)
	assert.Equal(t, meta.Type, decoded.Type)
	assert.Equal(t, meta.Barrels, decoded.Barrels)
	assert.Equal(t, meta.Capabilities, decoded.Capabilities)
	assert.True(t, meta.CreatedAt.Equal(decoded.CreatedAt))
}

func TestBarrelTrackSummary_Fields(t *testing.T) {
	summary := track.BarrelTrackSummary{
		BarrelName:      "folder-a",
		BarrelPath:      "../folder-a",
		Location:        track.LocationWorktree,
		Status:          track.TrackStatusInProgress,
		ActivePlanTasks: 5,
		CompletedTasks:  3,
		SpecDeltas:      []string{"api/spec.md"},
	}

	assert.Equal(t, "folder-a", summary.BarrelName)
	assert.Equal(t, track.LocationWorktree, summary.Location)
	assert.Equal(t, track.TrackStatusInProgress, summary.Status)
	assert.Equal(t, 60, summary.PercentComplete())
}

func TestBarrelTrackSummary_PercentComplete(t *testing.T) {
	tests := []struct {
		name      string
		summary   track.BarrelTrackSummary
		expected  int
	}{
		{
			name: "zero tasks and archived",
			summary: track.BarrelTrackSummary{
				Location: track.LocationArchive,
				Status:   track.TrackStatusCompleted,
			},
			expected: 100,
		},
		{
			name: "zero tasks not archived",
			summary: track.BarrelTrackSummary{
				Location: track.LocationActive,
				Status:   track.TrackStatusPlanning,
			},
			expected: 0,
		},
		{
			name: "partial completion",
			summary: track.BarrelTrackSummary{
				ActivePlanTasks: 10,
				CompletedTasks:  4,
			},
			expected: 40,
		},
		{
			name: "all completed",
			summary: track.BarrelTrackSummary{
				ActivePlanTasks: 8,
				CompletedTasks:  8,
			},
			expected: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.summary.PercentComplete())
		})
	}
}
