# Execution Plan: Model Context Protocol (MCP) Server Capability

Follow strict TDD (Red -> Green -> Refactor) and attach Git Notes to task commits. Maintain >80% test coverage.

---

## Phase 1: MCP Protocol Types & Stdio Server Core

- [x] Task 1.1: Define JSON-RPC 2.0 and MCP protocol message types (`Request`, `Response`, `Notification`, `Error`, Server Capabilities) in `internal/mcp/protocol.go` with unit tests. (2c4d7d5)
- [ ] Task 1.2: Implement stdio `Server` loop with initialization, ping, and request dispatch in `internal/mcp/server.go` with unit tests.
- [ ] Task 1.3: Add error resilience tests for malformed JSON, unknown methods, and graceful shutdown.

---

## Phase 2: Core MCP Tool Suite (`internal/mcp/tools.go`)

- [ ] Task 2.1: Implement `battery_status` and `battery_list_barrels` MCP tools with unit tests.
- [ ] Task 2.2: Implement `battery_init_track` and `battery_dispatch_track` MCP tools with unit tests.
- [ ] Task 2.3: Implement `battery_track_status` MCP tool with unit tests.

---

## Phase 3: MCP Resources & Prompts (`internal/mcp/resources.go`, `internal/mcp/prompts.go`)

- [ ] Task 3.1: Implement `resources/list` and `resources/read` for `battery://topology`, `battery://barrels/{name}/tech-stack`, and `battery://tracks/{track_id}` with unit tests.
- [ ] Task 3.2: Implement `prompts/list` and `prompts/get` for `plan_multi_barrel_track` with unit tests.

---

## Phase 4: CLI Command Suite (`cmd/mcp.go`)

- [ ] Task 4.1: Implement `battery mcp` / `battery serve` Cobra command routing to `internal/mcp.Server` with CLI tests.
- [ ] Task 4.2: Register `mcp` command in `cmd/root.go` and add help/flag options.

---

## Phase 5: Verification, Documentation, & Capabilities Spec Promotion

- [ ] Task 5.1: Run full test suite with coverage validation (`go test -v -cover ./...`) ensuring coverage >80%.
- [ ] Task 5.2: Update README and documentation with MCP integration examples for AI assistants (Cursor, Claude Code, Antigravity).
