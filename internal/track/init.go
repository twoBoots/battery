package track

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var validTrackIDRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// InitTrackOptions configures track creation.
type InitTrackOptions struct {
	Name         string
	Type         TrackType
	Capabilities []string
	Proposal     string
	Design       string
	Force        bool
}

// InitTrack initializes an active track inside battery's .cooper/active/<track_id>/ directory.
func InitTrack(cwd, trackID string, barrels []string, opts InitTrackOptions) (*TrackMetadata, error) {
	trackID = strings.TrimSpace(trackID)
	if trackID == "" {
		return nil, errors.New("track_id cannot be empty")
	}
	if !validTrackIDRegex.MatchString(trackID) {
		return nil, fmt.Errorf("invalid track_id %q: must contain only letters, numbers, hyphens, and underscores", trackID)
	}

	trackDir := filepath.Join(cwd, ".cooper", "active", trackID)
	if _, err := os.Stat(trackDir); err == nil && !opts.Force {
		return nil, fmt.Errorf("track %q already exists in .cooper/active/%s", trackID, trackID)
	}

	if err := os.MkdirAll(trackDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create track directory: %w", err)
	}

	trackName := opts.Name
	if trackName == "" {
		trackName = trackID
	}
	trackType := opts.Type
	if trackType == "" {
		trackType = TrackTypeFeature
	}

	now := time.Now().UTC()
	meta := TrackMetadata{
		TrackID:      trackID,
		Name:         trackName,
		Status:       TrackStatusInProgress,
		CreatedAt:    now,
		UpdatedAt:    now,
		Type:         trackType,
		Barrels:      barrels,
		Capabilities: opts.Capabilities,
	}

	metaBytes, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %w", err)
	}
	if err := os.WriteFile(filepath.Join(trackDir, "metadata.json"), metaBytes, 0644); err != nil {
		return nil, fmt.Errorf("failed to write metadata.json: %w", err)
	}

	proposalContent := opts.Proposal
	if proposalContent == "" {
		proposalContent = fmt.Sprintf("# Track Proposal: %s\n\n## Context & Motivation\nDescribe the cross-cutting problem and business intent.\n\n## Participating Barrels\n%s\n", trackName, formatBarrelsList(barrels))
	}
	if err := os.WriteFile(filepath.Join(trackDir, "proposal.md"), []byte(proposalContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write proposal.md: %w", err)
	}

	designContent := opts.Design
	if designContent == "" {
		designContent = fmt.Sprintf("# Technical Design: %s\n\n## Architecture & Shared Contracts\nDefine shared APIs, event schemas, and interface boundaries between barrels.\n", trackName)
	}
	if err := os.WriteFile(filepath.Join(trackDir, "design.md"), []byte(designContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write design.md: %w", err)
	}

	planContent := fmt.Sprintf("# Execution Plan: %s (Battery Master Plan)\n\n## Macro Cross-Barrel Milestones\n- [ ] Phase 1: Dispatch spec deltas and seed target barrels\n- [ ] Phase 2: Downstream barrel execution and TDD checkpoints\n- [ ] Phase 3: Integration verification\n", trackName)
	if err := os.WriteFile(filepath.Join(trackDir, "plan.md"), []byte(planContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write plan.md: %w", err)
	}

	return &meta, nil
}

func formatBarrelsList(barrels []string) string {
	if len(barrels) == 0 {
		return "- None specified"
	}
	var sb strings.Builder
	for _, b := range barrels {
		sb.WriteString(fmt.Sprintf("- %s\n", b))
	}
	return sb.String()
}
