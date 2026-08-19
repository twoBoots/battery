# Spec Delta: Model Context Protocol (MCP) Server

## Capability: `mcp-server`

### Requirement 1: Stdio MCP Server Execution (`battery mcp`)

#### Scenario 1.1: MCP Initialization
- **GIVEN** an MCP client connecting via stdin/stdout
- **WHEN** sending an `initialize` JSON-RPC request
- **THEN** the server MUST respond with server capabilities (`tools`, `resources`, `prompts`) and server info (`battery-mcp`) with dynamic versioning matching the Battery CLI binary version (`cmd.Version`).

### Requirement 2: MCP Tools for Barrel & Track Operations

#### Scenario 2.1: `battery_status` Tool
- **GIVEN** a valid workspace directory
- **WHEN** `tools/call` is invoked with `name: "battery_status"`
- **THEN** it MUST return structured JSON detailing configuration status, registered barrels, connectivity, active tracks, CLI version (`cli_version`), and config schema version (`config_version`).
