---
layout: home

hero:
  name: Battery 🔋
  text: Multi-Repository SDD Orchestration Protocol
  tagline: A Collection of Barrels. Coordinate living capability specs, multi-repo tracks, and isolated worktrees across distributed codebases.
  actions:
    - theme: brand
      text: Get Started
      link: /guide/getting-started
    - theme: alt
      text: View on GitHub
      link: https://github.com/twoBoots/battery

features:
  - icon: 🛢️
    title: Multi-Barrel Coordination
    details: Orchestrate cross-cutting feature epics across independent repositories or monorepo packages (barrels) with unified lifecycle tracking.
  - icon: 📐
    title: Cooper SDD Integration
    details: Built upon <a href="https://twoboots.github.io/cooper" target="_blank" rel="noopener">Cooper</a> Hybrid Spec-Driven Development, keeping living specs and capability requirements in sync across all barrels.
  - icon: 🐒
    title: Troop Worktree Isolation
    details: Dispatch track work directly into isolated Git worktrees powered by <a href="https://twoboots.github.io/troop" target="_blank" rel="noopener">Troop</a> (<code>.worktrees/&lt;track_id&gt;</code>), protecting your trunk branches.
  - icon: 🤖
    title: Native MCP Server
    details: Built-in Model Context Protocol (MCP) server exposing tools, resources, and prompt templates directly to modern AI coding assistants.
  - icon: ⚙️
    title: Layered Configuration
    details: Canonical team topologies in <code>.batteryrc</code> with personal local overrides in <code>.batteryrc.local</code> for flexible developer environments.
  - icon: 📜
    title: Contract-First Decoupled Planning
    details: Author macro interface contracts and living spec deltas at the root while delegating localized TDD implementation plans to autonomous barrel agents.
---

## Quickstart

Install Battery into your workspace with a single command:

```bash
curl -fsSL https://raw.githubusercontent.com/twoBoots/battery/main/install.sh | bash
```

The installer verifies your environment, sets up [Troop](https://twoboots.github.io/troop) Git aliases, scaffolds the [Cooper](https://twoboots.github.io/cooper) SDD infrastructure, installs the `battery` CLI binary, and guides you through initial `.batteryrc` workspace topology discovery.

## The Multi-Barrel Lifecycle

Battery orchestrates cross-repository features while preserving repository autonomy through a contract-first lifecycle:

<div class="lifecycle-flow">
  <div class="lifecycle-card">
    <span class="lifecycle-badge">Phase 1</span>
    <div class="lifecycle-step">battery track init</div>
    <div class="lifecycle-sub">Macro Contracts &amp; Spec Deltas</div>
  </div>
  <div class="lifecycle-arrow">→</div>
  <div class="lifecycle-card">
    <span class="lifecycle-badge">Phase 2</span>
    <div class="lifecycle-step">battery track dispatch</div>
    <div class="lifecycle-sub">Sync to Barrels (No plan.md)</div>
  </div>
  <div class="lifecycle-arrow">→</div>
  <div class="lifecycle-card">
    <span class="lifecycle-badge">Phase 3</span>
    <div class="lifecycle-step">git agent-start</div>
    <div class="lifecycle-sub">Isolated Barrel Worktrees</div>
  </div>
  <div class="lifecycle-arrow">→</div>
  <div class="lifecycle-card">
    <span class="lifecycle-badge">Phase 4</span>
    <div class="lifecycle-step">Autonomous TDD</div>
    <div class="lifecycle-sub">Red -> Green -> Refactor</div>
  </div>
  <div class="lifecycle-arrow">→</div>
  <div class="lifecycle-card">
    <span class="lifecycle-badge">Phase 5</span>
    <div class="lifecycle-step">PR &amp; git agent-stop</div>
    <div class="lifecycle-sub">Merge &amp; Teardown</div>
  </div>
</div>

1. `battery track init <track_id>` **Macro Orchestration**: Author system-level interface contracts (`design.md`) and living spec deltas at the Battery root.
2. `battery track dispatch <track_id>` **Decoupled Spec Dispatch**: Synchronize spec deltas to participating barrels while omitting `plan.md` to preserve local planning autonomy.
3. `git agent-start <track_id>` **Worktree Isolation**: Spawn isolated worktrees in each target barrel via [Troop](https://twoboots.github.io/troop) (`.worktrees/<track_id>`).
4. **Autonomous TDD Execution**: Barrel agents author local `plan.md` checklists and execute strict TDD using [Cooper](https://twoboots.github.io/cooper) SDD.
5. **Centralized Status & Teardown**: Monitor progress with `battery track status <track_id>`, submit PRs, and tear down worktrees with `git agent-stop <track_id>`.

## Learn More

- [Getting Started Guide](guide/getting-started.md)
- [Workflow & Lifecycle Details](guide/workflow.md)
- [Architecture & Topology Model](architecture.md)
- [Model Context Protocol (MCP) Server](mcp.md)
- [Installation & Setup Guide](installation.md)

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
  flex: 1 1 140px;
  min-width: 130px;
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
