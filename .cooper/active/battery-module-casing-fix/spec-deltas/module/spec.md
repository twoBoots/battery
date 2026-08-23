# Capability Spec Delta: Go Module Path Casing Standardization

## Capability: `module`

### Description
Standardizes the Go module path naming, internal package import paths, and build linker flags for the `battery` repository to strictly use the canonical `github.com/twoBoots/battery` casing, aligning with the `twoBoots` organization standard and sibling barrels (`twoBoots/bender`, `twoBoots/cooper`, `twoBoots/troop`).

---

### Added Requirements

#### Requirement: Canonical Module Path Casing
`+` The `battery` repository Go module declaration in `go.mod` SHALL be `module github.com/twoBoots/battery`.

##### Scenario: Go Module Declaration
- `+` GIVEN the `battery` Go project root
- `+` WHEN `go.mod` is inspected
- `+` THEN line 1 MUST declare `module github.com/twoBoots/battery`.

##### Scenario: Internal Package Imports
- `+` GIVEN any Go source file in `battery` (`main.go`, `cmd/*.go`, `internal/**/*.go`)
- `+` WHEN importing packages from within the `battery` module
- `+` THEN the import path MUST begin with `github.com/twoBoots/battery/`.

##### Scenario: CI Release Version Injection
- `+` GIVEN `.github/workflows/release.yml`
- `+` WHEN the release binary build step executes
- `+` THEN the linker ldflags MUST target `-X github.com/twoBoots/battery/cmd.Version=${VERSION}`.

---

### Removed / Deprecated Requirements

#### Requirement: Lowercase Organization Module Declaration
`-` The `battery` repository SHALL NOT declare `module github.com/twoboots/battery` in `go.mod` or import packages using lowercase `github.com/twoboots/battery/...`.
