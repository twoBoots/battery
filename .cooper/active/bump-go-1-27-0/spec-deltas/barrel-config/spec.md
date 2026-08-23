# Spec Delta: Bump to Go 1.27.0

## Capability: `barrel-config`

### Requirement: Per-Barrel Cooper Tech Stack Scaffolding & Auto-Inference
The system MUST provide automated discovery, inference, and generation of per-barrel Cooper tech stack definitions (`.cooper/definition/tech-stack.md`) using Go 1.27+ conventions for Go projects.

#### Scenario: Automated Go Marker-Based Inference
- `-` THEN the system MUST detect the primary programming language as Go 1.23+, recommend test runners, linters, and coverage threshold defaults (>80%).
- `+` THEN the system MUST detect Go barrels containing `go.mod` and infer `Language: "Go 1.27+"`, test runner `go test -v -cover ./...`, linter `golangci-lint run`, and coverage threshold `>80%`.
