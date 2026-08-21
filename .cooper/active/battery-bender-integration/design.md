# Design: Battery Bender CLI Integration

## Architectural Overview

`battery` transitions from maintaining internal updater and MCP protocol engines to importing them directly from `github.com/twoBoots/bender`.

```mermaid
flowchart TD
    subgraph Bender Core Library
        B_UPD["bender/pkg/updater<br/>• SelfUpdate()<br/>• CompareVersions()"]
        B_MCP["bender/pkg/mcp<br/>• Server (JSON-RPC 2.0 stdio)<br/>• RegisterTool / RegisterResource / RegisterPrompt<br/>• InstallClients() & GetSupportedClients()"]
    end

    subgraph Battery CLI Layer
        CMD_UPD["cmd/update.go<br/>(Calls updater.SelfUpdate)"]
        CMD_MCP["cmd/mcp.go<br/>(Calls mcp.InstallClients & Server.Serve)"]
    end

    subgraph Battery Domain MCP Layer
        BAT_SRV["internal/mcp/server.go<br/>(Creates bender.mcp.NewServer & registers handlers)"]
        BAT_TOOLS["internal/mcp/tools.go<br/>(battery_status, battery_dispatch_track, etc.)"]
        BAT_RES["internal/mcp/resources.go<br/>(battery://topology, etc.)"]
        BAT_PROMPTS["internal/mcp/prompts.go<br/>(plan_multi_barrel_track, etc.)"]
    end

    B_UPD --> CMD_UPD
    B_MCP --> CMD_MCP
    B_MCP --> BAT_SRV
    BAT_SRV --> BAT_TOOLS
    BAT_SRV --> BAT_RES
    BAT_SRV --> BAT_PROMPTS
```

## Component Breakdown

### 1. `cmd/update.go`
- Refactored to delegate directly to `updater.SelfUpdate`:
  ```go
  opts := updater.Options{
      Repo:           "twoBoots/battery",
      BinaryName:     "battery",
      CurrentVersion: Version,
      TargetVersion:  targetVersion,
      Force:          force,
      CheckOnly:      checkOnly,
  }
  res, err := updater.SelfUpdate(opts)
  ```
- Flags preserved: `--check` (`-c`), `--force` (`-f`), `--target-version` (`-t`).

### 2. `internal/mcp/` Refactoring
- **Deleted redundant protocol/installer files**:
  - `internal/mcp/protocol.go` & `protocol_test.go`
  - `internal/mcp/installer.go` & `installer_test.go`
- **Updated `internal/mcp/server.go`**:
  - Initializes `srv := mcp.NewServer("battery-mcp", version, cwd)`.
  - Re-registers Battery domain tools, resources, and prompts using Bender's `Tool`, `Resource`, `Prompt`, and handler interfaces.
- **Updated `cmd/mcp.go`**:
  - Server execution: Starts `internal/mcp.NewServer()` and serves over `os.Stdin` / `os.Stdout`.
  - Client installation: Uses `mcp.InstallClients()` and `mcp.GetSupportedClients()`.

### 3. Cleanup & Deletion of `internal/updater/`
- Remove `internal/updater/platform.go`, `release.go`, `semver.go`, `updater.go`, and their corresponding test files.

## Testing & Quality Strategy
- Update CLI unit tests in `cmd/update_test.go` and `cmd/mcp_test.go` to mock or verify interactions with Bender library functions.
- Update `internal/mcp/*_test.go` to verify tool, resource, and prompt registrations and execution through Bender's server.
- Run `go test -v -coverprofile=coverage.out ./...` to ensure coverage exceeds 80%.
