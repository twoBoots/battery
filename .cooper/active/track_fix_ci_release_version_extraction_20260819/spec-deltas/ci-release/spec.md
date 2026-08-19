# Spec Delta: CI & Release Workflows

## Capability: `ci-release`

### Requirement 2: Automated Cross-Platform Binary Releases

#### Scenario 2.2: Dynamic Release Version Extraction
- **GIVEN** an automated release build on `main`
- **WHEN** resolving the version from `cmd/root.go`
- **THEN** the extraction MUST match strictly single-line variable assignment to ensure valid `-ldflags` string parameterization.
