# Proposal: CLI Self-Updater & Release Versioning Enhancement

## Context & Motivation
Battery is evolving rapidly with multiple active tracks introducing new capabilities, specs, and commands. While GitHub Actions compiles and uploads cross-platform binaries to GitHub Releases on tags and main branch pushes, developers and agent workflows currently lack a first-class in-CLI mechanism to discover and install new updates seamlessly. Currently, users must re-run `curl | bash` or manually download release assets.

Additionally, ensuring strict semantic version stamping, clear update checking, and frictionless binary replacement is essential as new tracks land in `main`.

## Proposed Solution
1. **Self-Update Subcommand (`battery update` / `battery self-update` alias)**:
   - Provide an in-tool updater that queries GitHub Releases API for `twoBoots/battery`.
   - Compare current `cmd.Version` with the latest published release.
   - Download the target platform binary (`battery-<os>-<arch>`) to a temp file, set executable permissions, and atomically replace the current running executable resolved via `os.Executable()`.
   - Support `--check` (`-c`) to check for available updates without applying them.
   - Support `--force` (`-f`) to reinstall or overwrite even if already on the latest version.
   - Support `--version <tag>` to target a specific release version.
2. **Release Versioning & Build Improvements**:
   - Maintain robust semver comparison logic supporting `v` prefixes and prereleases.
   - Ensure clear, actionable error reporting if file permissions prevent in-place binary overwrite.
   - Refine release workflow validation.

## Out of Scope
- Package manager packaging (e.g. Homebrew tap, Scoop, Apt repo) — deferred to future track.
