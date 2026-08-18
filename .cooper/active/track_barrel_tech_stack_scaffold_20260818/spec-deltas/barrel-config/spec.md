# Spec Delta: Barrel Configuration & Project Discovery

## Capability: `barrel-config`

### Requirement 6: Per-Barrel Cooper Tech Stack Scaffolding & Auto-Inference
+ The system MUST provide automated discovery, inference, and generation of per-barrel Cooper tech stack definitions (`.cooper/definition/tech-stack.md`) for independent repositories and monorepo sub-packages.
+ 
+ #### Scenario 6.1: Automated Marker-Based Inference
+ - **GIVEN** a barrel directory containing project markers (such as `go.mod`, `package.json`, `Cargo.toml`, or `pyproject.toml`)
+ - **WHEN** tech stack inference is executed
+ - **THEN** the system MUST detect the primary programming language, recommend test runners, linters, and coverage threshold defaults (>80%).
+ 
+ #### Scenario 6.2: Scaffolding `.cooper/definition/tech-stack.md`
+ - **GIVEN** a target barrel path and optional specification overrides (language, framework, test runner, linter)
+ - **WHEN** `battery barrel init <path|name>` CLI command or scaffolding API is invoked
+ - **THEN** it MUST create `<barrel_path>/.cooper/definition/tech-stack.md` formatted to Cooper specifications, protecting existing files unless `--force` is specified.
