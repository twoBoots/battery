# Track Proposal: Model Context Protocol (MCP) Server Capability

## Context & Motivation
`battery` coordinates multi-repository Specification-Driven Development (SDD) across a collection of barrels for human developers and autonomous AI agents. Currently, AI agents must interact with `battery` by executing shell CLI commands (`battery status`, `battery track init`, `battery track dispatch`, `battery barrel list`) and parsing unstructured terminal text.

With modern AI agent ecosystems natively adopting the Model Context Protocol (MCP), exposing `battery`'s orchestration engine as an MCP server enables direct, high-precision tool calling, structured JSON-RPC communication, and rich contextual resource resolution without fragile shell invocations.

## Objectives
1. **MCP Server Command (`battery mcp` / `battery serve`)**:
   - Provide a stdio-based MCP server enabling seamless integration with tools like Antigravity, Claude Code, Cursor, Windsurf, VS Code, and other MCP clients.
2. **First-Class MCP Tool Suite**:
   - `battery_status`: Report workspace topology, barrel connectivity, and active tracks.
   - `battery_list_barrels`: List all registered barrels with their resolved Cooper tech stack guidelines.
   - `battery_init_track`: Initialize multi-barrel tracks with target barrel selections.
   - `battery_dispatch_track`: Dispatch spec deltas and contracts into barrel worktrees.
   - `battery_track_status`: Query aggregated progress, task counts, and phase status across all barrels.
3. **Living Context Resources**:
   - Expose `battery://topology`, `battery://barrels/{name}/tech-stack`, and `battery://tracks/{track_id}` as queryable MCP resources.
4. **Planning Prompts**:
   - Provide guided prompt templates for multi-barrel track inception and contract definition.
5. **MCP Client Auto-Configuration (`battery mcp install`)**:
   - Provide automated and interactive configuration generation for AI coding assistants (Cursor, Antigravity, Claude Desktop, Claude Code, Windsurf, VS Code) and interactive prompts during `battery init`.

## Value & Impact
* Gives AI agents direct, structured programmatic control over multi-repository SDD workflows.
* Eliminates CLI output parsing errors and shell escaping issues.
* Accelerates parallel track execution across heterogeneous multi-barrel environments.
* Delivers a frictionless setup experience for AI developers across any IDE or agent environment.

