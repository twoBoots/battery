# Spec Delta: CI & Release Workflows

## Capability: `ci-release`

### Requirement 3: Cross-Platform Binary Installation Integrity
The installation script MUST atomically replace binary executables and ensure macOS code signing and quarantine compliance to prevent `Killed: 9` execution errors.

#### Scenario 3.1: Atomic Download & Darwin Code Signing
- **GIVEN** a machine installing or upgrading the Battery CLI via `install.sh`
- **WHEN** downloading release binaries or building locally
- **THEN** the installer MUST:
+   - Download or build to a temporary file before moving it in place
+   - On Darwin (macOS), strip quarantine attributes (`com.apple.quarantine`)
+   - On Darwin (macOS), enforce ad-hoc code signature (`codesign -s - --force`)
+   - Atomically replace the destination binary using `mv -f` to preserve kernel inode signature integrity.
