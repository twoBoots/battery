# Technical Design: CLI Self-Updater & Release Versioning Enhancement

## Architecture Overview

### 1. Package Structure
```
internal/updater/
  ├── updater.go       # Core updater service: Release querying, semver comparison, binary downloading & replacement
  ├── platform.go      # OS & architecture normalizer matching release binary naming
  └── updater_test.go  # Unit tests with httptest server mocking GitHub Releases API

cmd/
  ├── update.go        # Cobra command implementation for `battery update` & `battery self-update`
  └── update_test.go   # Command-line integration tests
```

### 2. Platform Binary Naming Convention
Matches `.github/workflows/release.yml`:
- `darwin / amd64` -> `battery-darwin-x86_64`
- `darwin / arm64` -> `battery-darwin-aarch64`
- `linux / amd64`  -> `battery-linux-x86_64`
- `linux / arm64`  -> `battery-linux-aarch64`
- `windows / amd64`-> `battery-windows-x86_64.exe`

### 3. Workflow & Safety
1. **Resolution**: Determine current binary location with `os.Executable()`, resolving symlinks with `filepath.EvalSymlinks`.
2. **Comparison**: Query `https://api.github.com/repos/twoBoots/battery/releases/latest` (or specific tag release) with standard User-Agent. Compare tag against current `cmd.Version`.
3. **Download & Atomic Replace**:
   - Stream download to temporary file in the same directory as target binary (ensuring same filesystem for atomic rename).
   - `chmod 0755` on Unix systems.
   - Replace active binary via atomic `os.Rename` or swap mechanism.
4. **Permissions Handling**: If write permissions fail (e.g. installed in `/usr/local/bin`), provide clear diagnostic with suggested resolution (e.g. `sudo battery update` or running installer).
