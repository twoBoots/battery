# Battery Architecture & Topology 🏛️

**Battery** is an open, agent-agnostic and repository-agnostic multi-repository Specification-Driven Development (SDD) orchestration protocol. It coordinates multi-repository tracks, contracts, and living capability specs across a collection of independent repositories or monorepo packages (barrels 🛢️).

## 🛢️ Core Metaphor

<div class="arch-overview">
  <div class="arch-section">
    <span class="arch-badge">Orchestration Layer</span>
    <div class="arch-card primary">
      <div class="arch-title">🔋 Battery Orchestrator</div>
      <div class="arch-details"><code>.batteryrc</code> topology • Shared API / Event Contracts • Living Spec Deltas</div>
    </div>
  </div>
  <div class="arch-connector">↓ Coordinates via Cooper SDD &amp; Troop Worktrees ↓</div>
  <div class="arch-grid">
    <div class="arch-card">
      <span class="arch-subbadge">Barrel 1</span>
      <div class="arch-title">🛢️ Auth Service</div>
      <div class="arch-details">Go / Postgres Backend<br/><code>.worktrees/&lt;track_id&gt;</code></div>
    </div>
    <div class="arch-card">
      <span class="arch-subbadge">Barrel 2</span>
      <div class="arch-title">🛢️ Web Dashboard</div>
      <div class="arch-details">React / Vite Frontend<br/><code>.worktrees/&lt;track_id&gt;</code></div>
    </div>
    <div class="arch-card">
      <span class="arch-subbadge">Barrel 3</span>
      <div class="arch-title">🛢️ Billing API</div>
      <div class="arch-details">Node / Stripe Service<br/><code>.worktrees/&lt;track_id&gt;</code></div>
    </div>
  </div>
</div>

