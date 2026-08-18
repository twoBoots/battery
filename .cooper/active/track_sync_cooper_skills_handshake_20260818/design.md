# Technical Design: Sync Cooper Skills, Handshake Index & Installer Refinement

## Architecture & Integration

### 1. Project-Local Cooper Skills (`.agents/skills/`)
Package the 5 native Cooper skills directly under `.agents/skills/`:
- `cooper-setup/SKILL.md`
- `cooper-new-track/SKILL.md`
- `cooper-implement/SKILL.md`
- `cooper-review/SKILL.md`
- `cooper-status/SKILL.md`

### 2. Handshake Index (`.cooper/index.md`) & Tracks Registry (`.cooper/tracks.md`)
- `.cooper/index.md`: Maps all project context, including Battery Architecture (`BATTERY.md`), baseline definitions (`definition/`), living specs (`specs/`), tracks (`tracks.md`, `active/`, `archive/`), and capabilities (`.agents/skills/`).
- `.cooper/tracks.md`: Registers all active and completed tracks.

### 3. Battery Installer Protocol (`install.sh`)
- Step 1 delegates to `cooper/install.sh` (or remote raw script).
- Removes redundant directory creation / `TROOP.md` relocation from `battery/install.sh` since Cooper's installer natively handles them.
- Ensures `.cooper/index.md` created by Cooper has the Battery Architecture (`BATTERY.md`) linked.
- Installs `.cooper/BATTERY.md`, configures `AGENTS.md`, and runs CLI binary setup and `battery init`.

### 4. Agent Guidelines (`AGENTS.md` & `AGENTS.template.md`)
- Update Cooper section to include Rule 5: Project-Local Skills (`.agents/skills/`).
- Ensure Battery Agent Rules (Rules 1-4) precede Cooper guidelines.
