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
}
