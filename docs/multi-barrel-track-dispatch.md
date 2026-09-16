# Multi-Barrel Track Dispatch & Decoupled Planning Architecture

## Overview

**Battery** orchestrates cross-cutting feature epics across multiple independent repositories or monorepo packages (barrels 🛢️). This document describes the **Contract-First Decoupled Planning Architecture**, establishing why Battery authors macro interface contracts and requirements while delegating localized technical design and TDD task breakdown (`plan.md`) to autonomous agents inside each target barrel.

## 🏛️ Core Principles

### 1. Separation of "What" vs. "How"
* **Battery (The Orchestrator) defines the *WHAT* & *WHY*:** System-level business intent, acceptance criteria, cross-barrel interface contracts (OpenAPI/REST, gRPC, event schemas), and living Spec Deltas (`spec-deltas/`).
* **The Barrel Agent defines the *HOW*:** Local technical architecture (`design.md`), internal component/handler structure, database queries/migrations, test doubles, and granular TDD task checklists (`plan.md`).

### 2. Context Window Optimization
AI agents working inside a specific barrel (e.g. `folder-a` backend or `folder-b` frontend) only need deep context on their own repository's code, dependencies, and test suite. Keeping planning localized prevents context pollution, hallucinations about foreign file paths, and cross-repo coordination bottlenecks.

### 3. Tech Stack Autonomy
Each barrel defines its own language, frameworks, and linting rules in its local `.cooper/definition/tech-stack.md` and `.cooper/code_styleguides/`. Local agents plan and execute according to their repository's idioms (e.g. Go table tests vs. Vitest/MSW).

## 🗂️ Artifact Hierarchy Across Barrels

When a cross-cutting feature (e.g. `track_id = "user-profile-v2"`) is orchestrated across sibling barrels `folder-a` and `folder-b`:

<div class="hierarchy-flow">
  <div class="hierarchy-card orchestrator">
    <span class="hierarchy-badge">🔋 Battery Orchestrator</span>
    <div class="hierarchy-title">.cooper/active/user-profile-v2/</div>
    <div class="hierarchy-items">
      <code>metadata.json</code> (target barrels: a &amp; b) • <code>proposal.md</code> (macro epic intent) • <code>design.md</code> (shared contracts) • <code>plan.md</code> (macro milestones)
    </div>
  </div>
  <div class="hierarchy-arrow">↓ Dispatches Track Spec Deltas (Omitting plan.md) ↓</div>
  <div class="hierarchy-grid">
    <div class="hierarchy-card barrel">
      <span class="hierarchy-subbadge">🛢️ folder-a (Backend)</span>
      <div class="hierarchy-title">.worktrees/user-profile-v2/</div>
      <div class="hierarchy-items">
        <code>proposal.md</code> • <code>design.md</code> • <code>spec-deltas/</code> • <strong><code>plan.md</code> (TDD)</strong>
      </div>
    </div>
    <div class="hierarchy-card barrel">
      <span class="hierarchy-subbadge">🛢️ folder-b (Frontend)</span>
      <div class="hierarchy-title">.worktrees/user-profile-v2/</div>
      <div class="hierarchy-items">
        <code>proposal.md</code> • <code>design.md</code> • <code>spec-deltas/</code> • <strong><code>plan.md</code> (TDD)</strong>
      </div>
    </div>
  </div>
</div>

### 1. Macro Orchestrator Level (`battery`)

Stored in `battery/.cooper/active/<track_id>/`:
* **`metadata.json`**: Track ID, status, timestamps, and list of participating barrels (`["folder-a", "folder-b"]`).
* **`proposal.md`**: Macro epic summary explaining business rationale, scope, and cross-barrel boundaries.
* **`design.md`**: Interface contract specifications (HTTP routes, payloads, gRPC protos, async event payloads).
* **`plan.md`**: Cross-repo milestone roadmap (e.g. *Gate 1: folder-a publishes API schema -> Gate 2: folder-b integrates client*).

### 2. Target Barrel Level (`folder-a` & `folder-b`)

Created inside `.worktrees/<track_id>/.cooper/active/<track_id>/` via `git agent-start <track_id>`:
* **`metadata.json`**: Barrel-level track metadata referencing the upstream Battery track ID.
* **`proposal.md`**: Localized summary of this barrel's specific contribution.
* **`design.md`**: Barrel-specific technical architecture tailored to its tech stack.
* **`spec-deltas/<capability>/spec.md`**: Requirement diffs using `+` (additions) and `-` (removals) in GIVEN/WHEN/THEN format against the barrel's living specs.
* **`plan.md`**: Step-by-step TDD task checklist with Red-Green-Refactor phases, test coverage gates (>80%), and phase checkpoints.

