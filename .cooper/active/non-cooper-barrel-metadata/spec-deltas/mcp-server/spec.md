# Spec Delta: Model Context Protocol (MCP) Non-Cooper Documentation Resources

## Living Specification
- Target Capability: `mcp-server`
- Target Spec File: `.cooper/specs/mcp-server/spec.md`

## Spec Diffs

```diff
--- a/.cooper/specs/mcp-server/spec.md
+++ b/.cooper/specs/mcp-server/spec.md
@@ -31,3 +31,3 @@
 #### Scenario 2.2: `battery_list_barrels` Tool
 - **GIVEN** a valid workspace directory
 - **WHEN** `tools/call` is invoked with `name: "battery_list_barrels"`
-- **THEN** it MUST return a list of barrels along with resolved Cooper tech stack summaries.
+- **THEN** it MUST return a list of barrels along with metadata attributes (`role`, `tech`, `docs`, `jira`, `hasProfile`) and resolved hybrid tech stack / profile summaries.

@@ -69,3 +69,3 @@
 #### Scenario 3.1: Listing and Reading Resources
 - **GIVEN** a running MCP session
-- **WHEN** querying `resources/list` or `resources/read` for `battery://topology`, `battery://barrels/{name}/tech-stack`, or `battery://tracks/{track_id}`
+- **WHEN** querying `resources/list` or `resources/read` for `battery://topology`, `battery://barrels/{name}/tech-stack`, `battery://barrels/{name}/docs`, or `battery://tracks/{track_id}`
 - **THEN** the server MUST return valid resource payloads with correct MIME types.

+#### Scenario 3.3: Non-Cooper Barrel Documentation Resource & Fallbacks
+- **GIVEN** an active MCP session and a registered barrel
+- **WHEN** querying `resources/read` for `battery://barrels/{name}/docs`
+- **THEN** the server MUST return the markdown content of the resolved barrel profile document (`docs/barrels/<name>.md`, `.cooper/barrels/<name>.md`, or custom `docs` path) with MIME type `text/markdown`
+- **AND** when querying `battery://barrels/{name}/tech-stack` for a barrel lacking Cooper `tech-stack.md`, the server MUST fall back to returning the profile document content or summarized metadata.
```
