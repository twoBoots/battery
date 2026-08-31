# Proposal: Synchronize Cooper v1.1.0 Standards & Skills

## Rationale
The upstream Cooper framework was updated to v1.1.0 ([twoBoots/cooper#17](https://github.com/twoBoots/cooper/pull/17)), introducing:
1. **Interactive Question Protocol**: Mandating agents use interactive question tools (`ask_question`) for single/multi-choice selections, verifications, and approvals instead of text lists in chat.
2. **Native File Tools Mandate**: Requiring agents to use native file tools (`view_file`, `write_to_file`, `replace_file_content`) while prohibiting `sed`, `awk`, heredocs, and stream redirections.
3. **Troop Markdown Links**: Linking primary mentions of Troop to `https://github.com/twoBoots/troop`.

Battery serves as both a Cooper consumer and a multi-barrel orchestrator embedding framework templates. It must synchronize its local project skills, root workflow/agent rules, and embedded framework templates to remain 100% compliant with Cooper v1.1.0.

## Scope Boundaries
- Update project-local skills in `.agents/skills/`.
- Update repository guidelines in `AGENTS.md` and `.cooper/definition/workflow.md` (retaining Battery multi-barrel rules).
- Update embedded framework templates in `internal/framework/templates/`.
- Introduce and promote the `documentation` living capability specification under `.cooper/specs/documentation/spec.md`.
