# Proposal: Sync Cooper RFC PR Approval Protocols & Reviewer Guidance

## Motivation
Upstream repository `twoBoots/cooper` recently merged Pull Request #7 (`feat(cooper-rfc): Add PR approval detection protocols and reviewer instructions`) and Pull Request #6 (`docs(rfc): Document rationale and benefits of Draft PRs for RFC lifecycle`).

These updates introduce:
1. Structured reviewer guidance templates (`### Reviewer Actions`) in `prbody.md` when opening Draft RFC Pull Requests.
2. Automated PR approval detection via GitHub CLI native review decisions (`reviewDecision == "APPROVED"`) and discussion comment triggers (`/approve`).
3. Streamlined RFC graduation flow (`metadata.json`, track registration in `.cooper/tracks.md`, `gh pr ready`, and human maintainer merge handoff).
4. Synchronized architecture documentation in `.cooper/COOPER.md` and track execution rules in `.cooper/definition/workflow.md`.

Battery, as a multi-barrel SDD orchestrator and Cooper platform, should maintain full synchronization with native Cooper skills, documentation, and framework rules.

## Scope
- Synchronize `.agents/skills/cooper-rfc/SKILL.md` with upstream approval detection and reviewer guidance protocols.
- Synchronize `.cooper/COOPER.md` with the updated Two-Tier Planning Model description.
- Synchronize `.cooper/definition/workflow.md` with Rule 9 clarifications.
- Update living capability spec `.cooper/specs/barrel-config/spec.md` to reflect the enhanced RFC approval protocol.
