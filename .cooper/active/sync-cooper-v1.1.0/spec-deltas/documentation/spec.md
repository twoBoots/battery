# Capability Specification: Framework Documentation & Guidelines

## Purpose & Scope
Defines the standards and guidelines for project-level documentation, framework reference manuals, agent prompts, and templates within the Cooper & Battery ecosystem.

## Requirements

### Requirement: External Dependency Attribution & Repository Linking
All framework documentation, manuals, skill instructions, and templates referencing foundational external tools SHALL link directly to their upstream source repositories on primary introduction.

#### Scenario: Primary Mention of Troop in Markdown Documents
+ - GIVEN a Cooper markdown documentation file (`README.md`, `AGENTS.md`, `docs/*.md`, `.cooper/*.md`, `skills/*/*.md`, or `templates/*.md`)
+ - WHEN the document mentions the Troop worktree isolation tool
+ - THEN the first prominent mention of Troop MUST be formatted as a markdown link pointing to `https://github.com/twoBoots/troop` (e.g. `[Troop](https://github.com/twoBoots/troop)`).

#### Scenario: Tool Summary Tables & Capability Listings
+ - GIVEN a table, overview matrix, or capability list referencing Cooper foundation components
+ - WHEN Troop is listed as a tool or pillar
+ - THEN the entry MUST include a direct markdown link to `https://github.com/twoBoots/troop`.

### Requirement: Interactive Question & File Tool Protocols
When interacting with users or performing file operations during skill execution, agents MUST prioritize runtime interactive tools and native file tools over plain text or shell stream redirections.

#### Scenario: Interactive Question Tool Invocations
+ - GIVEN an agent executing a Cooper skill or workflow step requiring user input, selection, or confirmation
+ - AND the agent runtime provides an interactive question tool (such as `ask_question`)
+ - WHEN the agent presents single-choice or multiple-choice options
+ - THEN the agent MUST invoke the interactive question tool
+ - AND the agent MUST NOT output plain text numbered/bulleted option lists in chat.

#### Scenario: Fallback to Text Chat for Questions
+ - GIVEN an agent executing a Cooper skill or workflow step
+ - AND no interactive question tool is provided in the agent runtime environment
+ - WHEN the agent needs to ask questions or present choices
+ - THEN the agent SHALL format the question clearly in text chat, asking questions strictly one at a time.

#### Scenario: Native File Tools Enforcement
+ - GIVEN an agent runtime providing native file tools (`view_file`, `write_to_file`, `replace_file_content`)
+ - WHEN the agent needs to view, create, or edit files in the workspace
+ - THEN the agent MUST use native file tools
+ - AND the agent MUST NOT use shell commands with stream editors (`sed`, `awk`), heredocs, pipes, or stream redirections (`cat << 'EOF'`, `echo >`, `cat >`) to create or edit files.
