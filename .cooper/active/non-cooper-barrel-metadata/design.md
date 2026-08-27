# Technical Design: Lightweight Documentation & Context Metadata for Non-Cooper Barrels

## 1. Architectural Overview

This design enhances Battery's configuration, tech stack resolution engine, CLI commands, and MCP server to seamlessly support non-Cooper repositories with zero modifications to target barrels.

```
Battery Root (Orchestrator)
├── .batteryrc / .batteryrc.local
│   └── barrels: [ { name, path, role, tech, docs, jira, ...custom } ]
├── docs/barrels/<barrel-name>.md (Canonical Profile Directory)
│   └── (Fallback: .cooper/barrels/<barrel-name>.md or custom docs path)
└── internal/
    ├── config/          -> Preserves metadata & dynamic attributes
    ├── techstack/       -> Hybrid fallback resolver & profile parser
    └── mcp/             -> Exposes battery://barrels/{name}/docs
```

## 2. Configuration & Metadata Model (`internal/config`)

### 2.1 `BarrelConfig` Enhancements
```go
type BarrelConfig struct {
    Name  string     `json:"name"`
    Path  string     `json:"path"`
    Type  BarrelType `json:"type,omitempty"`
    Role  string     `json:"role,omitempty"`
    Tech  string     `json:"tech,omitempty"`
    Docs  string     `json:"docs,omitempty"`
    Jira  string     `json:"jira,omitempty"`
    
    // Dynamic attributes for custom extensions (e.g. CI links, owner teams)
    Extra map[string]json.RawMessage `json:"-"`
}
```
- Custom JSON Marshaling / Unmarshaling ensures unknown attributes in `.batteryrc` or `.batteryrc.local` are preserved round-trip across `AddBarrel` and `RemoveBarrel` operations.

### 2.2 `EffectiveBarrel` Enhancements
```go
type EffectiveBarrel struct {
    Name            string     `json:"name"`
    Path            string     `json:"path"`
    Type            BarrelType `json:"type,omitempty"`
    Source          string     `json:"source"` // "canonical" or "local"
    Role            string     `json:"role,omitempty"`
    Tech            string     `json:"tech,omitempty"`
    Docs            string     `json:"docs,omitempty"`
    Jira            string     `json:"jira,omitempty"`
    Exists          bool       `json:"exists,omitempty"`
    AbsolutePath    string     `json:"absolutePath,omitempty"`
    CooperTechStack string     `json:"cooperTechStack,omitempty"`
    ProfilePath     string     `json:"profilePath,omitempty"`
    HasProfile      bool       `json:"hasProfile,omitempty"`
    Extra           map[string]json.RawMessage `json:"extra,omitempty"`
}
```

## 3. Context & Tech Stack Fallback Engine (`internal/techstack`)

### 3.1 Profile Resolution Hierarchy
When querying documentation or architectural context for a barrel:
1. Custom path specified in `.batteryrc` `docs` field (resolved relative to Battery root).
2. Standard profile location: `docs/barrels/<name>.md`.
3. Alternative profile location: `.cooper/barrels/<name>.md`.

### 3.2 Hybrid Context Resolution
`ResolveBarrelContext(batteryCwd string, barrel config.EffectiveBarrel)`:
- If barrel has a Cooper `tech-stack.md`:
  - `Summary`: Summarized from Cooper spec.
  - `Source`: `"cooper"`.
- Else if `.batteryrc` has `Tech` or `Role`:
  - `Summary`: `<Tech> (<Role>)` or whichever is defined.
  - `Source`: `"metadata"`.
- Else if barrel profile markdown exists (`docs/barrels/<name>.md`):
  - `Summary`: Extracted highlight summary from the profile markdown.
  - `Source`: `"profile"`.
- Else:
  - `Summary`: `"No Cooper tech-stack or profile defined"`.
  - `Source`: `"none"`.

### 3.3 Profile Starter Scaffolding
A standardized profile template scaffolded at `docs/barrels/<name>.md`:
```markdown
# Barrel Profile: <Barrel Name>

## Role & Responsibilities
<Description of domain responsibility and purpose>

## Tech Stack & Runtime
- Language / Runtime: <e.g. Rust, Python 3.12, C++20, Docker>
- Build System: <e.g. Cargo, CMake, Poetry, Makefile>
- Key Dependencies: <e.g. ROS 2, PyTorch, CUDA>

## Development & Build Commands
- Build: `<command>`
- Test: `<command>`
- Lint: `<command>`

## Interface Contracts & Integration
- Interfaces: <APIs, Protobuf, gRPC, REST, IPC>
- Upstream / Downstream Barrels: <related barrels>

## AI Agent Guidelines
- <Notes for AI coding assistants regarding file structures, constraints, and conventions>
```

## 4. CLI Enhancements (`cmd/`)

1. **`battery barrel add <path>`**:
   - New flags: `--role`, `--tech`, `--docs`, `--jira`.
   - Populates metadata in `.batteryrc` or `.batteryrc.local`.
2. **`battery barrel doc init <name>` / `battery barrel profile init <name>`**:
   - Scaffolds `docs/barrels/<name>.md` with the starter template if it does not already exist (or overwrite with `--force`).
3. **`battery status` & `battery barrel list`**:
   - Displays role, tech, doc profile indicator, and resolved hybrid summary cleanly in console output.

## 5. MCP Server Integrations (`internal/mcp/`)

1. **Resource: `battery://barrels/{name}/docs`**:
   - Returns MIME `text/markdown` containing the resolved profile document.
2. **Resource: `battery://barrels/{name}/tech-stack`**:
   - If Cooper `tech-stack.md` is present, returns Cooper spec.
   - If absent, returns profile markdown content or structured fallback context with metadata.
3. **Resource: `battery://topology` & Tool `battery_list_barrels`**:
   - Enriched JSON payloads containing `role`, `tech`, `docs`, `jira`, `hasProfile`, and `profilePath`.
