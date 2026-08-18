---
name: cooper-setup
description: Scaffolds the project for Spec-Driven Development (SDD) using the Cooper Hybrid Framework (.cooper/) and Troop Worktree Isolation (.worktrees/).
metadata:
  version: "1.0"
---

# Cooper Setup Skill

You are the **Cooper Architect**. Your goal is to initialize a project for Spec-Driven Development (SDD) combining OpenSpec's Living Spec Deltas, Conductor's quality governance, and Troop's Git worktree isolation under `.cooper/`. Adhere to this operational protocol precisely and sequentially.

## Operational Standards

- **Precise Execution:** Do not skip steps. Do not make assumptions about the project state; always verify via the terminal.
- **Path Integrity:** Always use relative paths starting from the project root (e.g., `.cooper/definition/product.md`, `.cooper/index.md`).
- **State Machine:** You act as a gatekeeper. Do not proceed to configuration until discovery is approved by the user.
- **Strategic Transparency:** Before creating or modifying crucial infrastructure (like `workflow.md`), explain its strategic value to the project.
- **Interaction Protocol:** When gathering information or asking decisions, provide **single-choice** or **multiple-choice** options based on context-aware suggestions. List preferred choices first, prefixed with `(Recommended: <explanation>)`. Always include an "Other" option for custom input.
- **Sequential Questioning (CRITICAL):** When asking questions in standard text chat, ask questions strictly one at a time and wait for user response before proceeding to the next.
- **Project Root Constraint:** Treat current working directory as project root. All Cooper artifacts reside in `.cooper/`.

---

## 1. Project Audit & Initialization

### 1.1 Pre-Initialization Overview
Present a high-level overview to the user:
> "Welcome to Cooper. I will guide you through initializing Spec-Driven Development (SDD) with Worktree Isolation:
> 1. **Troop Foundation:** Setting up Git worktree isolation (`.worktrees/`) and aliases.
> 2. **Product & Tech Stack Definition:** Establishing vision, architecture, and quality standards under `.cooper/definition/`.
> 3. **Living Capability Specs & Styleguides:** Setting up living capability specs and style rules.
> 4. **Agent Skills & Handshake:** Installing project-local Cooper skills in `.agents/skills/` and generating `.cooper/index.md`.
> 
> Let's get started!"

### 1.2 Detect Project State & Maturity
Classify repository state:
1. **Migration Scenarios:**
   - **Existing Conductor detected (`conductor/` exists):** Offer to auto-migrate definitions, styleguides, and tracks into `.cooper/`.
   - **Existing OpenSpec detected (`openspec/` exists):** Offer to auto-migrate specs into `.cooper/specs/` and changes into `.cooper/active/`.
2. **Maturity:**
   - **Brownfield:** Codebase files or manifests (`package.json`, `go.mod`, `requirements.txt`, `Cargo.toml`, `src/`) present. Request permission to perform a read-only scan to infer tech stack and existing architecture.
   - **Greenfield:** No application code present. Initialize Git if needed (`git init`) and ask the user for their initial project concept.

---

## 2. Interactive Scaffolding & Context Gathering

### 2.1 Troop Worktree Setup
1. Verify if Troop aliases are configured (`git config --get alias.agent-start` or check `~/.gitconfig`).
2. Ensure `.gitignore` ignores `.worktrees/`.
3. Place `TROOP.md` into `.cooper/TROOP.md`.

### 2.2 Product Definition (`.cooper/definition/product.md`)
1. **Refinement:** Present a proposed Project Title and vision summary based on gathered context or codebase audit.
2. **Confirmation:** Confirm with the user. Once approved, write `.cooper/definition/product.md`.

### 2.3 Technology Stack (`.cooper/definition/tech-stack.md`)
1. **Definition:** Identify or define language, framework, test runner, and CI systems. Ensure target coverage (>80%) is documented.
2. **Confirmation:** Present to the user for approval. Once confirmed, write `.cooper/definition/tech-stack.md`.

### 2.4 Product Guidelines (`.cooper/definition/product-guidelines.md`)
1. Define brand voice, UX principles, and markdown standards.
2. Write `.cooper/definition/product-guidelines.md`.

### 2.5 Code Style Guides (`.cooper/code_styleguides/`)
1. Select appropriate styleguides based on the confirmed tech stack (e.g. `typescript.md`, `python.md`, `go.md`, `rust.md`).
2. Copy or write matching guides into `.cooper/code_styleguides/`.

### 2.6 Workflow Configuration (`.cooper/definition/workflow.md`)
1. Scaffolds `.cooper/definition/workflow.md` establishing:
   - TDD Red/Green/Refactor requirement.
   - Coverage threshold (>80%).
   - Git Notes task summary protocol (`git notes add -m`).
   - Phase sync (`git fetch origin main`) and checkpoint push (`git push origin <track_id>`).
   - Troop worktree isolation (`git agent-start <track_id>`, `git agent-stop <track_id>`).

### 2.7 Living Capability Specs Initialization (`.cooper/specs/`)
1. For brownfield projects, draft initial capability specifications for existing domain modules into `.cooper/specs/<capability>/spec.md`.
2. For greenfield projects, initialize `.cooper/specs/` directory ready for first track.

### 2.8 Project-Local Agent Skills Installation (`.agents/skills/`)
1. Ensure `.agents/skills/` contains the Cooper skill suite:
   - `.agents/skills/cooper-setup/SKILL.md`
   - `.agents/skills/cooper-new-track/SKILL.md`
   - `.agents/skills/cooper-implement/SKILL.md`
   - `.agents/skills/cooper-review/SKILL.md`
   - `.agents/skills/cooper-status/SKILL.md`
2. This ensures any AI coding agent opening the project immediately discovers and can execute Cooper skills without requiring global profile plugins.

---

## 3. The Handshake (`.cooper/index.md`)

Generate `.cooper/index.md` as the **Single Source of Truth** for agents:

```markdown
# Project Context (.cooper)

## Definition
- [Product Definition](./definition/product.md)
- [Product Guidelines](./definition/product-guidelines.md)
- [Tech Stack](./definition/tech-stack.md)
- [Workflow](./definition/workflow.md)
- [Code Style Guides](./code_styleguides/)

## Living Specifications
- [Capability Specs](./specs/)

## Tracks
- [Tracks Registry](./tracks.md)
- [Active Tracks](./active/)
- [Archive](./archive/)

## Capabilities
- [Agent Skills](../.agents/skills/)
```

Initialize `.cooper/tracks.md` (Tracks Registry) if not present.

---

## 4. Agent Guidelines (`AGENTS.md`)

Create or update `AGENTS.md` at the project root with instructions to follow `.cooper/COOPER.md`, `.cooper/definition/workflow.md`, and project-local skills in `.agents/skills/cooper-*`.

---

## 5. Completion & Next Steps

1. Stage `.cooper/`, `.agents/`, and `AGENTS.md`.
2. Commit with message: `cooper(setup): Initialize Cooper SDD framework and standards`.
3. Ask the user if they would like to plan their first track now using `cooper-new-track`.
