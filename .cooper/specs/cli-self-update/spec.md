# Capability Specification: CLI Self-Updater

## Description
Provides an in-CLI update mechanism (`battery update` and alias `battery self-update`) powered by the standardized Bender updater engine (`bender/pkg/updater`) that discovers latest releases from GitHub, checks version status, downloads platform-compatible prebuilt binaries, and replaces the running binary with atomic rollback and macOS code-signing.

## Requirements

### Requirement 1: Update Discovery and Check
The CLI MUST provide a command to inspect GitHub Releases for available version updates.

#### Scenario 1.1: Checking for updates when up to date
- **GIVEN** a Battery CLI running version `1.2.0`
- **WHEN** running `battery update --check` (or `-c`)
- **THEN** it MUST report that Battery is already up to date with the latest version.

#### Scenario 1.2: Checking for updates when an update is available
- **GIVEN** a Battery CLI running version `1.0.0` and latest release `1.2.0`
- **WHEN** running `battery update --check`
- **THEN** it MUST output that an update is available (`1.0.0` -> `1.2.0`) without modifying the local binary.

### Requirement 2: Binary Self-Update Execution
The CLI MUST download the target release binary for the host OS and architecture and replace the existing executable.

#### Scenario 2.1: Applying an update
- **GIVEN** an installed Battery binary and an available newer release
- **WHEN** running `battery update`
- **THEN** it MUST download the appropriate binary asset, set executable permissions, replace the binary, and report success.

#### Scenario 2.2: Forcing an update / reinstall
- **GIVEN** the CLI is already at the latest release version
- **WHEN** running `battery update --force` (or `-f`)
- **THEN** it MUST re-download and replace the binary regardless of current version match.

#### Scenario 2.3: Targeting a specific version
- **GIVEN** an installed Battery binary
- **WHEN** running `battery update --target-version <tag>` (or `-t`)
- **THEN** it MUST fetch and install the binary matching that specific release tag.
