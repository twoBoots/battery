# Proposal: Synchronize Cooper v1.2.0 Standards & Skills

## Rationale
The upstream Cooper framework was updated to v1.2.0 ([twoBoots/cooper#21](https://github.com/twoBoots/cooper/pull/21)), introducing several critical improvements:
1. **Prohibition on Unverified Checkpoint Attestations**: Prohibiting pre-filled Git Note verification templates (`Automated Tests: PASSED`, `Manual Verification: APPROVED by user`) and requiring explicit audit notes recording actual test command, pass/fail counts, measured coverage %, and explicit user responses.
2. **Two-Tiered SDD Architecture & Scoping**: Clarifying upstream architectural alignment (`cooper-rfc`) vs downstream execution (`cooper-new-track` & `workflow.md`), adding a silent scope check to fast-track isolated tasks, and fixing heading hierarchy (`### 6.2`).
3. **Validate-First Architecture & Link Hygiene**: Enforcing backticked path validation (`link/code-path-exists`) across documentation and templates, fixing invalid URI links in `README.md`, and keeping embedded templates strictly aligned with upstream.
4. **Tool Protocols Parity**: Ensuring `AGENTS.template.md` maintains parity with `AGENTS.md` regarding Rule 6 (Interactive Question Tools and Native File Tools Mandate).

As both an active consumer of Cooper and a multi-barrel orchestrator embedding canonical Cooper templates, Battery must synchronize its local project skills, embedded templates, and repository guides to remain 100% compliant with Cooper v1.2.0.

## Scope Boundaries
- Update project-local skills in `.agents/skills/cooper-implement/` and `.agents/skills/cooper-rfc/`.
- Update embedded framework templates in `internal/framework/templates/skills/` and `internal/framework/templates/docs/COOPER.md`.
- Synchronize root [.cooper/COOPER.md](.cooper/COOPER.md), [AGENTS.template.md](AGENTS.template.md), and fix [README.md](README.md) link reference.
- Promote new requirements to living capability spec [.cooper/specs/documentation/spec.md](.cooper/specs/documentation/spec.md).
- Remove deprecated `"cooper"` entry from `~/.gemini/config/mcp_config.json`.
