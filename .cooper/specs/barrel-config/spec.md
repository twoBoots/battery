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

#### Scenario 1.2: Local Overrides Merging (`.batteryrc.local`)
- **GIVEN** an existing `.batteryrc` and a local `.batteryrc.local` file
- **WHEN** configuration is loaded at runtime
- **THEN** settings and barrel paths in `.batteryrc.local` MUST override corresponding barrels in `.batteryrc` by barrel `name`, and `.batteryrc.local` MUST be ignored by git.

### Requirement 2: Reading Barrel Tech Stack from Cooper (`.cooper/definition/tech-stack.md`)
The system MUST respect each respective barrel's Cooper tech stack definition rather than hardcoding it in `.batteryrc`.

#### Scenario 2.1: Barrel Tech Stack Resolution
- **GIVEN** a registered barrel path
- **WHEN** inspecting barrel details or running multi-barrel status/planning
- **THEN** battery MUST read the barrel's `.cooper/definition/tech-stack.md` (or `conductor/tech-stack.md`) to resolve its language, framework, and tooling specifications.

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
- **WHEN** `battery barrel add <path> [--name <name>] [--type <barrel|battery>] [--local]` is executed
- **THEN** the barrel path MUST be normalized, name inferred or set, and appended to the target configuration file (`.batteryrc` or `.batteryrc.local`) without duplicates.

#### Scenario 5.2: Remove Barrel
- **GIVEN** a barrel name or path
- **WHEN** `battery barrel remove <identifier> [--local]` is executed
- **THEN** the matching barrel MUST be removed from the configuration.

#### Scenario 5.3: List Barrels
- **GIVEN** existing `.batteryrc` and optional `.batteryrc.local`
- **WHEN** `battery barrel list` or `battery status` is executed
- **THEN** all effective barrels MUST be displayed with their name, path, source (`canonical` vs `local override`), existence status, and resolved Cooper tech stack summary.
