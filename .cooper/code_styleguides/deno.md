# Deno & TypeScript Code Styleguide

## 1. Project & Task Configuration (`deno.json`)
- **No Ad-Hoc Build Scripts**: Always define tasks directly within `deno.json` under `"tasks"`.
- **Standard Tasks**:
  - `build`: `deno compile --allow-read --allow-write --allow-env -o bin/<name> src/cli/mod.ts`
  - `install`: `deno install -g --allow-read --allow-write --allow-env --force -n <name> src/cli/mod.ts`
  - `test`: `deno test --allow-read --allow-write --allow-env src/`
  - `test:coverage`: `deno test --allow-read --allow-write --allow-env --coverage=coverage src/ && deno coverage coverage`
  - `lint`: `deno lint src/`
  - `fmt`: `deno fmt src/`

## 2. Directory & Module Structure
- `src/mod.ts`: Canonical library entry point exporting public APIs and types.
- `src/types.ts`: TypeScript data structures, interfaces, and types.
- `src/cli/mod.ts`: CLI entry point and argument dispatcher with `import.meta.main` guard.
- **Colocated Tests**: Test files are colocated directly beside the source files they test using the `*.test.ts` naming convention (e.g. `src/config.ts` and `src/config.test.ts`).

## 3. Standard Library & Imports
- Use import map in `deno.json` for standard packages (e.g. `"@std/assert": "jsr:@std/assert@^1.0.0"`).
- Use `node:path` for standard path resolution.
- Keep dependencies minimal; prefer Deno built-in APIs (`Deno.readTextFile`, `Deno.writeTextFile`, `Deno.stat`, `Deno.readDir`).
- Avoid explicit `any`; use `unknown` and proper type guards.

## 4. Testing & Quality Gates
- Maintain >80% test coverage for all new modules.
- Ensure all tests pass with `deno test`.
- Code must be formatted with `deno fmt` and pass `deno lint` with 0 warnings.
