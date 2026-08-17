# Capability Specification: CI & Release Workflows

## Description
Automated continuous integration and binary release publishing pipelines for Battery using GitHub Actions, ensuring cross-platform binary compilation and zero runtime deprecation warnings.

## Requirements

### Requirement 1: Continuous Integration Validation
The CI pipeline MUST validate formatting, linting, and test suites across all PRs and main branch commits without deprecated Node runtime warnings.

#### Scenario 1.1: PR Lint and Test Check
- **GIVEN** a commit pushed to a branch or pull request targeting `main`
- **WHEN** GitHub Actions CI runs
- **THEN** it MUST use modern Node 24-compatible actions (`actions/checkout@v6`), verify `deno fmt`, `deno lint`, and execute tests with coverage >80%.

### Requirement 2: Automated Cross-Platform Binary Releases
The release pipeline MUST compile binaries for Linux, macOS, and Windows matrix targets and publish them to GitHub Releases.

#### Scenario 2.1: Release Asset Creation and Update
- **GIVEN** a push to `main` or a semantic tag `v*`
- **WHEN** `publish-release` runs
- **THEN** it MUST check if the release exists, creating it if absent (`gh release create`) or updating assets with clobber (`gh release upload --clobber`) if present, without invalid CLI flags.
