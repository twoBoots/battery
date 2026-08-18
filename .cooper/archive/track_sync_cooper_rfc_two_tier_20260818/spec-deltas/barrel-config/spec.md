# Spec Delta: Barrel Configuration & Project Discovery

### Requirement 4: Native Agent Skills & Two-Tier Planning Availability
- The system SHALL provide packaged project-local Cooper skills in `.agents/skills/cooper-{setup,new-track,implement,review,status}`.
+ The system SHALL provide the complete suite of 6 packaged project-local Cooper skills in `.agents/skills/cooper-{setup,rfc,new-track,implement,review,status}` supporting the Two-Tier SDD planning architecture (collaborative upstream RFCs and downstream tactical TDD tracks).

#### Scenario 4.1: Project-Local Skills Discovery & RFC Support
- GIVEN an AI coding assistant operating in the repository
- WHEN discovering project capabilities
+ GIVEN an AI coding assistant operating in the repository
+ WHEN discovering project capabilities
+ THEN skills in `.agents/skills/cooper-{setup,rfc,new-track,implement,review,status}` MUST be accessible and self-contained, enabling both `cooper-rfc` collaborative architectural workflows and `cooper-new-track` tactical implementation tracks.
