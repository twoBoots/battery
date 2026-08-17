# Product Definition & Vision

`battery` is an open, **agent-agnostic** multi-repository Specification-Driven Development (SDD) orchestration pattern for managing a collection of barrels (repositories).

## Overview
`battery` leverages the Cooper Hybrid SDD Framework (`.cooper/`) and Troop worktree isolation (`.worktrees/`) to enable cross-barrel track planning, living capability spec definition, and execution sequencing across arbitrary barrels for human developers and autonomous AI agents.

## Goals
- **Universal Multi-Barrel Planning**: Provide a standardized pattern for tracking epics across independent barrels using Cooper active tracks (`.cooper/active/<track_id>/`).
- **Contract-Driven Specification**: Author and validate cross-barrel capability specs and Spec Deltas (`spec-deltas/`) prior to downstream implementation in individual barrels.
- **Cooper Worktree Alignment**: Coordinate concurrent feature branches (`.worktrees/<track_id>`) across target barrels.
- **Traceable Checkpoints**: Provide verifiable progress tracking across multi-barrel execution phases using TDD, coverage gates (>80%), and Git Notes verification reports.

