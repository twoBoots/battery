# Capability Specification: Model Context Protocol (MCP) Server

## Description
Provides a native Model Context Protocol (MCP) server over standard input/output (stdio) conforming to the Model Context Protocol (JSON-RPC 2.0) powered by the standardized Bender MCP engine (`bender/pkg/mcp`). Exposes battery orchestration tools, living context resources, and prompt templates to AI coding assistants (e.g. Antigravity, Claude Code, Cursor, Windsurf).

## Requirements

### Requirement 1: Stdio MCP Server Execution (`battery mcp`)
The system MUST provide an MCP server command running over standard input/output (stdio) conforming to the Model Context Protocol (JSON-RPC 2.0).

#### Scenario 1.1: MCP Initialization
- **GIVEN** an MCP client connecting via stdin/stdout
- **WHEN** sending an `initialize` JSON-RPC request
- **THEN** the server MUST respond with server capabilities (`tools`, `resources`, `prompts`) and server info (`battery-mcp`).

#### Scenario 1.2: Stdio Message Loop & Error Resilience
- **GIVEN** a running MCP server session
- **WHEN** receiving JSON-RPC requests, notifications, or invalid payloads
- **THEN** valid requests MUST be dispatched to corresponding handlers and malformed requests MUST return standard JSON-RPC 2.0 error responses without crashing the server loop.

### Requirement 2: MCP Tools for Barrel & Track Operations
The system MUST expose Battery's core orchestration functions as callable MCP tools.

#### Scenario 2.1: `battery_status` Tool
- **GIVEN** a valid workspace directory
- **WHEN** `tools/call` is invoked with `name: "battery_status"`
- **THEN** it MUST return structured JSON detailing configuration status, registered barrels, connectivity, active tracks, and framework alignment status.

#### Scenario 2.2: `battery_list_barrels` Tool
- **GIVEN** a valid workspace directory
- **WHEN** `tools/call` is invoked with `name: "battery_list_barrels"`
- **THEN** it MUST return a list of barrels along with resolved Cooper tech stack summaries.

#### Scenario 2.3: `battery_init_track` Tool
- **GIVEN** arguments `track_id`, `barrels`, and optional `name`
- **WHEN** `tools/call` is invoked with `name: "battery_init_track"`
- **THEN** it MUST initialize `.cooper/active/<track_id>/` in the workspace with metadata, proposal, design, and plan files.

#### Scenario 2.4: `battery_dispatch_track` Tool
- **GIVEN** a valid `track_id`
- **WHEN** `tools/call` is invoked with `name: "battery_dispatch_track"`
- **THEN** it MUST seed target barrel worktree/active folders with spec deltas and metadata while preserving local planning autonomy.

#### Scenario 2.5: `battery_track_status` Tool
- **GIVEN** a valid `track_id`
- **WHEN** `tools/call` is invoked with `name: "battery_track_status"`
- **THEN** it MUST aggregate and return task progress and phase status across all participating barrels.

#### Scenario 2.6: `battery_init_barrel_tech_stack` Tool
- **GIVEN** a valid barrel name or directory path and optional parameters (`language`, `framework`, `test_runner`, `linter`, `coverage_threshold`, `force`)
- **WHEN** `tools/call` is invoked with `name: "battery_init_barrel_tech_stack"`
- **THEN** it MUST infer or apply the provided settings and scaffold `.cooper/definition/tech-stack.md` in the target barrel directory, returning structured confirmation of the generated files.

#### Scenario 2.7: `battery_framework_status` Tool
- **GIVEN** a workspace or registered barrel directory
- **WHEN** `tools/call` is invoked with `name: "battery_framework_status"` (and optional `barrel`)
- **THEN** it MUST compare local `.cooper/` and `.agents/skills/` against embedded upstream templates, returning structured status for each file (`up_to_date`, `customized_locally`, `outdated`, `missing`) and summary upgrade metadata.

#### Scenario 2.8: `battery_get_template` Tool
- **GIVEN** a valid template name (e.g. `skills/cooper-review`, `docs/BATTERY.md`)
- **WHEN** `tools/call` is invoked with `name: "battery_get_template"`
- **THEN** it MUST return the upstream template text content and description, or return a clear error if the template is not found.

### Requirement 3: MCP Resources for Living Context
The system MUST expose queryable URIs under `battery://` for workspace inspection.

#### Scenario 3.1: Listing and Reading Resources
- **GIVEN** a running MCP session
- **WHEN** querying `resources/list` or `resources/read` for `battery://topology`, `battery://barrels/{name}/tech-stack`, or `battery://tracks/{track_id}`
- **THEN** the server MUST return valid resource payloads with correct MIME types.

#### Scenario 3.2: Framework Alignment and Template Resources
- **GIVEN** an active MCP session
- **WHEN** querying `resources/read` for `battery://framework-status` or `battery://templates/{name}`
- **THEN** the server MUST return valid JSON framework status reports or raw template markdown content respectively.

### Requirement 4: MCP Prompts for Multi-Barrel Workflows
The system MUST expose standard prompt templates for guiding AI agents through multi-barrel track inception and framework upgrade workflows.

#### Scenario 4.1: Planning Prompt
- **GIVEN** a request for `prompts/get` with `name: "plan_multi_barrel_track"`
- **WHEN** the prompt is retrieved with argument `track_id`
- **THEN** it MUST return prompt messages instructing the agent on Cooper Hybrid and Troop multi-barrel protocols.

#### Scenario 4.2: Framework Upgrade Guide Prompt
- **GIVEN** a request for `prompts/get` with `name: "guide_framework_upgrade_track"`
- **WHEN** the prompt is retrieved with optional `track_id` and `barrel`
- **THEN** it MUST return prompt messages instructing the AI agent on how to inspect framework divergence, preserve local customizations, and guide the user through initializing an upgrade track in a Troop worktree.

### Requirement 5: MCP Client Auto-Configuration & Installation (`battery mcp install`)
The system MUST support automatic detection, interactive selection, and safe JSON configuration insertion of Battery MCP into supported AI assistants.

#### Scenario 5.1: Non-Interactive / Flag-Based Client Installation
- **GIVEN** a list of target clients (e.g. `--client cursor,claude-desktop,antigravity` or `--all`)
- **WHEN** `battery mcp install` is executed
- **THEN** it MUST create or update each client's configuration file with the `battery mcp` server definition without corrupting existing JSON keys. For Google Antigravity (`agy`), the target configuration file MUST be `~/.gemini/config/mcp_config.json` with display name referencing `mcp_config.json`.

#### Scenario 5.2: Interactive Client Discovery & Multi-Select
- **GIVEN** an interactive terminal session
- **WHEN** `battery mcp install` is executed without client flags
- **THEN** it MUST scan for known client configurations and present a multi-select prompt allowing the user to choose which clients to configure. Antigravity detection MUST check for `~/.gemini/config`, `~/.gemini`, `.agents`, and `.gemini`.

#### Scenario 5.3: Integration with `battery init`
- **GIVEN** an interactive `battery init` session
- **WHEN** workspace topology configuration completes
- **THEN** it MUST prompt the user whether to configure MCP for their AI assistant and invoke the installer upon confirmation.
