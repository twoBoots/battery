# Spec Delta: Barrel Configuration & Lightweight Profile Metadata

## Living Specification
- Target Capability: `barrel-config`
- Target Spec File: `.cooper/specs/barrel-config/spec.md`

## Spec Diffs

```diff
--- a/.cooper/specs/barrel-config/spec.md
+++ b/.cooper/specs/barrel-config/spec.md
@@ -19,4 +19,10 @@
     - `name`: string identifier for the barrel
     - `path`: relative or absolute directory path to the barrel
     - `type`: optional `"barrel"` (default) or `"battery"` (for hierarchical sub-batteries)
+    - `role`: optional string describing domain responsibility
+    - `tech`: optional string summarizing runtime and tech stack
+    - `docs`: optional path to orchestrator-level documentation profile
+    - `jira`: optional issue tracker / project key mapping
+    - Custom dynamic attributes MUST be preserved round-trip during serialization.

@@ -27,7 +33,7 @@
-### Requirement 2: Reading Barrel Tech Stack from Cooper (`.cooper/definition/tech-stack.md`)
-The system MUST respect each respective barrel's Cooper tech stack definition rather than hardcoding it in `.batteryrc`.
+### Requirement 2: Reading Barrel Tech Stack & Context Metadata
+The system MUST respect each respective barrel's Cooper tech stack definition when present, and seamlessly fall back to `.batteryrc` metadata and orchestrator documentation profiles when managing non-Cooper barrels.

-#### Scenario 2.1: Barrel Tech Stack Resolution
+#### Scenario 2.1: Barrel Tech Stack & Context Fallback Resolution
 - **GIVEN** a registered barrel path and configuration
 - **WHEN** inspecting barrel details or running multi-barrel status/planning
-- **THEN** battery MUST read the barrel's `.cooper/definition/tech-stack.md` (or `conductor/tech-stack.md`) to resolve its language, framework, and tooling specifications.
+- **THEN** battery MUST resolve the tech stack and context using the fallback hierarchy:
+  1. Barrel's `.cooper/definition/tech-stack.md` (or `conductor/tech-stack.md` / `tech-stack.md`)
+  2. `.batteryrc` / `.batteryrc.local` inline `tech` and `role` fields
+  3. Orchestrator-level barrel documentation profile (`docs/barrels/<name>.md` or `.cooper/barrels/<name>.md`)
+  4. Default message: `"No Cooper tech-stack or profile defined"`.

@@ -75,3 +81,3 @@
 #### Scenario 5.1: Add Barrel
 - **GIVEN** a valid directory path
-- **WHEN** `battery barrel add <path> [--name <name>] [--type <barrel|battery>] [--local]` is executed
+- **WHEN** `battery barrel add <path> [--name <name>] [--type <barrel|battery>] [--role <role>] [--tech <tech>] [--docs <docs>] [--jira <jira>] [--local]` is executed
 - **THEN** the barrel path MUST be normalized, name inferred or set, and appended to the target configuration file (`.batteryrc` or `.batteryrc.local`) with any specified metadata attributes without duplicates.

@@ -86,3 +92,3 @@
 #### Scenario 5.3: List Barrels
 - **GIVEN** existing `.batteryrc` and optional `.batteryrc.local`
 - **WHEN** `battery barrel list` or `battery status` is executed
-- **THEN** all effective barrels MUST be displayed with their name, path, source (`canonical` vs `local override`), existence status, and resolved Cooper tech stack summary.
+- **THEN** all effective barrels MUST be displayed with their name, path, source (`canonical` vs `local override`), existence status, metadata attributes (`role`, `tech`, `docs`, `jira`), and resolved hybrid tech stack / profile summary.

+### Requirement 8: Orchestrator-Level Barrel Documentation & Profile Scaffolding (`docs/barrels/`)
+The system MUST support lightweight, zero-intrusion documentation profiles for non-Cooper barrels stored directly in the Battery root repository.
+
+#### Scenario 8.1: Standard Profile Discovery & Custom Docs Resolution
+- **GIVEN** a registered barrel without `.cooper/`
+- **WHEN** context or documentation is queried
+- **THEN** Battery MUST discover profiles located at `docs/barrels/<name>.md`, `.cooper/barrels/<name>.md`, or the custom path specified in `docs`.
+
+#### Scenario 8.2: Profile Scaffolding CLI Subcommand
+- **GIVEN** a barrel name or target identifier
+- **WHEN** `battery barrel doc init <name>` or `battery barrel profile init <name>` is executed
+- **THEN** Battery MUST scaffold a starter Markdown profile at `docs/barrels/<name>.md` containing structured sections for Role & Responsibilities, Tech Stack & Runtime, Development & Build Commands, Interface Contracts, and AI Agent Guidelines without modifying the target repository.
```
