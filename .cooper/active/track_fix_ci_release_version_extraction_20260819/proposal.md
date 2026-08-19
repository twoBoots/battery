# Proposal: Fix Release Workflow Version Extraction Pattern

## Problem Statement
In `.github/workflows/release.yml`, the command:
`VERSION=$(grep 'Version = ' cmd/root.go | sed -E 's/.*"([^"]+)".*/\1/')`
matched both `Version = "1.4.0"` and `RootCmd.Version = Version`, creating a multi-line string with a newline. When passed to `-ldflags`, the newline corrupted the Go linker command arguments, failing the release build with `usage: link [options] main.o`.

## Proposed Solution
Update the grep pattern to strictly match only the top-level `Version = "..."` declaration:
`VERSION=$(grep -E '^\s*Version\s*=' cmd/root.go | head -n 1 | sed -E 's/.*"([^"]+)".*/\1/')`
