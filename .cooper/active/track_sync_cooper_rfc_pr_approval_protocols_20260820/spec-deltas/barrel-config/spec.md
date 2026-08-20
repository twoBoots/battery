# Spec Delta: Barrel Configuration & Project Discovery (RFC PR Approval Protocols)

## Capability: barrel-config

### Requirement 7: Native Agent Skills & Two-Tier Planning Availability

#### Scenario 7.1: Project-Local Skills Discovery & RFC Support
- **GIVEN** an AI coding assistant operating in the repository
- **WHEN** discovering project capabilities
- **THEN** skills in `.agents/skills/cooper-{setup,rfc,new-track,implement,review,status}` MUST be accessible and self-contained, enabling both `cooper-rfc` collaborative architectural workflows and `cooper-new-track` tactical implementation tracks.

+#### Scenario 7.2: Automated RFC Approval Detection & Reviewer Guidance
+- **GIVEN** an active Draft RFC PR created via `cooper-rfc`
+- **WHEN** PR feedback is inspected or approval is evaluated
+- **THEN** `cooper-rfc` MUST support automated detection of approvals via GitHub native review decisions (`reviewDecision == "APPROVED"`) or `/approve` comment triggers, embed reviewer action instructions in PR templates, and transition approved RFCs to ready for merge with registered child tracks.
