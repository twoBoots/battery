package techstack

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// InferredTechStack holds auto-detected tooling configurations.
type InferredTechStack struct {
	Language          string `json:"language"`
	Framework         string `json:"framework,omitempty"`
	TestRunner        string `json:"test_runner"`
	Linter            string `json:"linter,omitempty"`
	CoverageThreshold string `json:"coverage_threshold"`
	BuildTool         string `json:"build_tool,omitempty"`
}

// ScaffoldOptions provides overrides for scaffolding .cooper/definition/tech-stack.md.
type ScaffoldOptions struct {
	Language          string `json:"language,omitempty"`
	Framework         string `json:"framework,omitempty"`
	TestRunner        string `json:"test_runner,omitempty"`
	Linter            string `json:"linter,omitempty"`
	CoverageThreshold string `json:"coverage_threshold,omitempty"`
	Force             bool   `json:"force,omitempty"`
}

// ScaffoldResult represents the outcome of tech stack scaffolding.
type ScaffoldResult struct {
	BarrelPath     string             `json:"barrel_path"`
	TechStackPath  string             `json:"tech_stack_path"`
	StyleguidePath string             `json:"styleguide_path,omitempty"`
	Inferred       *InferredTechStack `json:"inferred,omitempty"`
	Created        bool               `json:"created"`
	Overwritten    bool               `json:"overwritten"`
}

// InferTechStack inspects project files in a barrel directory to infer tooling.
func InferTechStack(barrelPath string) InferredTechStack {
	// 1. Go
	if fileExists(filepath.Join(barrelPath, "go.mod")) {
		return InferredTechStack{
			Language:          "Go 1.27+",
			TestRunner:        "go test -v -cover ./...",
			Linter:            "golangci-lint run",
			CoverageThreshold: ">80%",
			BuildTool:         "go build",
		}
	}

	// 2. Deno
	if fileExists(filepath.Join(barrelPath, "deno.json")) || fileExists(filepath.Join(barrelPath, "deno.jsonc")) {
		return InferredTechStack{
			Language:          "TypeScript (Deno)",
			TestRunner:        "deno test",
			Linter:            "deno lint",
			CoverageThreshold: ">80%",
			BuildTool:         "deno compile",
		}
	}

	// 3. Node / TypeScript / JavaScript
	if fileExists(filepath.Join(barrelPath, "package.json")) {
		testCmd := "npm test"
		lintCmd := "npm run lint"
		buildCmd := "npm run build"

		if fileExists(filepath.Join(barrelPath, "pnpm-lock.yaml")) {
			testCmd = "pnpm test"
			lintCmd = "pnpm lint"
			buildCmd = "pnpm build"
		} else if fileExists(filepath.Join(barrelPath, "bun.lockb")) || fileExists(filepath.Join(barrelPath, "bun.lock")) {
			testCmd = "bun test"
			lintCmd = "bun lint"
			buildCmd = "bun build"
		} else if fileExists(filepath.Join(barrelPath, "yarn.lock")) {
			testCmd = "yarn test"
			lintCmd = "yarn lint"
			buildCmd = "yarn build"
		}

		return InferredTechStack{
			Language:          "TypeScript / JavaScript",
			TestRunner:        testCmd,
			Linter:            lintCmd,
			CoverageThreshold: ">80%",
			BuildTool:         buildCmd,
		}
	}

	// 4. Rust
	if fileExists(filepath.Join(barrelPath, "Cargo.toml")) {
		return InferredTechStack{
			Language:          "Rust",
			TestRunner:        "cargo test",
			Linter:            "cargo clippy",
			CoverageThreshold: ">80%",
			BuildTool:         "cargo build",
		}
	}

	// 5. Python
	if fileExists(filepath.Join(barrelPath, "pyproject.toml")) ||
		fileExists(filepath.Join(barrelPath, "requirements.txt")) ||
		fileExists(filepath.Join(barrelPath, "Pipfile")) {
		return InferredTechStack{
			Language:          "Python 3.11+",
			TestRunner:        "pytest --cov",
			Linter:            "ruff check",
			CoverageThreshold: ">80%",
			BuildTool:         "pip install",
		}
	}

	// 6. Java / Kotlin
	if fileExists(filepath.Join(barrelPath, "pom.xml")) || fileExists(filepath.Join(barrelPath, "build.gradle")) {
		return InferredTechStack{
			Language:          "Java / Kotlin",
			TestRunner:        "mvn test",
			Linter:            "mvn checkstyle:check",
			CoverageThreshold: ">80%",
			BuildTool:         "mvn package",
		}
	}

	// Fallback generic
	return InferredTechStack{
		Language:          "Generic Polyglot",
		TestRunner:        "make test",
		Linter:            "make lint",
		CoverageThreshold: ">80%",
		BuildTool:         "make build",
	}
}

