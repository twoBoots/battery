# Proposal: Preserve Existing Configuration on Battery Install & Init

## Rationale & Problem Statement
When Battery is installed via `install.sh` (or when `battery init` is executed) in a repository or directory where Battery configuration (`.batteryrc` or `.batteryrc.local`) already exists:
- Currently, `battery init` restarts topology discovery and prompts the user, overwriting existing `.batteryrc` settings and manually registered barrels unless carefully re-selected.
- In upgrade workflows (especially for versions prior to the self-update command, or when running `curl ... | bash` installer to install new Battery versions/scripts), the installer clobbers or forces re-configuration instead of preserving the existing barrel layout.

## Proposed Solution
1. **Interactive Detection & Choice**:
   - In interactive mode, if existing configuration (`.batteryrc` or `.batteryrc.local` when `--local` is passed) is detected, prompt the user with two clear options:
     1. **`Continue setup and preserve current config`** (Default / Recommended)
     2. **`Overwrite and start clean`** (Reconfigures topology and barrels from scratch)
   - When preserving, retain existing `structure` and `barrels`, update the file version/schema if necessary, display the preserved summary, and proceed to subsequent setup actions (such as MCP assistant configuration).
   - When overwriting, proceed with standard topology and candidate barrel selection prompts.

2. **Non-Interactive & CI Automation**:
   - In non-interactive mode (`--non-interactive`, `--yes`, `-y`, `CI=true`), if existing configuration is present, preserve current config by default without overwriting.
   - Introduce a `--force` / `-f` flag (and `--overwrite`) to allow non-interactive scripts to explicitly overwrite existing configuration when desired.

3. **Installer Resilience**:
   - Ensure `install.sh` properly delegates to `battery init` (which handles preservation/overwrite) and fallback defaults never overwrite an existing `.batteryrc`.

## User Benefit
- Seamless upgrades: Users can safely run `install.sh` to update Battery without losing their configured multi-barrel setups.
- Safer developer ergonomics: Accidental re-initialization won't wipe existing configuration.
