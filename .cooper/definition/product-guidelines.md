# Multi-Barrel Product Guidelines & Standards

## Documentation & Specification Standards
- **Agnostic & Modular**: Specifications must describe interfaces, data contracts, and barrel responsibilities clearly without lock-in to specific repositories.
- **Contract-First**: Every multi-barrel track must define data schemas, endpoints, or messaging interfaces in `spec.md` before implementation begins in target barrels.
- **Explicit Dependencies**: Inter-barrel dependencies must be explicitly declared in `metadata.json` and phased in `plan.md`.

## Architectural Principles
1. **Decoupled Barrels**: Barrels interact strictly through documented contracts (OpenAPI, gRPC, Protobuf, JSON-RPC/MCP).
2. **Backward Compatibility**: Interface additions must preserve backward compatibility across independent barrels.
3. **Cooper Worktree Alignment**: Branch naming across target barrels must match the `battery` track ID for consistent auditing.
4. **Verifiable Checkpoints**: Each downstream track must provide automated unit/integration tests for its contract implementation.
