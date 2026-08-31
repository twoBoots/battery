# Technical Design: Synchronize Cooper v1.1.0 Updates

## Context & Architecture
Battery utilizes Cooper SDD patterns for tracking development workflows and embeds Cooper scaffold templates in `internal/framework/templates/`.

Cooper v1.1.0 established two primary rules for agent interactions and execution:
1. **Interactive Question Protocol**: Direct tool invocation of `ask_question` whenever presenting options or decisions to the user.
2. **Native File Tools Mandate**: Exclusively using `view_file`, `write_to_file`, and `replace_file_content` instead of shell editing tools (`sed`, `awk`, heredocs, redirects).

## Component Updates

1. **Project Skills (`.agents/skills/`)**:
   - `cooper-implement`: Add interactive question protocol & native file tools mandate.
   - `cooper-new-track`: Add interactive question protocol & native file tools mandate.
   - `cooper-review`: Add interactive question protocol & native file tools mandate.
   - `cooper-rfc`: Add interactive question protocol & native file tools mandate.
   - `cooper-setup`: Add interactive question protocol & native file tools mandate.
   - `cooper-status`: Add interactive question protocol & native file tools mandate.

2. **Core Guidelines & Workflow**:
   - `AGENTS.md`: Add Rule 6 "Interaction & Native Tool Protocols", update Troop links, retain multi-barrel orchestrator header rules.
   - `.cooper/definition/workflow.md`: Add Guiding Principles 10 & 11, update Troop links, update phase verification and version bump instructions.

3. **Embedded Framework Templates (`internal/framework/templates/`)**:
   - Sync all corresponding skills under `internal/framework/templates/skills/`.
   - Sync `internal/framework/templates/definition/workflow.md`.

4. **Living Specs**:
   - Introduce `.cooper/specs/documentation/spec.md` reflecting the new Cooper v1.1.0 documentation requirements.
