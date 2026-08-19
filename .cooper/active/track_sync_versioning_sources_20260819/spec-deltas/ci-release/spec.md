# Spec Delta: CI & Release Workflows

## Capability: `ci-release`

### Requirement 2: Automated Cross-Platform Binary Releases

#### Scenario 2.2: Dynamic Release Version Extraction
- **GIVEN** a push to `main` without an explicit git tag
- **WHEN** the GitHub Actions release workflow compiles binaries
- **THEN** it MUST dynamically parse the version from `cmd/root.go` rather than hardcoding a default fallback (`1.0.0`).
