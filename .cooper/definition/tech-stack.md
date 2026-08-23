# Technology Stack & Platform Contracts

## Pattern Components

| Component | Role | Metaphor | Tech Stack / Tooling |
| :--- | :--- | :--- | :--- |
| **Multi-Barrel Orchestration** | `battery` | Collection of Barrels | Cooper SDD Framework (`.cooper/`), Living Specs & Plan Templates |
| **CLI & Runtime** | `battery` CLI | Tooling Engine | Go 1.27+ (`spf13/cobra`, `charmbracelet/huh`, `charmbracelet/lipgloss`, `go test`, `go vet`, `gofmt`) |
| **Worktree Manager & SDD** | `twoBoots/cooper` | The Barrel Maker | Cooper Hybrid SDD Workflow, Troop Git Worktree Aliases |
| **Individual Repos** | `barrel` | Single Barrel | Repositories / Microservices (independent tech stacks defined in their own `.cooper/`) |
| **Configuration** | `battery` Config | Local / Team Registry | JSON (`.batteryrc` committed, `.batteryrc.local` git-ignored) |
| **Interface Contracts** | Contracts | Shared Specs | OpenSpec Living Specs (`spec.md`), Spec Deltas (`spec-deltas/`), OpenAPI, Protobuf, MCP |
| **Track Metadata** | Registry | Battery Dashboard | JSON (`metadata.json`), Markdown (`plan.md`, `proposal.md`, `design.md`) |

## Multi-Barrel Worktree Integration
- **Isolation Directory**: `.worktrees/<track_id>` in each target barrel.
- **Workflow Tools**: Git worktree automation (`git agent-start`, `git agent-stop`, `git troop`).
- **Binary Tooling**: `bin/battery` compiled native binary via `go build`.
