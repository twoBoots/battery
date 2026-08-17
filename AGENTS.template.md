# Battery Agent Rules (Multi-Barrel Orchestrator)

1. **Isolation Protocol**: Multi-barrel tracks defined in `battery` dispatch tasks into target barrel worktrees under `.worktrees/<track_id>`.
2. **Branching Strategy**: Track IDs must be consistent across `battery` and target barrels.
3. **Execution & Cleanup**: Use `git agent-start <track_id>` in target barrels for isolated work, and `git agent-stop <track_id>` upon PR merge.
4. **Configuration & Topology**: Respect `.batteryrc` (canonical) and `.batteryrc.local` (local overrides). Do not hardcode tech stacks; resolve each barrel's `.cooper/definition/tech-stack.md`.

See [.cooper/BATTERY.md](.cooper/BATTERY.md) and [.cooper/COOPER.md](.cooper/COOPER.md) for full context.
