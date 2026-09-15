# Technical Design: VitePress GitHub Pages Documentation Site

## Architecture & System Structure

### 1. VitePress Framework & Node.js Tooling
- **Engine**: VitePress `^1.6.3`
- **Configuration**: `docs/.vitepress/config.mts`
  - Site Title: `Battery 🔋`
  - Base Path: `process.env.VITEPRESS_BASE || '/battery/'`
  - Search: Local search provider
  - Nav & Sidebar: Structured into "Introduction" and "Deep Dive" sections.
  - Social Links: GitHub link to `https://github.com/twoBoots/battery`.
- **Scripts**:
  - `docs:dev`: Local live-reload VitePress server.
  - `docs:build`: Production static site compilation.
  - `docs:preview`: Local preview of compiled build.
  - `test`: Node.js native test runner `node --test tests/*.test.mjs`.

### 2. Information Architecture (`docs/`)
```
docs/
├── .vitepress/
│   └── config.mts                 # VitePress configuration & theme
├── index.md                       # Landing page with hero & feature pillars
├── guide/
│   ├── getting-started.md         # Fast onboarding, CLI basics & status check
│   └── workflow.md                # Multi-barrel track lifecycle (init -> dispatch -> worktree)
├── architecture.md                # Batteries & Barrels model, topology & boundaries
├── mcp.md                         # Model Context Protocol server tools & prompts
└── installation.md                # CLI curl installer, compile & troubleshooting
```

### 3. CI/CD Deployment (`.github/workflows/pages.yml`)
- Triggered on push to `main` for paths:
  - `docs/**`
  - `package*.json`
  - `.github/workflows/pages.yml`
- Environment: `github-pages`
- Uses official GitHub Actions:
  - `actions/checkout@v4`
  - `actions/setup-node@v4` (Node 20, cache: npm)
  - `actions/configure-pages@v5`
  - `actions/upload-pages-artifact@v3` (path: `docs/.vitepress/dist`)
  - `actions/deploy-pages@v4`

### 4. Quality Verification & Testing (`tests/`)
- `tests/docs-package.test.mjs`: Tests package configuration, dependencies, and scripts.
- `tests/docs-config.test.mjs`: Tests VitePress configuration structure, base path, title, and links.
- `tests/docs-landing.test.mjs`: Tests landing page layout, hero, quickstart command, and feature pillars.
- `tests/docs-guides.test.mjs`: Tests presence of all guide markdown files and markdown link integrity.
- `tests/docs-workflow.test.mjs`: Tests GitHub Pages workflow configuration, triggers, and permissions.
