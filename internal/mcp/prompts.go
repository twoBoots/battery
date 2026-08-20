package mcp

import (
	"context"
	"fmt"
	"strings"
)

// RegisterDefaultPrompts registers prompt templates on the MCP server.
func RegisterDefaultPrompts(s *Server) {
	// 1. plan_multi_barrel_track
	s.RegisterPrompt(Prompt{
		Name:        "plan_multi_barrel_track",
		Description: "Guided prompt for AI agents to plan and dispatch a multi-barrel track adhering to Cooper Hybrid and Troop worktree protocols.",
		Arguments: []PromptArgument{
			{
				Name:        "track_id",
				Description: "Unique track identifier (e.g. track_oauth_20260817)",
				Required:    true,
			},
			{
				Name:        "goal",
				Description: "Brief description of the cross-cutting feature or change",
				Required:    false,
			},
		},
	}, func(ctx context.Context, args map[string]string) (GetPromptResult, error) {
		trackID := strings.TrimSpace(args["track_id"])
		if trackID == "" {
			return GetPromptResult{}, fmt.Errorf("track_id argument is required")
		}
		goal := strings.TrimSpace(args["goal"])
		if goal == "" {
			goal = "Implement cross-cutting capability across targeted barrels"
		}

		promptText := fmt.Sprintf(`You are an autonomous AI software engineer coordinating a multi-repository SDD track using Battery, Cooper, and Troop.

### Track Context
* **Track ID**: %s
* **Goal**: %s

### Multi-Barrel SDD Protocol
1. **Discover & Inspect Barrels**:
   - Call MCP tool 'battery_list_barrels' to verify registered barrels and read their individual Cooper tech stack definitions.
2. **Initialize Track**:
   - Call MCP tool 'battery_init_track' with track_id '%s' and targeted barrels to scaffold '.cooper/active/%s/' in the orchestrator.
3. **Author Living Spec Deltas**:
   - Define interface contracts and requirements in GIVEN/WHEN/THEN format inside '.cooper/active/%s/spec-deltas/<capability>/spec.md'.
4. **Dispatch Contracts**:
   - Call MCP tool 'battery_dispatch_track' with track_id '%s'. This seeds spec deltas into barrel worktrees under '.worktrees/%s/' while preserving local planning autonomy (omitting barrel-level plan.md).
5. **Monitor Progress**:
   - Periodically check cross-barrel status using 'battery_track_status' with track_id '%s'.
`, trackID, goal, trackID, trackID, trackID, trackID, trackID, trackID)

		return GetPromptResult{
			Description: fmt.Sprintf("Multi-Barrel Track Planning Guide for %s", trackID),
			Messages: []PromptMessage{
				{
					Role: "user",
					Content: ContentItem{
						Type: "text",
						Text: promptText,
					},
				},
			},
		}, nil
	})

	// 2. guide_framework_upgrade_track
	s.RegisterPrompt(Prompt{
		Name:        "guide_framework_upgrade_track",
		Description: "Guided prompt for AI agents to inspect Cooper/Battery framework alignment, preserve local customizations, and guide the user through adopting upgrades in an isolated track worktree.",
		Arguments: []PromptArgument{
			{
				Name:        "track_id",
				Description: "Suggested track identifier for the upgrade (e.g. track_upgrade_cooper_20260820)",
				Required:    false,
			},
			{
				Name:        "barrel",
				Description: "Optional target barrel name or path to inspect and upgrade",
				Required:    false,
			},
		},
	}, func(ctx context.Context, args map[string]string) (GetPromptResult, error) {
		trackID := strings.TrimSpace(args["track_id"])
		if trackID == "" {
			trackID = "track_upgrade_cooper_standards"
		}
		barrel := strings.TrimSpace(args["barrel"])
		barrelContext := "workspace root"
		if barrel != "" {
			barrelContext = fmt.Sprintf("barrel '%s'", barrel)
		}

		promptText := fmt.Sprintf(`You are an autonomous AI software engineer guiding the adoption of updated Cooper SDD & Battery framework standards for %s.

### Upgrade Context
* **Track ID**: %s
* **Target**: %s

### Non-Destructive Framework Upgrade Protocol
1. **Inspect Framework Status**:
   - Call MCP tool 'battery_framework_status'%s to detect version alignment and identify which files are up to date, customized locally, or missing.
2. **Retrieve Upstream Canonical Templates**:
   - Call MCP tool 'battery_get_template' for any diverged or missing skills (e.g. 'skills/cooper-rfc', 'skills/cooper-review') or documents (e.g. 'docs/COOPER.md', 'docs/BATTERY.md').
3. **Analyze Local Customizations (Crucial)**:
   - Compare local workspace files against the retrieved upstream templates.
   - Explicitly identify project-specific rules, linters, commands, and workflow modifications that MUST be preserved.
4. **Propose Dedicated Upgrade Track**:
   - Present a clear summary to the user outlining upstream improvements vs. local customizations.
   - Propose initializing track '%s'.
5. **Execute in Isolated Worktree**:
   - Call 'battery_init_track' or use 'cooper-new-track' to scaffold the track under '.worktrees/%s/'.
   - Perform a 3-way semantic merge: incorporate new upstream capabilities while preserving all team-specific customizations.
   - Run project test and lint suites, then open a Pull Request for team review.
`, barrelContext, trackID, barrelContext, func() string {
			if barrel != "" {
				return fmt.Sprintf(" with barrel: %q", barrel)
			}
			return ""
		}(), trackID, trackID)

		return GetPromptResult{
			Description: fmt.Sprintf("Framework & Standards Upgrade Guide (%s)", barrelContext),
			Messages: []PromptMessage{
				{
					Role: "user",
					Content: ContentItem{
						Type: "text",
						Text: promptText,
					},
				},
			},
		}, nil
	})
}
