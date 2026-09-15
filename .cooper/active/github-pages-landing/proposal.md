# Proposal: VitePress GitHub Pages Documentation Site

## Rationale
Sibling repositories [`twoBoots/cooper`](https://github.com/twoBoots/cooper) and [`twoBoots/troop`](https://github.com/twoBoots/troop) recently launched public-facing, responsive VitePress documentation websites hosted on GitHub Pages with automated CI/CD deployment pipelines.

Battery is the multi-repository Spec-Driven Development (SDD) orchestrator that binds Cooper and Troop across a collection of barrels. To complete the framework's documentation ecosystem, Battery requires the same clean, responsive documentation site:
1. **Clear Onboarding**: Immediate value proposition, quickstart one-liner, and visual workflow diagrams.
2. **Comprehensive Guides**: Getting started, multi-barrel workflow lifecycle, layered topology (`.batteryrc` & `.batteryrc.local`), and spec-only boundaries.
3. **Deep Dives**: Full architectural breakdown of the Batteries & Barrels model and Model Context Protocol (MCP) server capabilities for AI coding assistants.
4. **Automated CI/CD**: Seamless GitHub Actions deployment to GitHub Pages on push to `main`.
5. **Ecosystem Parity**: Automated tests and consistent styling matching Cooper and Troop documentation sites.

## Scope Boundaries
- Setup `package.json` with `vitepress` and Node native test scripts.
- Configure VitePress in `docs/.vitepress/config.mts` with base path `/battery/` and local search.
- Create landing page `docs/index.md` with Hero banner, quickstart command, feature pillars, and Mermaid workflow diagram.
- Author documentation guides:
  - `docs/guide/getting-started.md`
  - `docs/guide/workflow.md`
  - `docs/architecture.md`
  - `docs/mcp.md`
  - `docs/installation.md`
- Configure GitHub Pages CI/CD workflow in `.github/workflows/pages.yml`.
- Author integration tests in `tests/*.test.mjs` verifying configuration, link integrity, and layout.
- Update `README.md` to link to the documentation site.
