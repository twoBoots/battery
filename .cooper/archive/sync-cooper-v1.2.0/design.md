# Technical Design: Synchronize Cooper v1.2.0 Standards & Skills

## Architecture & Component Breakdown

### 1. Project-Local Agent Skills (`.agents/skills/`)
- **`cooper-implement/SKILL.md`**:
  - Replace Section 3.4 Git Notes verification note command with the audit record template that mandates recording actual command, results, measured coverage %, and manual verification response.
  - Remove the spurious `69: ` line prefix artifact on line 69.
- **`cooper-rfc/SKILL.md`**:
  - Update agent identity to "Cooper System Architect".
  - Insert "The Two-Tiered SDD Architecture" section detailing upstream alignment vs downstream track execution.
  - Insert "Silent Scope Check" in Section 2, enabling agents to offer fast-tracking via `cooper-new-track` for isolated bugs/features.
  - Fix heading regression `## 6.2` to `### 6.2`.

### 2. Embedded Framework Templates (`internal/framework/templates/`)
- Maintain byte-for-byte fidelity with the updated canonical Cooper skills in:
  - `internal/framework/templates/skills/cooper-implement/SKILL.md`
  - `internal/framework/templates/skills/cooper-rfc/SKILL.md`
- Synchronize `internal/framework/templates/docs/COOPER.md` (and `.cooper/COOPER.md`) to standardize Troop links: `[Troop's](https://github.com/twoBoots/troop)` and `| **[Troop](https://github.com/twoBoots/troop)** |`.

### 3. Repository Guidelines & Templates
- **`AGENTS.template.md`**:
  - Incorporate missing Rule 6 ("Interaction & Native Tool Protocols") defining interactive question tools and native file tools mandates.
  - Standardize `[Troop](https://github.com/twoBoots/troop)` link in headings.
- **`README.md`**:
  - Replace `file:///.batteryrc` with `.batteryrc`.

### 4. Specification Promotion
- Update living capability specification `.cooper/specs/documentation/spec.md` with:
  - Requirement: Prohibition on Unverified Checkpoint Attestations
  - Requirement: Two-Tiered SDD Architecture & Collaborative RFCs

### 5. Verification
- Run `go test ./...` across all packages including `internal/framework` and `internal/mcp`.
- Ensure tests verify `framework.InspectFrameworkStatus` detects zero divergence.
