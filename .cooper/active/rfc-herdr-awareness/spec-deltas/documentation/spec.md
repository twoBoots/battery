# Spec Delta: documentation (rfc-herdr-awareness)

Capability: `.cooper/specs/documentation/spec.md`

```diff
+ ### Requirement: Multiplexer-Assisted Dispatch Reference
+ The documentation site MUST describe Battery's dispatch and optional multiplexer behaviour completely enough for an operator without Herdr to use dispatch, and for an operator with Herdr to use spawn and live status.
+
+ #### Scenario: Dispatch and Spawn Reference
+ - **GIVEN** the documentation site
+ - **WHEN** a reader looks up `battery track dispatch`
+ - **THEN** it MUST document `battery track init --map`, `--force`, `--create-worktree`, `--spawn`, `--create-workspace`, `--wait-gates`, `--gate-timeout`, and `--agent-kind`
+ - **AND** the per-delta states (`up-to-date`, `added`, `stale`, `updated`, `orphaned`)
+ - **AND** the per-barrel outcomes (`stale`, `prompted`, `waiting-for-operator`, `operator-active`, `already-running`, `unmapped`, `failed`, including the `agent_not_detected` and `agent_exited` codes)
+ - **AND** the exit-status table (0, 1, 3).
+
+ #### Scenario: Environment Variables
+ - **GIVEN** the documentation site
+ - **WHEN** a reader looks up configuration
+ - **THEN** it MUST document `BATTERY_HERDR` and `BATTERY_AGENT_KIND`, and state that Herdr support is detected only for `dispatch --spawn`, `status --live`, and `battery status`.
+
+ #### Scenario: Operator Authority Statement
+ - **GIVEN** the documentation of `--spawn`
+ - **WHEN** it describes agents stopped at a startup dialog
+ - **THEN** it MUST state that Battery never answers dialogs or approvals and that the operator sends the printed kickoff text after clearing one.
```
