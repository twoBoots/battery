# Track Proposal: Multi-Barrel Track Dispatch & Decoupled Planning

## Context & Motivation
Battery coordinates feature epics across multiple independent repositories or monorepo packages (barrels). When proposing cross-cutting features spanning multiple barrels (e.g. `folder-a` backend and `folder-b` frontend), attempting to centrally generate fine-grained TDD implementation plans from the orchestrator causes context window bloat, stale codebase assumptions, and tech stack collisions.

## Objectives
Implement the `battery track` command suite and underlying dispatch engine to support **Contract-First Decoupled Planning**:
1. **Track Dispatch**: Battery creates/dispatches macro specifications, interface contracts, and initial Living Spec Deltas (`spec-deltas/`) to target barrels under matching track IDs.
2. **Local Planning Autonomy**: Barrels receive track specs and their local agents or developers autonomously author local technical designs (`design.md`) and TDD task plans (`plan.md`).
3. **Multi-Barrel Status Aggregation**: Battery aggregates track progress across all participating barrels (worktree status, active tasks, phase checkpoints).

## Value & Impact
* Decouples macro system contracts ("What" and "Why") from localized implementation details ("How").
* Optimizes AI agent context windows for precision and high test coverage (>80%).
* Enables parallel, asynchronous track execution across heterogeneous tech stacks.
