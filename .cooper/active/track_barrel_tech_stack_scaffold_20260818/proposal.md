# Track Proposal: Per-Barrel Tech Stack Scaffolding & Auto-Inference Engine (CLI + MCP)

## Context & Motivation
Battery coordinates multi-repository and monorepo workspaces consisting of heterogeneous packages and services (barrels). Each barrel possesses its own programming language, framework, test runner, coverage thresholds, and linting rules, defined within its `.cooper/definition/tech-stack.md`.

When developers or autonomous AI agents add new barrels or work within a polyglot monorepo (e.g. adding `apps/api` in Go, `apps/web` in Next.js/TypeScript, or `packages/worker` in Python), there is currently no automated scaffolding mechanism to generate standardized Cooper tech stack contracts. AI agents must either manually create markdown files or guess tooling commands.

## Objectives
1. **Core Tech Stack Scaffolding & Auto-Inference (`internal/techstack`)**:
   - Inspect package markers (`package.json`, `go.mod`, `Cargo.toml`, `pyproject.toml`, `requirements.txt`, `deno.json`, `pom.xml`, etc.) to auto-detect language, frameworks, test commands, linters, and coverage parameters.
   - Implement `ScaffoldBarrelTechStack(barrelPath, opts)` to write standardized `.cooper/definition/tech-stack.md` and default code styleguides (`.cooper/code_styleguides/`).
2. **CLI Subcommand (`battery barrel init <path|name>`)**:
   - Provide interactive and flag-driven CLI commands to scaffold a barrel's Cooper tech stack with options for `--language`, `--framework`, `--test-runner`, `--linter`, `--coverage-threshold`, and `--force`.
3. **MCP Tool (`battery_init_barrel_tech_stack`)**:
   - Expose programmatic tech stack scaffolding to AI agents over the Model Context Protocol, enabling agents in monorepos to initialize package-level tech stacks on the fly.

## Value & Impact
* **Eliminates AI Hallucinations**: Gives AI agents exact test and lint commands per package in polyglot monorepos.
* **Standardized Governance**: Ensures all barrels adhere to consistent Cooper SDD contracts.
* **Frictionless Onboarding**: Automates new service and package initialization in seconds.