## 🔄 End-to-End Track Lifecycle

<div class="lifecycle-flow">
  <div class="lifecycle-card">
    <span class="lifecycle-badge">Step 1</span>
    <div class="lifecycle-step">Macro Proposal</div>
    <div class="lifecycle-sub">Shared Contract &amp; Spec Deltas</div>
  </div>
  <div class="lifecycle-arrow">→</div>
  <div class="lifecycle-card">
    <span class="lifecycle-badge">Step 2</span>
    <div class="lifecycle-step">Spec Dispatch</div>
    <div class="lifecycle-sub">Sync to Barrels (No plan.md)</div>
  </div>
  <div class="lifecycle-arrow">→</div>
  <div class="lifecycle-card">
    <span class="lifecycle-badge">Step 3</span>
    <div class="lifecycle-step">Barrel Agent-Start</div>
    <div class="lifecycle-sub">Spawn Isolated Worktrees</div>
  </div>
  <div class="lifecycle-arrow">→</div>
  <div class="lifecycle-card">
    <span class="lifecycle-badge">Step 4</span>
    <div class="lifecycle-step">TDD Execution</div>
    <div class="lifecycle-sub">Red -&gt; Green -&gt; Refactor</div>
  </div>
  <div class="lifecycle-arrow">→</div>
  <div class="lifecycle-card">
    <span class="lifecycle-badge">Step 5</span>
    <div class="lifecycle-step">Status &amp; Teardown</div>
    <div class="lifecycle-sub">PR Merge &amp; git agent-stop</div>
  </div>
</div>

### Step 1: Macro Epic Proposal in Battery
Battery defines the track and writes the shared contract specification.

### Step 2: Spec Dispatch to Target Barrels
Battery writes the initial track metadata and `spec-deltas/` into target barrels.

### Step 3: Localized Planning by Barrel Agents
Agents inside each barrel run `git agent-start <track_id>` using [Troop](https://twoboots.github.io/troop), inspect their local codebase and living capability specs via [Cooper](https://twoboots.github.io/cooper), and construct their own executable `plan.md`.

### Step 4: TDD Execution & Checkpoints
Barrels execute their local TDD loops independently. At phase completion:
1. Run local test suite (`CI=true`).
2. Commit with Git Notes (`git notes add -m "<summary>" <hash>`).
3. Push checkpoint branch (`git push origin <track_id>`).

### Step 5: Merge, Living Spec Integration & Teardown
When each barrel PR merges:
1. Spec Deltas are merged into `.cooper/specs/<capability>/spec.md`.
2. Active track moves to `.cooper/archive/<track_id>/`.
3. Worktree canopy is cleaned up with `git agent-stop <track_id>`.

## 💻 CLI Command Vision

```bash
# Initialize and dispatch a new cross-barrel track spec
battery track init <track_id> --barrels folder-a,folder-b

# Inspect multi-barrel progress across all participating worktrees
battery track status [<track_id>]

# Validate contract alignment between barrel spec deltas
battery track verify <track_id>
```

<style>
.hierarchy-flow {
  background-color: var(--vp-c-bg-soft);
  border: 1px solid var(--vp-c-divider);
  border-radius: 12px;
  padding: 16px;
  margin: 20px 0;
}
.hierarchy-badge, .hierarchy-subbadge {
  display: inline-block;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: var(--vp-c-brand-1);
  margin-bottom: 4px;
}
.hierarchy-card {
  background-color: var(--vp-c-bg);
  border: 1px solid var(--vp-c-divider);
  border-radius: 8px;
  padding: 12px;
  box-sizing: border-box;
}
.hierarchy-card.orchestrator {
  text-align: center;
  max-width: 540px;
  margin: 0 auto;
}
.hierarchy-title {
  font-family: var(--vp-font-family-mono);
  font-size: 13px;
  font-weight: 700;
  color: var(--vp-c-text-1);
  margin-bottom: 4px;
}
.hierarchy-items {
  font-size: 12px;
  color: var(--vp-c-text-2);
  line-height: 1.4;
}
.hierarchy-arrow {
  text-align: center;
  font-size: 12px;
  font-weight: 600;
  color: var(--vp-c-text-3);
  margin: 12px 0;
  user-select: none;
}
.hierarchy-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 12px;
}
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
