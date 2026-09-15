---
layout: home

hero:
  name: Battery
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
  - title: 🛢️ Multi-Barrel Coordination
    details: Orchestrate cross-cutting feature epics across independent repositories or monorepo packages (barrels) with unified lifecycle tracking.
  - title: 📐 Cooper SDD Integration
    details: Built upon [Cooper](https://github.com/twoBoots/cooper) Hybrid Spec-Driven Development, keeping living specs and capability requirements in sync across all barrels.
  - title: 🐒 Troop Worktree Isolation
    details: Dispatch track work directly into isolated Git worktrees powered by [Troop](https://github.com/twoBoots/troop) (`.worktrees/<track_id>`), protecting your trunk branches.
  - title: 🤖 Native MCP Server
    details: Built-in Model Context Protocol (MCP) server exposing tools, resources, and prompt templates directly to modern AI coding assistants.
  - title: ⚙️ Layered Configuration
    details: Canonical team topologies in `.batteryrc` with personal local overrides in `.batteryrc.local` for flexible developer environments.
  - title: 📜 Contract-First Decoupled Planning
    details: Author macro interface contracts and living spec deltas at the root while delegating localized TDD implementation plans to autonomous barrel agents.
---

## 🚀 Quickstart Installation

Install Battery into your workspace with a single command:

```bash
curl -fsSL https://raw.githubusercontent.com/twoBoots/battery/main/install.sh | bash
```

The installer verifies your environment, sets up [Troop](https://github.com/twoBoots/troop) Git aliases, scaffolds the [Cooper](https://github.com/twoBoots/cooper) SDD infrastructure, installs the `battery` CLI binary, and guides you through initial `.batteryrc` workspace topology discovery.

---

## 🔄 Multi-Barrel Lifecycle Workflow

Battery orchestrates cross-repository features while preserving repository autonomy through a contract-first lifecycle:

```mermaid
flowchart TD
    subgraph Root ["🔋 Battery Orchestrator"]
        A["battery track init &lt;track_id&gt;"] --> B["Author Macro Contracts & Living Spec Deltas"]
        B --> C["battery track dispatch &lt;track_id&gt;"]
    end

    subgraph BarrelA ["🛢️ Barrel A (e.g. Backend)"]
        C --> D1["git agent-start &lt;track_id&gt;"]
        D1 --> E1["Author Local plan.md (TDD)"]
        E1 --> F1["Execute TDD: Red -> Green -> Refactor"]
        F1 --> G1["Submit PR & git agent-stop &lt;track_id&gt;"]
    end

    subgraph BarrelB ["🛢️ Barrel B (e.g. Frontend)"]
        C --> D2["git agent-start &lt;track_id&gt;"]
        D2 --> E2["Author Local plan.md (TDD)"]
        E2 --> F2["Execute TDD: Red -> Green -> Refactor"]
        F2 --> G2["Submit PR & git agent-stop &lt;track_id&gt;"]
    end
```

Read the full [Getting Started Guide](guide/getting-started.md) or explore the [Workflow Guide](guide/workflow.md) for in-depth guidance.

