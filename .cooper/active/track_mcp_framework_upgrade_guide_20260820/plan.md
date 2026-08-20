# Implementation Plan: MCP Framework & Standards Upgrade Guide

## Phase 1: Embedded Framework Templates & Inspection Engine (`internal/framework`)
- [~] Task 1.1: Framework Inspection Engine Unit Tests (TDD - Red)
  - [ ] Sub-task: Write unit tests for embedded template catalog retrieval and missing template error handling in `internal/framework/framework_test.go`
  - [ ] Sub-task: Write unit tests for workspace scanning and file status classification (`up_to_date`, `customized_locally`, `outdated`, `missing`)
- [ ] Task 1.2: Implement Embedded Templates & Framework Engine (TDD - Green)
  - [ ] Sub-task: Create `internal/framework/templates/` with embedded canonical Cooper skills and framework documentation
  - [ ] Sub-task: Implement `GetTemplate`, `ListTemplates`, and `InspectFrameworkStatus` in `internal/framework/framework.go`
  - [ ] Sub-task: Refactor & ensure coverage >80% in `internal/framework` (Refactor)
- [ ] Task 1.3: Phase 1 Verification & Checkpoint
  - [ ] Sub-task: Run unit tests and record git note checkpoint

## Phase 2: MCP Tools, Resources & Prompt Templates (`internal/mcp`)
- [ ] Task 2.1: MCP Tools & Resources Unit Tests (TDD - Red)
  - [ ] Sub-task: Write unit tests in `internal/mcp/framework_tools_test.go` for `battery_framework_status` and `battery_get_template`
  - [ ] Sub-task: Write unit tests for `battery://framework-status` and `battery://templates/{name}` resources
  - [ ] Sub-task: Write unit tests for `guide_framework_upgrade_track` prompt template
- [ ] Task 2.2: Implement MCP Handlers & Prompt Definition (TDD - Green)
  - [ ] Sub-task: Register `battery_framework_status` and `battery_get_template` in `internal/mcp/tools.go`
  - [ ] Sub-task: Register `battery://framework-status` and `battery://templates/{name}` in `internal/mcp/resources.go`
  - [ ] Sub-task: Register `guide_framework_upgrade_track` in `internal/mcp/prompts.go`
  - [ ] Sub-task: Enrich `battery_status` tool output with summary framework status metadata
  - [ ] Sub-task: Refactor & ensure coverage >80% across `internal/mcp` (Refactor)
- [ ] Task 2.3: Phase 2 Verification & Checkpoint
  - [ ] Sub-task: Run full test suite and record git note checkpoint

## Phase 3: End-to-End Verification & Living Spec Consolidation
- [ ] Task 3.1: Stdio MCP Protocol Integration Testing
  - [ ] Sub-task: Write end-to-end stdio JSON-RPC test verifying tools, resources, and prompt retrieval in a simulated client session
  - [ ] Sub-task: Run `go test ./...` with race detection and linter verification
- [ ] Task 3.2: Living Capability Spec Update
  - [ ] Sub-task: Merge spec delta into `.cooper/specs/mcp-server/spec.md`
- [ ] Task 3.3: Phase 3 Verification & Final Checkpoint
  - [ ] Sub-task: Verify all tests pass, stage track artifacts, and prepare track summary
