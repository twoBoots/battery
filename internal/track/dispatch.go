package track

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/twoboots/battery/internal/config"
)

// DispatchTrackOptions configures track dispatching.
type DispatchTrackOptions struct {
	Force         bool
	BarrelCapMap  map[string][]string // Optional mapping of barrel -> specific capabilities
}

// DispatchedBarrelResult captures the outcome of dispatching to a barrel.
type DispatchedBarrelResult struct {
	BarrelName string `json:"barrel_name"`
	BarrelPath string `json:"barrel_path"`
	TargetDir  string `json:"target_dir"`
	Created    bool   `json:"created"`
	Error      string `json:"error,omitempty"`
}

// DispatchTrack seeds target barrels with track metadata, proposal, and spec deltas.
// IMPORTANT: It deliberately DOES NOT create plan.md, empowering local barrel agents to plan autonomously.
func DispatchTrack(cwd, trackID string, opts DispatchTrackOptions) ([]DispatchedBarrelResult, error) {
	// 1. Read battery track metadata
	batteryTrackDir := filepath.Join(cwd, ".cooper", "active", trackID)
	metaBytes, err := os.ReadFile(filepath.Join(batteryTrackDir, "metadata.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to read battery track metadata for %q: %w", trackID, err)
	}

	var meta TrackMetadata
	if err := json.Unmarshal(metaBytes, &meta); err != nil {
		return nil, fmt.Errorf("failed to parse battery track metadata: %w", err)
	}

	// 2. Load effective barrels from .batteryrc
	effCfg, err := config.GetEffectiveConfig(cwd)
	if err != nil {
		return nil, fmt.Errorf("failed to load battery configuration: %w", err)
	}

	barrelMap := make(map[string]config.EffectiveBarrel)
	for _, b := range effCfg.Barrels {
		barrelMap[b.Name] = b
	}

	var results []DispatchedBarrelResult

	for _, barrelName := range meta.Barrels {
		barrel, exists := barrelMap[barrelName]
		if !exists {
			results = append(results, DispatchedBarrelResult{
				BarrelName: barrelName,
				Created:    false,
				Error:      fmt.Sprintf("barrel %q not found in .batteryrc configuration", barrelName),
			})
			continue
		}

		targetTrackDir, err := seedBarrelTrack(cwd, barrel, meta, opts)
		if err != nil {
			results = append(results, DispatchedBarrelResult{
				BarrelName: barrelName,
				BarrelPath: barrel.Path,
				Created:    false,
				Error:      err.Error(),
			})
			continue
		}

		results = append(results, DispatchedBarrelResult{
			BarrelName: barrelName,
			BarrelPath: barrel.Path,
			TargetDir:  targetTrackDir,
			Created:    true,
		})
	}

	return results, nil
}

func seedBarrelTrack(cwd string, barrel config.EffectiveBarrel, meta TrackMetadata, opts DispatchTrackOptions) (string, error) {
	barrelAbs := barrel.AbsolutePath
	if barrelAbs == "" {
		if filepath.IsAbs(barrel.Path) {
			barrelAbs = barrel.Path
		} else {
			barrelAbs = filepath.Join(cwd, barrel.Path)
		}
	}

	// Check if barrel worktree exists (.worktrees/<track_id>)
	worktreeRoot := filepath.Join(barrelAbs, ".worktrees", meta.TrackID)
	var targetTrackDir string
	if isDir(worktreeRoot) {
		targetTrackDir = filepath.Join(worktreeRoot, ".cooper", "active", meta.TrackID)
	} else {
		targetTrackDir = filepath.Join(barrelAbs, ".cooper", "active", meta.TrackID)
	}

	if _, err := os.Stat(targetTrackDir); err == nil && !opts.Force {
		// Already exists
		return targetTrackDir, nil
	}

	if err := os.MkdirAll(targetTrackDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create barrel track directory: %w", err)
	}

	// 1. Write barrel metadata.json
	barrelMeta := TrackMetadata{
		TrackID:      meta.TrackID,
		Name:         meta.Name,
		Status:       TrackStatusPlanning,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
		Type:         meta.Type,
		Capabilities: meta.Capabilities,
	}

	bMetaBytes, err := json.MarshalIndent(barrelMeta, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal barrel metadata: %w", err)
	}
	if err := os.WriteFile(filepath.Join(targetTrackDir, "metadata.json"), bMetaBytes, 0644); err != nil {
		return "", fmt.Errorf("failed to write barrel metadata.json: %w", err)
	}

	// 2. Write barrel proposal.md
	barrelProposal := fmt.Sprintf("# Proposal: %s (%s)\n\n## Intent\nImplement barrel-specific requirements for %s dispatched from battery.\n\n## Next Steps\nRun local planning to generate `design.md` and `plan.md` adhering to this repository's tech stack.\n", meta.Name, barrel.Name, meta.TrackID)
	if err := os.WriteFile(filepath.Join(targetTrackDir, "proposal.md"), []byte(barrelProposal), 0644); err != nil {
		return "", fmt.Errorf("failed to write barrel proposal.md: %w", err)
	}

	// 3. Seed spec-deltas/
	caps := meta.Capabilities
	if customCaps, ok := opts.BarrelCapMap[barrel.Name]; ok && len(customCaps) > 0 {
		caps = customCaps
	}
	if len(caps) == 0 {
		caps = []string{meta.TrackID}
	}

	for _, capName := range caps {
		specDeltaDir := filepath.Join(targetTrackDir, "spec-deltas", capName)
		if err := os.MkdirAll(specDeltaDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create spec-deltas dir for %s: %w", capName, err)
		}

		specDeltaContent := fmt.Sprintf("# Spec Delta: %s\n\n## Capability: `%s`\n\n### Requirement 1: %s\n+ Description of requirement additions and changes.\n+\n+ #### Scenario 1.1: Core Behavior\n+ - **GIVEN** a precondition\n+ - **WHEN** an action occurs\n+ - **THEN** expected outcome occurs.\n", meta.Name, capName, capName)
		if err := os.WriteFile(filepath.Join(specDeltaDir, "spec.md"), []byte(specDeltaContent), 0644); err != nil {
			return "", fmt.Errorf("failed to write spec delta for %s: %w", capName, err)
		}
	}

	// NOTICE: plan.md IS INTENTIONALLY OMITTED. Local barrel agents must create their own plan.md.

	return targetTrackDir, nil
}
