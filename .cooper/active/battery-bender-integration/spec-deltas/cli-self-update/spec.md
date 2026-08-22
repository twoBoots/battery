# Spec Delta: CLI Self-Updater (Bender Integration)

## Capability: `cli-self-update`

### Requirement 1: Update Discovery and Check
~ The CLI MUST delegate update discovery and version checking to `bender/pkg/updater`.

#### Scenario 1.1: Checking for updates when up to date
- **GIVEN** a Battery CLI running version `1.2.0`
- **WHEN** running `battery update --check` (or `-c`)
- **THEN** it MUST report that Battery is already up to date with the latest release via Bender updater engine.

#### Scenario 1.2: Checking for updates when an update is available
- **GIVEN** a Battery CLI running version `1.0.0` and latest release `1.2.0`
- **WHEN** running `battery update --check`
- **THEN** it MUST report that an update is available without modifying the local executable.

### Requirement 2: Binary Self-Update Execution
~ The CLI MUST delegate binary download, atomic in-place replacement, and OS permissions management to `bender/pkg/updater`.

#### Scenario 2.1: Applying an update
- **GIVEN** an installed Battery binary and available release
- **WHEN** running `battery update`
- **THEN** it MUST perform atomic binary replacement via Bender's updater engine.
