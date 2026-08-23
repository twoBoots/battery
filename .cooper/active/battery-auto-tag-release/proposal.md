# Proposal: Battery Auto-Tag Release Pipeline (battery)

## Intent
Port `bender`'s `auto-tag` GitHub Actions job into `battery/.github/workflows/release.yml` so that whenever `cmd/root.go` has its version bumped in a PR, merging to `main` automatically inspects `origin`, tags `v<Version>` if missing, and publishes release assets under both the semantic version tag and `latest`.

## Next Steps
Run local planning (`cooper-implement` or local track planning) in `battery` to generate `design.md` and `plan.md` adhering to this repository's tech stack, then execute.
