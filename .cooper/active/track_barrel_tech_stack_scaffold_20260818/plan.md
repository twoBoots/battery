# Execution Plan: Per-Barrel Tech Stack Scaffolding & Auto-Inference Engine (CLI + MCP)

Follow strict TDD (Red -> Green -> Refactor) and attach Git Notes to task commits. Maintain >80% test coverage.

---

## Phase 1: Tech Stack Inference & Scaffolding Engine (`internal/techstack/scaffold.go`)

- [x] Task 1.1: Implement marker-based inference logic (`InferTechStack(barrelPath)`) supporting Go, TypeScript/JS, Rust, Python, Deno, Java with unit tests. (c02b195)
- [x] Task 1.2: Implement `ScaffoldBarrelTechStack(barrelPath, opts)` to write `.cooper/definition/tech-stack.md` and styleguides with unit tests. (c02b195)
- [x] Task 1.3: Add tests validating collision prevention (not overwriting without `force: true`). (c02b195)

---

## Phase 2: CLI Integration (`cmd/barrel.go`)

- [x] Task 2.1: Add `battery barrel init <path|name>` Cobra subcommand with flags (`--language`, `--framework`, `--test-runner`, `--linter`, `--coverage-threshold`, `--force`) and tests. (03bbb00)

---

## Phase 3: MCP Tool Integration (`internal/mcp/tools.go`)

- [x] Task 3.1: Implement `battery_init_barrel_tech_stack` MCP tool with parameter parsing, inference, and execution with unit tests. (3bcb88d)
- [x] Task 3.2: Verify integration with dynamic resource `battery://barrels/{name}/tech-stack`. (3bcb88d)

---

## Phase 4: Verification, Documentation, & Capabilities Spec Promotion

- [x] Task 4.1: Run full test suite with coverage validation (`go test -v -cover ./...`) ensuring coverage >80%. (f13ab28)
- [x] Task 4.2: Update README, `docs/mcp-setup-guide.md`, and promote living capability specs (`.cooper/specs/barrel-config/spec.md` and `.cooper/specs/mcp-server/spec.md`). (f13ab28)
