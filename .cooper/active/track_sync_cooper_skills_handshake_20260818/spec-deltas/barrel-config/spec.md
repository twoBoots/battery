# Spec Delta: Barrel Configuration & Project Discovery

### Requirement 1: Layered Configuration Schema & Project Handshake
- The system SHALL rely on ad-hoc document discovery without a unified Handshake index.
+ The system SHALL maintain a unified project Handshake index (`.cooper/index.md`) and tracks registry (`.cooper/tracks.md`) that links to Battery architecture, definitions, living specs, tracks, and native skills.

#### Scenario 1.3: Project Handshake Index
+ GIVEN a Battery workspace
+ WHEN inspecting project architecture and capabilities
+ THEN `.cooper/index.md` MUST provide direct links to `.cooper/BATTERY.md`, `.cooper/definition/`, `.cooper/specs/`, `.cooper/tracks.md`, and `.agents/skills/`.

### Requirement 4: Native Agent Skills Availability
+ The system SHALL provide packaged project-local Cooper skills under `.agents/skills/cooper-*` to enable agentic workflows without external plugins.

#### Scenario 4.1: Project-Local Skills Discovery
+ GIVEN an AI coding assistant operating in the repository
+ WHEN discovering project capabilities
+ THEN skills in `.agents/skills/cooper-{setup,new-track,implement,review,status}` MUST be accessible and self-contained.
