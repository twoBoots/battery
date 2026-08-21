# Spec Delta: Model Context Protocol (MCP) Server (Bender Integration)

## Capability: `mcp-server`

### Requirement 1: Stdio MCP Server Execution (`battery mcp`)
~ The system MUST utilize `bender/pkg/mcp` as the JSON-RPC 2.0 stdio engine while retaining Battery domain capabilities.

#### Scenario 1.1: MCP Initialization
- **GIVEN** an MCP client connecting via stdin/stdout
- **WHEN** sending an `initialize` JSON-RPC request
- **THEN** the server MUST respond with server capabilities (`tools`, `resources`, `prompts`) and server info (`battery-mcp`) powered by Bender engine.

### Requirement 5: MCP Client Auto-Configuration & Installation (`battery mcp install`)
~ The system MUST delegate MCP client auto-configuration and discovery to `bender/pkg/mcp.InstallClients`.

#### Scenario 5.1: Non-Interactive / Flag-Based Client Installation
- **GIVEN** target client identifiers (e.g. `--client cursor,claude-desktop,antigravity` or `--all`)
- **WHEN** `battery mcp install` is executed
- **THEN** it MUST delegate configuration updates to Bender's multi-client installer engine, preserving existing configurations.
