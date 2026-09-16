# Multi-Barrel Track Lifecycle Workflow 🔄

Battery provides a structured, contract-first orchestration lifecycle for feature epics spanning multiple repositories or monorepo packages (barrels 🛢️). By combining **[Cooper](https://twoboots.github.io/cooper)** for Spec-Driven Development (SDD) and **[Troop](https://twoboots.github.io/troop)** for Git worktree isolation, Battery ensures cross-barrel consistency without sacrificing repository autonomy.

## 🧭 The Multi-Barrel Lifecycle

<div class="lifecycle-flow">
  <div class="lifecycle-card">
    <span class="lifecycle-badge">Step 1</span>
    <div class="lifecycle-step">battery track init</div>
    <div class="lifecycle-sub">Macro Contracts &amp; Spec Deltas</div>
  </div>
  <div class="lifecycle-arrow">→</div>
  <div class="lifecycle-card">
    <span class="lifecycle-badge">Step 2</span>
    <div class="lifecycle-step">battery track dispatch</div>
    <div class="lifecycle-sub">Sync to Barrels (No plan.md)</div>
  </div>
  <div class="lifecycle-arrow">→</div>
  <div class="lifecycle-card">
    <span class="lifecycle-badge">Step 3</span>
    <div class="lifecycle-step">git agent-start</div>
    <div class="lifecycle-sub">Spawn Barrel Worktrees</div>
  </div>
  <div class="lifecycle-arrow">→</div>
  <div class="lifecycle-card">
    <span class="lifecycle-badge">Step 4</span>
    <div class="lifecycle-step">Local TDD Loop</div>
    <div class="lifecycle-sub">Red -&gt; Green -&gt; Refactor</div>
  </div>
  <div class="lifecycle-arrow">→</div>
  <div class="lifecycle-card">
    <span class="lifecycle-badge">Step 5</span>
    <div class="lifecycle-step">battery track status</div>
    <div class="lifecycle-sub">PR, Merge &amp; Teardown</div>
  </div>
</div>

## 1. Track Initialization

To initiate a cross-cutting epic across barrels:

```bash
battery track init <track_id> --barrels <barrel-1>,<barrel-2> --name "Feature Title"
```

This scaffolds a dedicated track directory at the Battery root:
* `.cooper/active/<track_id>/metadata.json`: Stores participating barrels, status, and metadata.
* `.cooper/active/<track_id>/proposal.md`: Defines high-level business goals, requirements, and scope boundaries.
* `.cooper/active/<track_id>/design.md`: Specifies shared contracts (REST endpoints, protobuf schemas, events).
* `.cooper/active/<track_id>/spec-deltas/`: Scaffolds living specification additions (`+`) and removals (`-`).

## 2. Decoupled Planning Protocol

> [!IMPORTANT]
> **Orchestrator Boundary Rule**: When operating from the Battery multi-barrel root, agents and developers author macro contracts and spec deltas. The Battery root **MUST NOT** author barrel-internal `plan.md` files or write target barrel code directly.

Local implementation planning is intentionally delegated to autonomous sessions running inside each target barrel. This guarantees:
1. **Context Window Optimization**: Barrel agents only need context on their own repository codebase.
2. **Tech Stack Autonomy**: Each barrel uses its own idioms, test runners, and styleguides defined in `.cooper/definition/tech-stack.md`.
3. **Collision Avoidance**: Barrel agents create their own TDD checklists without conflicting with cross-repo plans.

## 3. Spec Dispatch

Once contracts and spec deltas are ready, dispatch the track into participating barrels:

```bash
battery track dispatch <track_id>
```

Battery synchronizes the track artifacts and spec deltas to each barrel's `.cooper/active/<track_id>/` directory while strictly omitting `plan.md` to preserve local planning autonomy.

## 4. Worktree Isolation with [Troop](https://twoboots.github.io/troop)

Inside each target barrel, developers or agents spawn an isolated worktree based off `main`:

```bash
# Navigate to target barrel
cd ../backend

# Spawn isolated worktree off main
git agent-start <track_id>
```

This ensures feature implementation is completely isolated from the main trunk in `.worktrees/<track_id>`.

## 5. Autonomous TDD Implementation

Within the isolated barrel worktree:
1. Ground implementation in living capability specs using [Cooper](https://twoboots.github.io/cooper) SDD.
2. Author local `plan.md` with granular TDD tasks.
3. Follow the strict **Red -> Green -> Refactor** cycle:
   * **Red**: Write failing unit or integration tests.
   * **Green**: Write minimal code to pass tests.
   * **Refactor**: Optimize code, verify styleguides, ensure >80% test coverage.
4. Record Git Notes metadata on commits:
   ```bash
   git notes add -m "Task: <task_name>\nScope: <files>\nSummary: <details>" <commit_sha>
   ```

## 6. Centralized Status Aggregation

From the Battery root, track status across all participating barrels can be aggregated in real time:

```bash
# Check status of specific multi-barrel track
battery track status <track_id>

# List all active and completed tracks
battery track list
```

Battery parses each barrel's local `plan.md` and phase checkpoints to display overall multi-repo progress.

## 7. Review, Merge & Teardown

1. Submit pull requests in each target barrel repository.
2. After PR merge, cleanly tear down the isolated worktree:
   ```bash
   git agent-stop <track_id>
   ```
3. Mark the multi-barrel track complete in Battery:
   ```bash
   battery track status <track_id>
   ```

## 🔗 Related Resources

* [Getting Started Guide](getting-started.md)
* [Architecture & Topology Model](../architecture.md)
* [Model Context Protocol (MCP) Server](../mcp.md)
* [Installation Guide](../installation.md)
* [Cooper Framework](https://twoboots.github.io/cooper)
* [Troop Isolation Tool](https://twoboots.github.io/troop)

<style>
.lifecycle-flow {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin: 24px 0;
  flex-wrap: wrap;
}
.lifecycle-card {
  flex: 1 1 130px;
  min-width: 120px;
  background-color: var(--vp-c-bg-soft);
  border: 1px solid var(--vp-c-divider);
  border-radius: 8px;
  padding: 12px;
  text-align: center;
  box-sizing: border-box;
}
.lifecycle-badge {
  display: inline-block;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: var(--vp-c-brand-1);
  margin-bottom: 4px;
}
.lifecycle-step {
  font-family: var(--vp-font-family-mono);
  font-size: 13px;
  font-weight: 700;
  color: var(--vp-c-text-1);
  margin-bottom: 4px;
}
.lifecycle-sub {
  font-size: 12px;
  color: var(--vp-c-text-2);
  line-height: 1.3;
}
.lifecycle-arrow {
  font-size: 18px;
  font-weight: bold;
  color: var(--vp-c-text-3);
  user-select: none;
}
@media (max-width: 640px) {
  .lifecycle-flow {
    flex-direction: column;
    align-items: stretch;
  }
  .lifecycle-arrow {
    text-align: center;
    transform: rotate(90deg);
  }
  .lifecycle-card {
    min-width: 100%;
  }
}
</style>
