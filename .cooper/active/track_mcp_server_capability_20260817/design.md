# Technical Design: Model Context Protocol (MCP) Server Capability

## Architecture Overview

The `internal/mcp` package encapsulates the MCP protocol server and handlers, exposing Battery's domain packages (`internal/config`, `internal/techstack`, `internal/track`, `internal/discovery`) over the Model Context Protocol stdio transport.

```
cmd/
├── mcp.go                     # Cobra CLI command: battery mcp [--transport stdio]
└── root.go                    # Registers mcp command

internal/
├── config/                    # Configuration loading and barrel resolution
├── techstack/                 # Cooper tech-stack detector
├── track/                     # Track lifecycle & status engine
├── discovery/                 # Workspace topology discovery
└── mcp/                       # Model Context Protocol engine
    ├── protocol.go            # JSON-RPC 2.0 / MCP message types and protocol definitions
    ├── server.go              # Stdio MCP Server loop, lifecycle & error handlers
    ├── tools.go               # Implementation of battery_* MCP tools
    ├── resources.go           # Implementation of battery:// URI resource providers
    ├── prompts.go             # Implementation of prompt templates
    └── server_test.go         # Comprehensive unit tests (>80% coverage)
```

## Protocol & Message Handling

The MCP server adheres to the Model Context Protocol Specification (JSON-RPC 2.0 transport over stdin/stdout):

1. **Lifecycle Methods**:
   - `initialize`: Returns server info (`battery-mcp`, version `v1.2.0`), protocol version, and server capabilities (`tools`, `resources`, `prompts`).
   - `notifications/initialized`: Acknowledged by client to start normal operation.
   - `ping`: Standard liveness check returning `{}`.

2. **Tools (`tools/list` & `tools/call`)**:
   - `battery_status`: Calls `internal/config.LoadConfig` and checks barrel connectivity.
   - `battery_list_barrels`: Calls `internal/config.LoadConfig` + `internal/techstack.DetectTechStack`.
   - `battery_init_track`: Calls `internal/track.InitTrack(cwd, trackID, barrels, opts)`.
   - `battery_dispatch_track`: Calls `internal/track.DispatchTrack(cwd, trackID, opts)`.
   - `battery_track_status`: Calls `internal/track.GetMultiBarrelTrackStatus(cwd, trackID)`.

3. **Resources (`resources/list` & `resources/read`)**:
   - `battery://topology`: Returns JSON of active merged configuration.
   - `battery://barrels/{name}/tech-stack`: Returns the raw `.cooper/definition/tech-stack.md` content.
   - `battery://tracks/{track_id}`: Returns JSON representation of track status and spec deltas.

4. **Prompts (`prompts/list` & `prompts/get`)**:
   - `plan_multi_barrel_track`: Returns a structured prompt guiding agents on authoring spec-deltas and multi-barrel tracks.

5. **MCP Client Auto-Configuration (`internal/mcp/installer.go`)**:
   - Detects existing client configuration directories (Cursor `.cursor/`, Antigravity `.gemini/`, Claude Desktop, Claude Code `~/.claude.json`, Windsurf `~/.codeium/windsurf/`, VS Code `.vscode/`).
   - Safely reads, merges, and writes the `"battery"` entry under `"mcpServers"` in each client's config file without mutating other servers.
   - Provides interactive multi-select prompt when run interactively, and flag-based selection for scripted workflows (`battery mcp install --client=cursor,claude`).
   - Integrates with `battery init` to prompt developers at the completion of workspace setup.

