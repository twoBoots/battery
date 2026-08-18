# Spec Delta: Barrel Configuration & Project Discovery

## Capability: `barrel-config`

### Requirement 4: Interactive & Non-Interactive Initialization
The system MUST support both interactive prompting and non-interactive/CI initialization, with explicit handling for existing configuration files.

#### Scenario 4.1: Interactive Multi-Select Prompting
- **GIVEN** a terminal session with interactive input or piped installer execution (`curl ... | bash`)
- **WHEN** `battery init` or `install.sh` is run
- **THEN** it MUST interactively prompt the user:
+   - If existing configuration (`.batteryrc` or `.batteryrc.local`) is present:
+     - Prompt user with choice to "Continue setup and preserve current config" [default] or "Overwrite and start clean"
+     - If preserving: retain existing structure and barrels, skip topology prompts, and proceed to subsequent setup actions
+     - If overwriting: proceed to prompt structure and barrel discovery
    - Select project structure (`multi-repo`, `monorepo`, or `custom`) using arrow keys
    - Interactively select candidate barrel folders with arrow keys and spacebar checkboxes, with candidates unselected by default so the user explicitly chooses which folders to register
    - Allow adding additional custom folder paths interactively if needed
    - Write the configuration to `.batteryrc` (or `.batteryrc.local` if `--local` is specified).

#### Scenario 4.2: Non-Interactive Setup
- **GIVEN** `CI=true` or `--non-interactive` or `-y` flag
- **WHEN** `battery init` is executed
- **THEN** it MUST automatically apply detected structure and barrels or write default configuration without prompting:
+   - If existing configuration is present and `--force` / `--overwrite` is NOT provided, it MUST preserve existing configuration by default
+   - If `--force` or `--overwrite` is provided, it MUST overwrite with auto-discovered topology.
