# Implementation Plan: Synchronize Cooper v1.2.0 Updates

## Phase 1: Project-Local Skills & Guidelines Sync
- [x] Task 1.1: Update `.agents/skills/cooper-implement/SKILL.md` (audit attestation template & remove line 69 typo) (90a8abb)
- [x] Task 1.2: Update `.agents/skills/cooper-rfc/SKILL.md` (Two-Tiered SDD, silent scope check, heading 6.2 fix) (582f98c)
- [x] Task 1.3: Synchronize `AGENTS.template.md`, `.cooper/COOPER.md`, and fix `README.md` link (198a546)
- [x] Task 1.4: Phase 1 Verification & Checkpoint [checkpoint: b4bf8aa]

## Phase 2: Embedded Framework Templates Sync & Tests
- [x] Task 2.1: Sync `internal/framework/templates/skills/cooper-implement/SKILL.md` (1b054de)
- [x] Task 2.2: Sync `internal/framework/templates/skills/cooper-rfc/SKILL.md` (c63041c)
- [x] Task 2.3: Sync `internal/framework/templates/docs/COOPER.md` (d96371d)
- [x] Task 2.4: Add guard test in `internal/framework/` ensuring no fabricated attestation templates exist (850ea56)
- [x] Task 2.5: Phase 2 Verification & Checkpoint [checkpoint: aa02757]

## Phase 3: Spec Promotion, MCP Cleanup & Final Verification
- [x] Task 3.1: Promote living capability spec in `.cooper/specs/documentation/spec.md` (192c57f)
- [x] Task 3.2: Remove deprecated `"cooper"` server entry from `~/.gemini/config/mcp_config.json`
- [x] Task 3.3: Final test suite verification, checkpoint & completion report [checkpoint: 050c35f]
