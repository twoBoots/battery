# Spec Delta: Battery Auto-Tag Release Pipeline

## Capability: `ci-release`

### Requirement: Automated Git SemVer Tagging on Merge to Main
The CI/CD release workflow SHALL automatically detect the application semantic version from `cmd/root.go` and publish a Git tag matching `v<Version>` upon push or merge to `main` if the tag does not already exist.

#### Scenario: Tagging on Merge to Main
- `+` GIVEN a pull request merged to `main` with `cmd.Version` set in `cmd/root.go`
- `+` WHEN the `Release Binary` GitHub Actions workflow triggers on `refs/heads/main`
- `+` THEN it MUST verify if `v<Version>` exists on `origin`
- `+` AND IF missing, it MUST create and push Git tag `v<Version>` to `origin`
- `+` AND it MUST publish the release assets under `v<Version>` as well as updating `latest`.

#### Scenario: Idempotent Tag Handling
- `+` GIVEN a push to `main` where Git tag `v<Version>` already exists on `origin`
- `+` WHEN the release workflow runs
- `+` THEN it MUST proceed without error, updating assets for `v<Version>` and `latest`.
