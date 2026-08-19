# Spec Delta: CI & Release Workflows

## Capability: `ci-release`

### Requirement 1: Continuous Integration Validation

#### Scenario 1.1: PR Lint and Test Check
- **GIVEN** a commit pushed to a branch or pull request targeting `main`
- **WHEN** GitHub Actions CI runs
- **THEN** it MUST use modern Node 24-compatible actions (`actions/checkout@v6`, `actions/setup-go@v6`), verify formatting, linting, and execute tests with coverage >80%.

### Requirement 2: Automated Cross-Platform Binary Releases

#### Scenario 2.3: Node 24 Native Artifact Packaging
- **GIVEN** a push to `main` or semantic release tag
- **WHEN** GitHub Actions release workflow packages cross-platform binaries
- **THEN** it MUST use modern Node 24-compatible actions (`actions/upload-artifact@v7`, `actions/download-artifact@v8`).
