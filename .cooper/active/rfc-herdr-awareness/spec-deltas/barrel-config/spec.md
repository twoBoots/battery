# Spec Delta: barrel-config (rfc-herdr-awareness)

Capability: `.cooper/specs/barrel-config/spec.md`

## Requirement 5: Barrel Path Management (Add, Remove, List)

```diff
+ #### Scenario 5.4: Multiplexer Line in `battery status`
+ - **GIVEN** `battery status` is executed
+ - **WHEN** multiplexer detection (track-dispatch Requirement 5) selects the Herdr backend
+ - **THEN** the output MUST include one line naming the multiplexer and its server version, and, for each barrel, whether a source workspace is open for it
+ - **AND** when the no-op backend is selected, the output MUST include one line stating that no multiplexer was detected
+ - **AND** the line MUST NOT include any Herdr workspace, tab, or pane identifier.
```

`battery barrel list` output is unchanged.
