# Capability Specification Delta: Framework Documentation & Guidelines

## Capability: documentation

## Requirements

### Requirement: Prohibition on Unverified Checkpoint Attestations
+ Cooper and Battery agent skills and workflow instructions MUST NOT provide or encourage pre-filled verification attestations, requiring agents to record actual test outcomes, coverage metrics, and user confirmation.

#### Scenario: Audit Verification Notes
+ - GIVEN an agent executing a phase checkpoint in a Cooper track
+ - WHEN the agent drafts and attaches a Git Note verification record
+ - THEN the note body MUST reflect the actual test command executed, the actual test suite result, and the measured coverage percentage
+ - AND the note body MUST record the explicit verification steps presented and the user's actual recorded response
+ - AND the agent instruction template MUST NOT supply pre-filled assertions such as `Automated Tests: PASSED` or `Manual Verification: APPROVED by user`.

### Requirement: Two-Tiered SDD Architecture & Collaborative RFCs
+ The framework documentation and skills SHALL explicitly distinguish upstream architectural alignment from downstream track execution.

#### Scenario: Upstream Architectural Alignment via RFC
+ - GIVEN an epic, major refactor, or multi-capability initiative
+ - WHEN the initiative requires consensus, trade-off analysis, or spec delta approval before code implementation
+ - THEN the initiative MUST be planned upstream via `cooper-rfc` inside an isolated worktree
+ - AND a Draft Pull Request SHALL be published for team review and comment resolution before implementation tracks are created.

#### Scenario: Silent Scope Assessment
+ - GIVEN an agent invoked with `cooper-rfc`
+ - WHEN the requested change is an isolated single-capability bug fix or localized improvement
+ - THEN the agent SHALL silently assess scope and prompt the user with an option to fast-track directly into `cooper-new-track` without RFC ceremony overhead.
