# Track Design: MCP Framework & Standards Upgrade Guide

## Architecture Overview

```mermaid
flowchart TD
    subgraph Battery Binary [Battery Go Binary (v1.x.x)]
        E["Embedded Upstream Templates (embed.FS)\n- .agents/skills/cooper-*\n- .cooper/COOPER.md, BATTERY.md\n- .cooper/definition/workflow.md"]
        F["Framework Analyzer Engine (internal/framework)"]
        F --- E
    end

    subgraph MCP Server [Battery MCP Server (stdio)]
        T1["Tool: battery_framework_status"]
        T2["Tool: battery_get_template"]
        R1["Resource: battery://framework-status"]
        R2["Resource: battery://templates/{name}"]
        P1["Prompt: guide_framework_upgrade_track"]
        T1 & T2 & R1 & R2 & P1 --- F
    end

    subgraph AI Assistant [AI Coding Assistant (Antigravity / Cursor / Claude)]
        A["Agent inspects status & reads templates"]
        A --> T1 & T2 & R1 & R2 & P1
        A --> B["Initiates dedicated Cooper Track: track_upgrade_cooper_1_x_x"]
        B --> C["Troop Worktree: .worktrees/track_upgrade_cooper_1_x_x"]
        C --> D["Semantic 3-Way Merge (Preserves custom rules & adopts new standards)"]
    end
```

## Component Details

### 1. Embedded Templates & Framework Engine (`internal/framework`)
- **Package**: `internal/framework`
- **Embedding**: Go 1.16+ `//go:embed templates/*`
- **Functions**:
  - `GetTemplate(name string) (string, error)`: Retrieves template content by identifier (e.g. `skills/cooper-rfc`, `docs/BATTERY.md`).
  - `ListTemplates() []TemplateInfo`: Returns catalog of embedded templates with their upstream version and descriptions.
  - `InspectFrameworkStatus(cwd string, targetPath string) (*FrameworkStatusReport, error)`:
    - Scans `.cooper/` and `.agents/skills/` in the target workspace.
    - Computes file hashes against embedded templates to classify each file into:
      - `up_to_date`: File matches upstream template.
      - `customized_locally`: File differs from base upstream template, indicating team customizations.
      - `outdated`: File matches an older version's hash or lacks modern headers.
      - `missing`: Standard template not yet present in workspace.

### 2. MCP Tools
- **`battery_framework_status`**:
  - Arguments:
    - `barrel` (string, optional): Target registered barrel name or path. If omitted, checks current battery workspace root.
  - Returns: JSON representation of `FrameworkStatusReport`.
- **`battery_get_template`**:
  - Arguments:
    - `name` (string, required): Template name (e.g. `skills/cooper-review`, `docs/COOPER.md`, `definition/workflow.md`).
  - Returns: Text of the template and format type (`markdown`).

### 3. MCP Resources
- **`battery://framework-status`**:
  - Read active workspace framework alignment report.
- **`battery://templates/{name}`**:
  - Read specific upstream template content.

### 4. MCP Prompts
- **`guide_framework_upgrade_track`**:
  - Arguments:
    - `track_id` (string, optional): Suggested track ID (e.g. `track_upgrade_cooper_standards_20260820`).
    - `barrel` (string, optional): Target barrel name.
  - Prompt instructions:
    - Guides the agent through reading `battery_framework_status`, retrieving upstream template diffs with `battery_get_template`, analyzing team-specific customizations, formulating a track proposal & spec delta, and performing an isolated upgrade in a Troop worktree.
