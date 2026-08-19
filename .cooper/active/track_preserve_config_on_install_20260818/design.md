# Technical Design: Preserve Existing Configuration on Battery Install & Init

## Overview
This design outlines the architecture for detecting existing configuration during `battery init` and `install.sh`, allowing users to either preserve current configuration (default) or overwrite and start clean.

## Architecture & Flow

### 1. Existing Configuration Detection (`cmd/init.go`)
- Target configuration file determination:
  - If `--local` is specified: check for `filepath.Join(cwd, config.LocalConfigFilename)` (`.batteryrc.local`).
  - Otherwise: check for `filepath.Join(cwd, config.ConfigFilename)` (`.batteryrc`).
- Check if file exists and contains valid configuration.

### 2. Interactive Flow
When running interactively (`!isNonInteractive`):
- If target config file exists:
  - Prompt user using `huh.Select[string]`:
    - Title: `"Existing Battery configuration detected in <filename>:"`
    - Option 1: `"Continue setup and preserve current config (Recommended)"` -> `action = "preserve"`
    - Option 2: `"Overwrite and start clean"` -> `action = "overwrite"`
  - If `action == "preserve"`:
    - Load existing config via `config.LoadConfig(cwd)` (or `config.LoadLocalConfig(cwd)`).
    - Output: `"  [✓] Preserving current configuration (<structure>, <N> barrels registered)"`.
    - Set `selectedStructure = existing.Structure`, `finalBarrels = existing.Barrels`.
    - Skip `promptStructure` and `promptBarrels`.
    - (Optionally touch/save to ensure latest schema/version, or keep intact).
    - Proceed to MCP prompt step if applicable.
  - If `action == "overwrite"`:
    - Output: `"  [!] Overwriting existing configuration..."`
    - Run standard `discovery.DiscoverCandidateBarrels(cwd)`, `promptStructure`, `promptBarrels`.
    - Save updated configuration.

### 3. Non-Interactive / CI Flow
- Flags added to `initCmd`:
  - `--force`, `-f`, `--overwrite`: Bool flag to force overwrite of existing configuration.
- Behavior:
  - If target config file exists:
    - If `--force` / `--overwrite` is set: proceed with topology discovery and save fresh config.
    - Otherwise (default): load existing config, log `"  [✓] Existing configuration detected; preserving current config."`, and keep configuration intact.

### 4. Installer Script (`install.sh`)
- Step 4 (`battery init`):
  - Pass along standard flags.
  - In binary fallback mode (if battery binary fails or isn't built yet):
    - If `.batteryrc` exists, log `"  [✓] Existing .batteryrc detected; preserving current configuration."` instead of overwriting.
    - If missing, create default `.batteryrc`.

## Test Plan
- Unit test for non-interactive init with existing `.batteryrc` verifying preservation by default.
- Unit test for non-interactive init with `--force` / `--overwrite` verifying overwrite.
- Unit test for interactive prompt simulation / helper functions.
