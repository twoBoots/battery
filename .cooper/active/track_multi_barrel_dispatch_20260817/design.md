# Technical Design: Multi-Barrel Track Dispatch & Decoupled Planning

## Architecture Overview

The `track` package (`internal/track`) provides track lifecycle management and multi-barrel dispatch for Battery, interacting with target barrels configured via `.batteryrc` / `.batteryrc.local`.

```
cmd/
└── track.go                 # Cobra CLI commands: battery track (init, dispatch, status, list)

internal/
├── config/                  # Existing configuration & barrel resolution
├── techstack/               # Existing Cooper tech stack detector
└── track/                   # New track management engine
    ├── track.go             # Track models, state inspection, and validation
    ├── dispatch.go          # Spec delta seeding & downstream barrel initialization
    ├── status.go            # Cross-barrel progress aggregation & checkpoint inspector
    └── track_test.go        # Comprehensive unit & integration test suite (>80% coverage)
```

## Data Structures

```go
type TrackMetadata struct {
    TrackID      string    `json:"track_id"`
    Name         string    `json:"name"`
    Status       string    `json:"status"` // "in-progress", "completed", "archived"
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
    Type         string    `json:"type"`   // "feature", "fix", "refactor"
    Barrels      []string  `json:"barrels"`
    Capabilities []string  `json:"capabilities"`
}

type BarrelTrackSummary struct {
    BarrelName      string   `json:"barrel_name"`
    BarrelPath      string   `json:"barrel_path"`
    Location        string   `json:"location"` // "worktree", "active", "archive", "missing"
    Status          string   `json:"status"`   // "not-started", "planning", "in-progress", "completed"
    ActivePlanTasks int      `json:"active_plan_tasks"`
    CompletedTasks  int      `json:"completed_tasks"`
    SpecDeltas      []string `json:"spec_deltas"`
}
```

## Dispatch Engine Mechanics

1. **Track Initialization (`battery track init`)**:
   - Creates `.cooper/active/<track_id>/` in `battery` with `metadata.json`, `proposal.md`, `design.md`, and `plan.md`.
2. **Track Dispatch (`battery track dispatch`)**:
   - Iterates through target barrels registered in `.batteryrc`.
   - Checks if `.worktrees/<track_id>` or barrel `.cooper/active/<track_id>` exists.
   - Generates boilerplate `metadata.json`, `proposal.md`, and `spec-deltas/<capability>/spec.md` with GIVEN/WHEN/THEN templates.
   - Leaves `plan.md` to be created by the barrel's local agent.
3. **Status Aggregation (`battery track status`)**:
   - Inspects each target barrel across three locations in order:
     1. `.worktrees/<track_id>/.cooper/active/<track_id>/` (worktree active)
     2. `.cooper/active/<track_id>/` (trunk active)
     3. `.cooper/archive/<track_id>/` (completed & archived)
   - If found in `.cooper/archive/`, marks the barrel as `completed` / `archived` (100%).
   - If found in active folders, parses `plan.md` task checklists (`[ ]`, `[~]`, `[x]`) to calculate completion percentage and current phase.
   - Outputs a unified multi-barrel progress table.
