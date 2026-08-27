# Capability Specification: Barrel Configuration & Project Discovery

## Description
Provides project structure detection, interactive initialization, layered configuration storage (`.batteryrc` committed + `.batteryrc.local` git-ignored), and barrel path management (add, remove, list) across multi-repository, monorepo, and custom workspace topologies, respecting each barrel's individual Cooper tech stack specification (`.cooper/definition/tech-stack.md`).

## Requirements

### Requirement 1: Layered Configuration Schema (`.batteryrc` & `.batteryrc.local`)
The system MUST support a layered configuration model where canonical project settings are stored in `.batteryrc` (committed to git) and developer-specific path overrides are stored in `.batteryrc.local` (git-ignored).

#### Scenario 1.1: Canonical Configuration Schema (`.batteryrc`)
- **GIVEN** a battery project workspace
- **WHEN** configuration is initialized or saved
- **THEN** `.batteryrc` MUST adhere to the schema:
  - `version`: semantic version string (e.g. `"1.0.0"`)
  - `structure`: `"multi-repo"` | `"monorepo"` | `"custom"`
  - `barrels`: Array of barrel objects, each containing:
    - `name`: string identifier for the barrel
    - `path`: relative or absolute directory path to the barrel
    - `type`: optional `"barrel"` (default) or `"battery"` (for hierarchical sub-batteries)
    - `role`: optional string describing domain responsibility
    - `tech`: optional string summarizing runtime and tech stack
    - `docs`: optional path to orchestrator-level documentation profile
    - `jira`: optional issue tracker / project key mapping
    - Custom dynamic attributes MUST be preserved round-trip during serialization.

#### Scenario 1.2: Local Overrides Merging (`.batteryrc.local`)
- **GIVEN** an existing `.batteryrc` and a local `.batteryrc.local` file
- **WHEN** configuration is loaded at runtime
- **THEN** settings and barrel paths in `.batteryrc.local` MUST override corresponding barrels in `.batteryrc` by barrel `name`, and `.batteryrc.local` MUST be ignored by git.

### Requirement 2: Reading Barrel Tech Stack & Context Metadata
The system MUST respect each respective barrel's Cooper tech stack definition when present, and seamlessly fall back to `.batteryrc` metadata and orchestrator documentation profiles when managing non-Cooper barrels.

#### Scenario 2.1: Barrel Tech Stack & Context Fallback Resolution
- **GIVEN** a registered barrel path and configuration
- **WHEN** inspecting barrel details or running multi-barrel status/planning
- **THEN** battery MUST resolve the tech stack and context using the fallback hierarchy:
  1. Barrel's `.cooper/definition/tech-stack.md` (or `conductor/tech-stack.md` / `tech-stack.md`)
  2. `.batteryrc` / `.batteryrc.local` inline `tech` and `role` fields
  3. Orchestrator-level barrel documentation profile (`docs/barrels/<name>.md` or `.cooper/barrels/<name>.md`)
  4. Default message: `"No Cooper tech-stack or profile defined"`.

### Requirement 3: Project Structure Discovery
The system MUST detect whether the surrounding environment is a multi-repo landscape (sibling repositories), a monorepo workspace with sub-packages, or a custom layout.

#### Scenario 3.1: Sibling Multi-Repo Detection
- **GIVEN** sibling directories containing `.git` or `.cooper` repositories
- **WHEN** project discovery is executed
- **THEN** candidate sibling repositories MUST be identified and suggested as barrels with structure `"multi-repo"`.

#### Scenario 3.2: Monorepo Package Detection
- **GIVEN** subdirectories within folders like `packages/`, `apps/`, `services/`, `libs/`, or root directories containing project markers (`package.json`, `go.mod`, `Cargo.toml`, `pyproject.toml`, `requirements.txt`)
- **WHEN** project discovery is executed
- **THEN** package subdirectories MUST be identified and suggested as barrels with structure `"monorepo"`.

#### Scenario 3.3: Custom / Hybrid / Hierarchical Layout
- **GIVEN** user-specified non-standard, heterogeneous, or nested battery directories
- **WHEN** custom structure is selected
- **THEN** barrel paths MUST be configured and resolved relative to the battery root.

### Requirement 4: Interactive & Non-Interactive Initialization
The system MUST support both interactive prompting and non-interactive/CI initialization.

#### Scenario 4.1: Interactive Multi-Select Prompting
- **GIVEN** a terminal session with interactive input or piped installer execution (`curl ... | bash`)
- **WHEN** `battery init` or `install.sh` is run
- **THEN** it MUST interactively prompt the user:
  - Select project structure (`multi-repo`, `monorepo`, or `custom`) using arrow keys
  - Interactively select candidate barrel folders with arrow keys and spacebar checkboxes, with candidates unselected by default so the user explicitly chooses which folders to register
  - Allow adding additional custom folder paths interactively if needed
  - Write the configuration to `.batteryrc` (or `.batteryrc.local` if `--local` is specified).

