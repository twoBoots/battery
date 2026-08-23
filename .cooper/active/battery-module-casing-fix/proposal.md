# Proposal: Battery Go Module Casing Standardization

## Intent
Standardize the Go module path and package import paths across the `battery` repository to `github.com/twoBoots/battery`, matching the canonical casing of the GitHub organization `twoBoots` and aligning with sibling repositories (`twoBoots/bender`, `twoBoots/cooper`, `twoBoots/troop`).

## Scope Boundaries
- **Module Declaration:** Update line 1 of `go.mod` from `module github.com/twoboots/battery` to `module github.com/twoBoots/battery`.
- **Internal Imports:** Rewrite all internal import statements in `main.go`, `cmd/**/*.go`, and `internal/**/*.go` to use `github.com/twoBoots/battery/...`.
- **CI / Build Configuration:** Update ldflags in `.github/workflows/release.yml` so that binary version injection targets `-X github.com/twoBoots/battery/cmd.Version=${VERSION}`.
- **Verification:** Ensure all tests pass with `go test ./...`, linting passes with `go vet ./...`, and the binary builds cleanly.
