# Proposal: Bump to Go 1.27.0

## Intent
Upgrade the Battery toolchain, Go module definitions, GitHub Actions CI & Release workflows, Cooper tech stack platform contracts, and barrel tech stack scaffolding inference engine from Go 1.23 to Go 1.27.0.

## Motivation & Benefits
- **Modern Language & Compiler Features**: Adopt the latest Go compiler optimizations, performance enhancements, and standard library updates.
- **Unified Baseline**: Keep `go.mod`, GitHub Actions CI matrix, binary distribution pipelines, and living capability specifications aligned with the latest upstream Go release.
- **Scaffolding Accuracy**: Update Battery's automatic inference engine so newly scaffolded Go barrels default to the Go 1.27+ ecosystem standard.

## Scope Boundaries
- **In Scope**:
  - `go.mod` directive update to `go 1.27.0`.
  - `.github/workflows/ci.yml` and `.github/workflows/release.yml` Go version update to `1.27`.
  - `.cooper/definition/tech-stack.md` specification update to `Go 1.27+`.
  - `internal/techstack/scaffold.go` and associated unit tests updated to infer `Go 1.27+`.
  - Spec deltas for `ci-release` and `barrel-config`.
- **Out of Scope**:
  - Unrelated dependency major version breaking upgrades.