#### Scenario 4.2: Non-Interactive Setup
- **GIVEN** `CI=true` or `--non-interactive` or `-y` flag
- **WHEN** `battery init` is executed
- **THEN** it MUST automatically apply detected structure and barrels or write default configuration without prompting.

### Requirement 5: Barrel Path Management (Add, Remove, List)
The system MUST provide commands and programmatic APIs to manage barrels across `.batteryrc` and `.batteryrc.local`.

#### Scenario 5.1: Add Barrel
- **GIVEN** a valid directory path
- **WHEN** `battery barrel add <path> [--name <name>] [--type <barrel|battery>] [--role <role>] [--tech <tech>] [--docs <docs>] [--jira <jira>] [--local]` is executed
- **THEN** the barrel path MUST be normalized, name inferred or set, and appended to the target configuration file (`.batteryrc` or `.batteryrc.local`) with any specified metadata attributes without duplicates.

#### Scenario 5.2: Remove Barrel
- **GIVEN** a barrel name or path
- **WHEN** `battery barrel remove <identifier> [--local]` is executed
- **THEN** the matching barrel MUST be removed from the configuration.

#### Scenario 5.3: List Barrels
- **GIVEN** existing `.batteryrc` and optional `.batteryrc.local`
- **WHEN** `battery barrel list` or `battery status` is executed
- **THEN** all effective barrels MUST be displayed with their name, path, source (`canonical` vs `local override`), existence status, metadata attributes (`role`, `tech`, `docs`, `jira`), and resolved hybrid tech stack / profile summary.

### Requirement 6: Per-Barrel Cooper Tech Stack Scaffolding & Auto-Inference
The system MUST provide automated discovery, inference, and generation of per-barrel Cooper tech stack definitions (`.cooper/definition/tech-stack.md`) for independent repositories and monorepo sub-packages.

#### Scenario 6.1: Automated Marker-Based Inference
- **GIVEN** a barrel directory containing project markers (such as `go.mod`, `package.json`, `Cargo.toml`, or `pyproject.toml`)
- **WHEN** tech stack inference is executed
- **THEN** the system MUST detect the primary programming language (inferring `Go 1.27+` for Go modules), recommend test runners, linters, and coverage threshold defaults (>80%).

#### Scenario 6.2: Scaffolding `.cooper/definition/tech-stack.md`
- **GIVEN** a target barrel path and optional specification overrides (language, framework, test runner, linter)
- **WHEN** `battery barrel init <path|name>` CLI command or scaffolding API is invoked
- **THEN** it MUST create `<barrel_path>/.cooper/definition/tech-stack.md` formatted to Cooper specifications, protecting existing files unless `--force` is specified.

### Requirement 7: Native Agent Skills & Two-Tier Planning Availability
The system MUST provide the complete suite of 6 packaged project-local Cooper skills in `.agents/skills/cooper-{setup,rfc,new-track,implement,review,status}` supporting the Two-Tier SDD planning architecture (collaborative upstream RFCs and downstream tactical TDD tracks).

#### Scenario 7.1: Project-Local Skills Discovery & RFC Support
- **GIVEN** an AI coding assistant operating in the repository
- **WHEN** discovering project capabilities
- **THEN** skills in `.agents/skills/cooper-{setup,rfc,new-track,implement,review,status}` MUST be accessible and self-contained, enabling both `cooper-rfc` collaborative architectural workflows and `cooper-new-track` tactical implementation tracks.

#### Scenario 7.2: Automated RFC Approval Detection & Reviewer Guidance
- **GIVEN** an active Draft RFC PR created via `cooper-rfc`
- **WHEN** PR feedback is inspected or approval is evaluated
- **THEN** `cooper-rfc` MUST support automated detection of approvals via GitHub native review decisions (`reviewDecision == "APPROVED"`) or `/approve` comment triggers, embed reviewer action instructions in PR templates, and transition approved RFCs to ready for merge with registered child tracks.

### Requirement 8: Orchestrator-Level Barrel Documentation & Profile Scaffolding (`docs/barrels/`)
The system MUST support lightweight, zero-intrusion documentation profiles for non-Cooper barrels stored directly in the Battery root repository.

#### Scenario 8.1: Standard Profile Discovery & Custom Docs Resolution
- **GIVEN** a registered barrel without `.cooper/`
- **WHEN** context or documentation is queried
- **THEN** Battery MUST discover profiles located at `docs/barrels/<name>.md`, `.cooper/barrels/<name>.md`, or the custom path specified in `docs`.

#### Scenario 8.2: Profile Scaffolding CLI Subcommand
- **GIVEN** a barrel name or target identifier
- **WHEN** `battery barrel doc init <name>` or `battery barrel profile init <name>` is executed
- **THEN** Battery MUST scaffold a starter Markdown profile at `docs/barrels/<name>.md` containing structured sections for Role & Responsibilities, Tech Stack & Runtime, Development & Build Commands, Interface Contracts, and AI Agent Guidelines without modifying the target repository.


