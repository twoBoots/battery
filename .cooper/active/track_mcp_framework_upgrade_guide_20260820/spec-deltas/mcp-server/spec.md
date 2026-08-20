# Spec Delta: Model Context Protocol (MCP) Server

## Modifications to Capability: MCP Server

### Requirement 2: MCP Tools for Barrel & Track Operations

+ #### Scenario 2.7: `battery_framework_status` Tool
+ - **GIVEN** a workspace or registered barrel directory
+ - **WHEN** `tools/call` is invoked with `name: "battery_framework_status"` (and optional `barrel`)
+ - **THEN** it MUST compare local `.cooper/` and `.agents/skills/` against embedded upstream templates, returning structured status for each file (`up_to_date`, `customized_locally`, `outdated`, `missing`) and summary upgrade metadata.

+ #### Scenario 2.8: `battery_get_template` Tool
+ - **GIVEN** a valid template name (e.g. `skills/cooper-review`, `docs/BATTERY.md`)
+ - **WHEN** `tools/call` is invoked with `name: "battery_get_template"`
+ - **THEN** it MUST return the upstream template text content and description, or return a clear error if the template is not found.

### Requirement 3: MCP Resources for Living Context

+ #### Scenario 3.2: Framework Alignment and Template Resources
+ - **GIVEN** an active MCP session
+ - **WHEN** querying `resources/read` for `battery://framework-status` or `battery://templates/{name}`
+ - **THEN** the server MUST return valid JSON framework status reports or raw template markdown content respectively.

### Requirement 4: MCP Prompts for Multi-Barrel Workflows

+ #### Scenario 4.2: Framework Upgrade Guide Prompt
+ - **GIVEN** a request for `prompts/get` with `name: "guide_framework_upgrade_track"`
+ - **WHEN** the prompt is retrieved with optional `track_id` and `barrel`
+ - **THEN** it MUST return prompt messages instructing the AI agent on how to inspect framework divergence, preserve local customizations, and guide the user through initializing an upgrade track in a Troop worktree.