* **[Cooper](https://twoboots.github.io/cooper)**: The barrel maker managing SDD specifications (`.cooper/`) and worktree lifecycle for human developers and autonomous AI agents.
* **[Troop](https://twoboots.github.io/troop)**: Worktree isolation tool providing shared Git aliases for human developers and AI code monkeys working in isolated worktrees (`.worktrees/`).
* **Barrel**: An individual repository, package, or microservice within the system landscape (barrels 🛢️).
* **Battery**: A collection of barrels orchestrated together for cross-repository feature epics.

## ⚙️ Layered Configuration (`.batteryrc` & `.batteryrc.local`)

Battery adopts a two-layer configuration model:
1. **`.batteryrc` (Canonical Topology)**: Committed to version control, defining standard team barrels, relative paths, and workspace structure.
2. **`.batteryrc.local` (Local Overrides)**: Git-ignored file allowing individual developers to override barrel paths to match local filesystem directories.

### Example `.batteryrc`
```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "version": "1.0.0",
  "structure": "multi-repo",
  "barrels": [
    {
      "name": "auth-service",
      "path": "../auth-service"
    },
    {
      "name": "web-dashboard",
      "path": "../web-dashboard"
    }
  ]
}
```

## 🗺️ Supported Topologies

Battery supports multiple workspace structures out of the box:

### 1. Multi-Repo (`"structure": "multi-repo"`)
Barrels reside as sibling repositories on the filesystem:
```text
workspace/
├── battery/           # Orchestrator root (.batteryrc)
├── auth-service/      # Barrel 1 (sibling repo)
└── web-dashboard/     # Barrel 2 (sibling repo)
```

### 2. Monorepo (`"structure": "monorepo"`)
Barrels reside as packages inside subdirectories within a single repository:
```text
monorepo-root/
├── .batteryrc         # Orchestrator topology
├── packages/
│   ├── auth/          # Barrel 1 (package)
│   └── ui/            # Barrel 2 (package)
└── apps/
    └── web/           # Barrel 3 (package)
```

### 3. Custom (`"structure": "custom"`)
Hybrid layouts, distributed disk paths, or Git submodules:
```text
custom-workspace/
├── battery/           # Orchestrator root (.batteryrc)
├── core/
│   └── api-gateway/   # Barrel 1 (nested submodule)
└── external/
    └── ../services/payment/  # Barrel 2 (arbitrary relative path)
```

### 4. Hierarchical Sub-Batteries (`"type": "battery"`)
A barrel can itself be a composite battery orchestrating nested sub-barrels:
```text
enterprise-workspace/
├── platform-battery/  # Root Battery Orchestrator
├── commerce-battery/  # Sub-Battery (Barrel of platform, orchestrates sub-barrels)
│   ├── .batteryrc     # Sub-Battery topology
│   ├── cart-service/  # Child barrel
│   └── checkout/      # Child barrel
└── auth-service/      # Sibling barrel
```

## 🛠️ Decoupled Barrel Tech Stacks

Battery intentionally avoids centralizing tech stack definitions. Instead, Battery dynamically resolves each barrel's local [Cooper](https://twoboots.github.io/cooper) definition:
```
<barrel_path>/.cooper/definition/tech-stack.md
```

Each barrel maintains complete autonomy:
* **Backend Barrel**: Go, Postgres, Docker, GitHub Actions, `golangci-lint`.
* **Frontend Barrel**: TypeScript, React, Vite, Vitest, ESLint.

When Battery dispatches track requirements, barrel agents inspect their own local tech stack and author TDD tasks matching their specific language idioms.

## 📐 Decoupled Planning Boundary

<div class="boundary-cards">
  <div class="boundary-card">
    <span class="boundary-badge">Macro Orchestration</span>
    <h3 class="boundary-title">🔋 Battery Orchestrator</h3>
    <p class="boundary-scope"><strong>Scope:</strong> Cross-Barrel (System-Wide)</p>
    <p class="boundary-desc"><strong>Defines WHAT &amp; WHY:</strong> System-level requirements, OpenAPI schemas, shared protobufs, cross-repo milestone gates, and living Spec Deltas.</p>
    <div class="boundary-artifacts">
      <span class="artifact-tag">proposal.md</span>
      <span class="artifact-tag">design.md</span>
      <span class="artifact-tag">spec-deltas/</span>
    </div>
  </div>
  <div class="boundary-card">
    <span class="boundary-badge">Local Implementation</span>
    <h3 class="boundary-title">🛢️ Target Barrels</h3>
    <p class="boundary-scope"><strong>Scope:</strong> Single Barrel (Isolated Worktree)</p>
    <p class="boundary-desc"><strong>Defines HOW:</strong> Internal component architecture, database migrations, testing doubles, and granular TDD task checklists.</p>
    <div class="boundary-artifacts">
      <span class="artifact-tag">.worktrees/&lt;track_id&gt;/</span>
      <span class="artifact-tag">plan.md (TDD)</span>
    </div>
  </div>
</div>

By isolating local implementation plans to barrel worktrees managed by [Troop](https://twoboots.github.io/troop), Battery prevents AI agent hallucinations and context window overflow.

## 🔗 Related Resources

* [Getting Started Guide](guide/getting-started.md)
* [Multi-Barrel Workflow Guide](guide/workflow.md)
* [Model Context Protocol (MCP) Server](mcp.md)
* [Installation Guide](installation.md)
* [Cooper SDD Framework](https://twoboots.github.io/cooper)
* [Troop Worktree Isolation](https://twoboots.github.io/troop)

<style>
.arch-overview {
  background-color: var(--vp-c-bg-soft);
  border: 1px solid var(--vp-c-divider);
  border-radius: 12px;
  padding: 16px;
  margin: 20px 0;
}
.arch-section {
  text-align: center;
}
.arch-badge {
  display: inline-block;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: var(--vp-c-brand-1);
  margin-bottom: 6px;
}
.arch-subbadge {
  display: inline-block;
  font-size: 11px;
  font-weight: 600;
  color: var(--vp-c-brand-1);
  margin-bottom: 4px;
}
.arch-card {
  background-color: var(--vp-c-bg);
  border: 1px solid var(--vp-c-divider);
  border-radius: 8px;
  padding: 12px;
  text-align: center;
  box-sizing: border-box;
}
.arch-card.primary {
  max-width: 480px;
  margin: 0 auto;
}
.arch-title {
  font-size: 14px;
  font-weight: 700;
  color: var(--vp-c-text-1);
  margin-bottom: 4px;
}
.arch-details {
  font-size: 12px;
  color: var(--vp-c-text-2);
  line-height: 1.4;
}
.arch-connector {
  text-align: center;
  font-size: 12px;
  font-weight: 600;
  color: var(--vp-c-text-3);
  margin: 12px 0;
  user-select: none;
}
.arch-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 12px;
}
.boundary-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 16px;
  margin: 20px 0;
}
.boundary-card {
  background-color: var(--vp-c-bg-soft);
  border: 1px solid var(--vp-c-divider);
  border-radius: 10px;
  padding: 16px;
  box-sizing: border-box;
}
.boundary-badge {
  display: inline-block;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: var(--vp-c-brand-1);
  margin-bottom: 6px;
}
.boundary-title {
  margin: 0 0 8px 0;
  font-size: 16px;
  font-weight: 700;
  color: var(--vp-c-text-1);
}
.boundary-scope {
  font-size: 13px;
  color: var(--vp-c-text-2);
  margin: 0 0 6px 0;
}
.boundary-desc {
  font-size: 13px;
  line-height: 1.5;
  color: var(--vp-c-text-1);
  margin: 0 0 12px 0;
}
.boundary-artifacts {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}
.artifact-tag {
  font-family: var(--vp-font-family-mono);
  font-size: 11px;
  background-color: var(--vp-c-bg-alt);
  border: 1px solid var(--vp-c-divider);
  border-radius: 4px;
  padding: 2px 6px;
  color: var(--vp-c-brand-1);
}
</style>
