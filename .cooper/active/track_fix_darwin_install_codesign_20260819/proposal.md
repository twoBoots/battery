# Proposal: Fix macOS `Killed: 9` Error on Battery Installation

## Problem Statement
When running `install.sh` on macOS (Darwin arm64 / x86_64), the installer encounters:
`bash: line 246: 77656 Killed: 9 "$BATTERY_BIN" init "${INIT_ARGS[@]}" < /dev/tty`

### Root Cause
1. **In-place Truncation**: When `curl -o "${INSTALL_BIN_DIR}/battery"` downloads directly over an existing binary on disk, it truncates the inode in-place. On Darwin/macOS, the kernel's code-signing subsystem (`CS_VALID` / `CS_KILL`) retains cached page/signature hashes for existing vnodes. Modifying the executable in-place invalidates the kernel's code-signature cache, resulting in the kernel immediately terminating the process with `SIGKILL` (`Killed: 9`) upon `execve()`.
2. **Missing Quarantine Removal & Ad-Hoc Signature**: Downloaded binaries via `curl`/`wget` on macOS can be flagged by Gatekeeper or lack clean ad-hoc signatures for Apple Silicon execution.

## Proposed Solution
1. **Atomic Binary Replacement**:
   - Download or build to a temporary file (`${INSTALL_BIN_DIR}/.battery.tmp.$$`).
   - On macOS (`[ "$OS" = "darwin" ]`), strip quarantine attribute (`xattr -d com.apple.quarantine`) and apply/refresh ad-hoc code signature (`codesign -s - --force`).
   - Atomically replace the destination binary using `mv -f` (which assigns a new inode).
2. **Resilient TTY Fallback**:
   - Ensure `init` execution falls back cleanly if `/dev/tty` is closed or piped.
