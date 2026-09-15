# Capability Specification: Framework Documentation & Guidelines

## Purpose & Scope
Defines the standards and guidelines for project-level documentation, framework reference manuals, agent prompts, and templates within the Cooper & Battery ecosystem.

## Requirements

### Requirement: External Dependency Attribution & Repository Linking
All framework documentation, manuals, skill instructions, and templates referencing foundational external tools SHALL link directly to their upstream source repositories on primary introduction.

#### Scenario: Primary Mention of Troop in Markdown Documents
- GIVEN a Cooper markdown documentation file (`README.md`, `AGENTS.md`, `docs/*.md`, `.cooper/*.md`, `skills/*/*.md`, or `templates/*.md`)
- WHEN the document mentions the Troop worktree isolation tool
- THEN the first prominent mention of Troop MUST be formatted as a markdown link pointing to `https://github.com/twoBoots/troop` (e.g. `[Troop](https://github.com/twoBoots/troop)`).

#### Scenario: Tool Summary Tables & Capability Listings
- GIVEN a table, overview matrix, or capability list referencing Cooper foundation components
- WHEN Troop is listed as a tool or pillar
- THEN the entry MUST include a direct markdown link to `https://github.com/twoBoots/troop`.

### Requirement: Interactive Question & File Tool Protocols
When interacting with users or performing file operations during skill execution, agents MUST prioritize runtime interactive tools and native file tools over plain text or shell stream redirections.

#### Scenario: Interactive Question Tool Invocations
- GIVEN an agent executing a Cooper skill or workflow step requiring user input, selection, or confirmation
- AND the agent runtime provides an interactive question tool (such as `ask_question`)
- WHEN the agent presents single-choice or multiple-choice options
- THEN the agent MUST invoke the interactive question tool
- AND the agent MUST NOT output plain text numbered/bulleted option lists in chat.

#### Scenario: Fallback to Text Chat for Questions
- GIVEN an agent executing a Cooper skill or workflow step
- AND no interactive question tool is provided in the agent runtime environment
- WHEN the agent needs to ask questions or present choices
- THEN the agent SHALL format the question clearly in text chat, asking questions strictly one at a time.

#### Scenario: Native File Tools Enforcement
- GIVEN an agent runtime providing native file tools (`view_file`, `write_to_file`, `replace_file_content`)
- WHEN the agent needs to view, create, or edit files in the workspace
- THEN the agent MUST use native file tools
- AND the agent MUST NOT use shell commands with stream editors (`sed`, `awk`), heredocs, pipes, or stream redirections (`cat << 'EOF'`, `echo >`, `cat >`) to create or edit files.

### Requirement: Prohibition on Unverified Checkpoint Attestations
Cooper and Battery agent skills and workflow instructions MUST NOT provide or encourage pre-filled verification attestations, requiring agents to record actual test outcomes, coverage metrics, and user confirmation.

#### Scenario: Audit Verification Notes
- GIVEN an agent executing a phase checkpoint in a Cooper track
- WHEN the agent drafts and attaches a Git Note verification record
- THEN the note body MUST reflect the actual test command executed, the actual test suite result, and the measured coverage percentage
- AND the note body MUST record the explicit verification steps presented and the user's actual recorded response
- AND the agent instruction template MUST NOT supply pre-filled assertions such as `Automated Tests: PASSED` or `Manual Verification: APPROVED by user`.

### Requirement: Two-Tiered SDD Architecture & Collaborative RFCs
The framework documentation and skills SHALL explicitly distinguish upstream architectural alignment from downstream track execution.

#### Scenario: Upstream Architectural Alignment via RFC
- GIVEN an epic, major refactor, or multi-capability initiative
- WHEN the initiative requires consensus, trade-off analysis, or spec delta approval before code implementation
- THEN the initiative MUST be planned upstream via `cooper-rfc` inside an isolated worktree
- AND a Draft Pull Request SHALL be published for team review and comment resolution before implementation tracks are created.

#### Scenario: Silent Scope Assessment
- GIVEN an agent invoked with `cooper-rfc`
- WHEN the requested change is an isolated single-capability bug fix or localized improvement
- THEN the agent SHALL silently assess scope and prompt the user with an option to fast-track directly into `cooper-new-track` without RFC ceremony overhead.

