# Spec Delta: Bump to Go 1.27.0

## Capability: `ci-release`

### Requirement: Continuous Integration Validation
The CI pipeline MUST validate formatting, linting, and test suites across all PRs and main branch commits using the Go 1.27 runtime environment.

#### Scenario: PR Lint and Test Check
- `-` THEN it MUST use modern Node 24-compatible actions (`actions/checkout@v6`), verify `deno fmt`, `deno lint`, and execute tests with coverage >80%.
- `+` THEN it MUST use modern Go 1.27 setup (`actions/setup-go@v6` with `go-version: "1.27"`), check formatting (`gofmt`), run linter (`go vet`), and execute tests with race detection and coverage >80%.

### Requirement: Automated Cross-Platform Binary Releases
The release pipeline MUST compile binaries for Linux, macOS, and Windows matrix targets using Go 1.27 and publish them to GitHub Releases.

#### Scenario: Go 1.27 Matrix Compilation
- `+` GIVEN a semantic tag `v*` or commit pushed to `main`
- `+` WHEN the release pipeline executes matrix compilation
- `+` THEN it MUST build cross-platform binaries (Linux x86_64/aarch64, macOS x86_64/arm64, Windows x86_64) using the Go 1.27 compiler.
