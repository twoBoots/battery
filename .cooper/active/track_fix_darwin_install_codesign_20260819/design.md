# Technical Design: Atomic Binary Installation & Darwin Code Signing

## Architecture & Implementation Details

### 1. Atomic Binary Replacement in `install.sh`
Instead of downloading directly to `${INSTALL_BIN_DIR}/battery`, `install.sh` performs atomic replacement:
1. Generate unique temporary path: `TMP_FILE="${INSTALL_BIN_DIR}/.battery.tmp.$$$"`
2. Download binary via `curl -fsSL "$RELEASE_URL" -o "$TMP_FILE"` (or `wget`).
3. Set executable bit: `chmod +x "$TMP_FILE"`.
4. If running on Darwin (`OS="darwin"`):
   - Remove quarantine attribute: `xattr -d com.apple.quarantine "$TMP_FILE" 2>/dev/null || true`
   - Re-sign with ad-hoc signature: `codesign -s - --force "$TMP_FILE" 2>/dev/null || true`
5. Atomically move into place: `mv -f "$TMP_FILE" "${INSTALL_BIN_DIR}/battery"`.
6. Clean up temporary files on failure.

### 2. Local Build Atomic Replacement
For Tier 1 local build (`go build`):
1. Build to temporary binary in target bin dir.
2. Apply codesign on Darwin if needed.
3. Atomically move into place via `mv -f`.
