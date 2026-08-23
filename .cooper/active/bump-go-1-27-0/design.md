# Design: Bump to Go 1.27.0

## Technical Architecture & Impact

### 1. Go Module & Compiler Baseline
- Update `go.mod` to specify `go 1.27.0`.
- Verify standard library imports and module dependencies build cleanly without deprecated API issues.

### 2. CI & Release Workflows
- Update `.github/workflows/ci.yml`:
  - `actions/setup-go@v6` parameter `go-version` configured to `"1.27"`.
- Update `.github/workflows/release.yml`:
  - `actions/setup-go@v6` parameter `go-version` configured to `"1.27"`.
  - Matrix cross-compilation builds (Linux, macOS, Windows) verified under Go 1.27 toolchain.

### 3. Tech Stack Specification & Inference Engine
- Update `.cooper/definition/tech-stack.md` to declare `Go 1.27+` for the CLI & Runtime tooling engine.
- Update `internal/techstack/scaffold.go`:
  - Modify `InferTechStack` Go rule to return `Language: "Go 1.27+"`.
  - Update unit test assertions in `internal/techstack/scaffold_test.go` and `internal/techstack/techstack_test.go`.

### 4. Living Spec Deltas
- `spec-deltas/ci-release/spec.md`: Document requirement update for CI & cross-platform releases to build on Go 1.27.
- `spec-deltas/barrel-config/spec.md`: Document requirement update for marker-based tech stack inference to infer `Go 1.27+`.
