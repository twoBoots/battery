# Technical Design: Per-Barrel Tech Stack Scaffolding & Auto-Inference Engine (CLI + MCP)

## Architecture Overview

The scaffolding engine resides in `internal/techstack/scaffold.go`, exposing detection and generation APIs to both the Cobra CLI (`cmd/barrel.go`) and the MCP Server (`internal/mcp/tools.go`).

```
cmd/
└── barrel.go                  # Adds 'battery barrel init <path|name>' subcommand

internal/
├── techstack/
│   ├── techstack.go           # Existing resolution & markdown summary
│   ├── scaffold.go            # New marker detection, inference & Cooper tech-stack generator
│   └── scaffold_test.go       # Comprehensive test suite (>80% coverage)
└── mcp/
    ├── tools.go               # Adds 'battery_init_barrel_tech_stack' tool definition & handler
    └── tools_test.go          # MCP tool tests
```

## Data Structures & Options

```go
type InferredTechStack struct {
    Language          string `json:"language"`
    Framework         string `json:"framework,omitempty"`
    TestRunner        string `json:"test_runner"`
    Linter            string `json:"linter,omitempty"`
    CoverageThreshold string `json:"coverage_threshold"`
    BuildTool         string `json:"build_tool,omitempty"`
}

type ScaffoldOptions struct {
    Language          string `json:"language,omitempty"`
    Framework         string `json:"framework,omitempty"`
    TestRunner        string `json:"test_runner,omitempty"`
    Linter            string `json:"linter,omitempty"`
    CoverageThreshold string `json:"coverage_threshold,omitempty"`
    Force             bool   `json:"force,omitempty"`
}
```

## Marker Detection Matrix

| Project Marker | Inferred Language | Default Test Runner | Default Linter | Coverage Threshold |
| :--- | :--- | :--- | :--- | :--- |
| `go.mod` | `Go` | `go test -v -cover ./...` | `golangci-lint run` | `80%` |
| `package.json` | `TypeScript / JavaScript` | `npm test` (or `pnpm test` / `bun test`) | `npm run lint` / `eslint .` | `80%` |
| `Cargo.toml` | `Rust` | `cargo test` | `cargo clippy` | `80%` |
| `pyproject.toml` / `requirements.txt` | `Python` | `pytest --cov` | `ruff check` / `flake8` | `80%` |
| `deno.json` / `deno.jsonc` | `TypeScript (Deno)` | `deno test` | `deno lint` | `80%` |
| `pom.xml` / `build.gradle` | `Java / Kotlin` | `mvn test` / `./gradlew test` | `mvn checkstyle:check` | `80%` |

## Generated Artifacts

When `ScaffoldBarrelTechStack(barrelPath, opts)` is executed:
1. Creates `<barrelPath>/.cooper/definition/tech-stack.md` formatted with standard Cooper Tech Stack markdown table and sections.
2. Creates `<barrelPath>/.cooper/code_styleguides/<language>.md` with default conventions if not already present.
3. Overwrites existing definitions only if `Force` is true.
