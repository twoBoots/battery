# Execution Plan: Multi-Barrel Track Dispatch & Decoupled Planning

Follow strict TDD (Red -> Green -> Refactor) and attach Git Notes to task commits. Maintain >80% test coverage.

---

## Phase 1: Core Track Models & Barrel Location Scanner

- [x] Task 1.1: Define Track and Barrel Track models (`TrackMetadata`, `BarrelTrackSummary`, location enums) in `internal/track/models.go` with unit tests. (d5d8c67)
- [x] Task 1.2: Implement `LocateBarrelTrack(barrelPath, trackID)` to scan `.worktrees/<track_id>/.cooper/active/`, `.cooper/active/`, and `.cooper/archive/` with unit tests. (0d724a3)
- [x] Task 1.3: Add tests verifying correct detection of completed/archived barrel tracks. (0d724a3)

---

## Phase 2: Track Initialization & Spec Dispatch Engine

- [x] Task 2.1: Implement `InitTrack(cwd, trackID, barrels, opts)` to scaffold `.cooper/active/<track_id>/` in Battery with unit tests. (7d9ccbf)
- [x] Task 2.2: Implement `DispatchTrack(cwd, trackID, opts)` to seed target barrels with `metadata.json`, `proposal.md`, and `spec-deltas/` while omitting `plan.md`. (7d9ccbf)
- [x] Task 2.3: Add tests verifying that `plan.md` is omitted in dispatched barrels to preserve local planning autonomy. (7d9ccbf)

---

## Phase 3: Status Aggregation & Plan Task Parser

- [ ] Task 3.1: Implement `ParsePlanTasks(planContent)` to extract total, completed `[x]`, in-progress `[~]`, and pending `[ ]` tasks with unit tests.
- [ ] Task 3.2: Implement `GetMultiBarrelTrackStatus(cwd, trackID)` aggregating status across active worktrees, main trunks, and archived tracks.
- [ ] Task 3.3: Implement `ListTracks(cwd)` to scan and list all active and archived tracks in Battery.

---

## Phase 4: CLI Command Suite (`cmd/track.go`)

- [ ] Task 4.1: Implement `battery track init <track_id> --barrels <names>` Cobra command with tests.
- [ ] Task 4.2: Implement `battery track dispatch <track_id>` Cobra command with tests.
- [ ] Task 4.3: Implement `battery track status [<track_id>]` with formatted table output supporting active and archived barrels.
- [ ] Task 4.4: Implement `battery track list` command with tests.

---

## Phase 5: End-to-End Verification & Documentation

- [ ] Task 5.1: Run full test suite with coverage validation (`go test -v -cover ./...`) ensuring coverage >80%.
- [ ] Task 5.2: Execute manual verification workflow using sibling mock barrel directories.
