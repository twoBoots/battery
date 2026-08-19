# Proposal: Upgrade GitHub Actions to Native Node 24 Versions

## Problem Statement
GitHub Actions workflow runs displayed deprecation warnings because `actions/setup-go@v5` and `actions/download-artifact@v6` targeted the deprecated Node 20 runtime.

## Proposed Solution
Upgrade all actions across `.github/workflows/ci.yml` and `.github/workflows/release.yml`:
- `actions/setup-go`: `v5` -> `v6`
- `actions/upload-artifact`: `v6` -> `v7`
- `actions/download-artifact`: `v6` -> `v8`
- `actions/checkout`: `v6` (already on Node 24)