// ScaffoldBarrelTechStack generates .cooper/definition/tech-stack.md in the barrel directory.
func ScaffoldBarrelTechStack(barrelPath string, opts ScaffoldOptions) (*ScaffoldResult, error) {
	if _, err := os.Stat(barrelPath); err != nil {
		return nil, fmt.Errorf("barrel directory %q does not exist: %w", barrelPath, err)
	}

	inferred := InferTechStack(barrelPath)

	lang := inferred.Language
	if strings.TrimSpace(opts.Language) != "" {
		lang = opts.Language
	}

	framework := inferred.Framework
	if strings.TrimSpace(opts.Framework) != "" {
		framework = opts.Framework
	}

	testRunner := inferred.TestRunner
	if strings.TrimSpace(opts.TestRunner) != "" {
		testRunner = opts.TestRunner
	}

	linter := inferred.Linter
	if strings.TrimSpace(opts.Linter) != "" {
		linter = opts.Linter
	}

	coverage := inferred.CoverageThreshold
	if strings.TrimSpace(opts.CoverageThreshold) != "" {
		coverage = opts.CoverageThreshold
	}

	cooperDefDir := filepath.Join(barrelPath, ".cooper", "definition")
	if err := os.MkdirAll(cooperDefDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory %s: %w", cooperDefDir, err)
	}

	techStackPath := filepath.Join(cooperDefDir, "tech-stack.md")
	exists := fileExists(techStackPath)

	if exists && !opts.Force {
		return nil, fmt.Errorf("tech-stack.md already exists at %s (use force to overwrite)", techStackPath)
	}

	var sb strings.Builder
	sb.WriteString("# Technology Stack & Platform Contracts\n\n")
	sb.WriteString("## Pattern Components\n\n")
	sb.WriteString("| Component | Role | Metaphor | Tech Stack / Tooling |\n")
	sb.WriteString("| :--- | :--- | :--- | :--- |\n")
	sb.WriteString(fmt.Sprintf("| **Core Runtime** | Application Logic | Barrel Core | %s |\n", lang))
	if framework != "" {
		sb.WriteString(fmt.Sprintf("| **Framework** | Architecture & Routing | Frame | %s |\n", framework))
	}
	sb.WriteString(fmt.Sprintf("| **Testing Engine** | Unit & Integration TDD | Test Rig | `%s` (Threshold: %s) |\n", testRunner, coverage))
	if linter != "" {
		sb.WriteString(fmt.Sprintf("| **Code Quality** | Linter & Style Enforcement | Quality Gate | `%s` |\n", linter))
	}
	if inferred.BuildTool != "" {
		sb.WriteString(fmt.Sprintf("| **Build & Distribution** | Artifact Compiler | Build Engine | `%s` |\n", inferred.BuildTool))
	}
	sb.WriteString("\n## Quality Gates & Execution\n")
	sb.WriteString(fmt.Sprintf("- **Test Command**: `%s`\n", testRunner))
	sb.WriteString(fmt.Sprintf("- **Coverage Gate**: `%s` statement/branch coverage\n", coverage))
	if linter != "" {
		sb.WriteString(fmt.Sprintf("- **Lint Command**: `%s`\n", linter))
	}
	sb.WriteString("\n---\n*Generated by Battery Multi-Barrel Orchestration Protocol*\n")

	if err := os.WriteFile(techStackPath, []byte(sb.String()), 0644); err != nil {
		return nil, fmt.Errorf("failed to write tech-stack.md: %w", err)
	}

	// Create default styleguide if not present
	styleguideDir := filepath.Join(barrelPath, ".cooper", "code_styleguides")
	var styleguidePath string
	if err := os.MkdirAll(styleguideDir, 0755); err == nil {
		cleanLang := strings.ToLower(strings.Fields(lang)[0])
		styleguidePath = filepath.Join(styleguideDir, cleanLang+".md")
		if !fileExists(styleguidePath) {
			guideContent := fmt.Sprintf("# Code Styleguide: %s\n\n- Adhere to idiomatic %s conventions.\n- Enforce strict unit test coverage (%s).\n- Run `%s` before committing.\n", lang, lang, coverage, linter)
			_ = os.WriteFile(styleguidePath, []byte(guideContent), 0644)
		}
	}

	return &ScaffoldResult{
		BarrelPath:     barrelPath,
		TechStackPath:  techStackPath,
		StyleguidePath: styleguidePath,
		Inferred:       &inferred,
		Created:        !exists,
		Overwritten:    exists,
	}, nil
}

func fileExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && !fi.IsDir()
}
