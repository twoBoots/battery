import { defineConfig } from 'vitepress';

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: 'Battery',
  description: 'Multi-Repository Specification-Driven Development (SDD) Orchestrator for Humans and AI Agents',
  base: process.env.VITEPRESS_BASE || '/battery/',
  themeConfig: {
    siteTitle: 'Battery 🔋',
    nav: [
      { text: 'Guide', link: '/guide/getting-started' },
      { text: 'Workflow', link: '/guide/workflow' },
      { text: 'Architecture', link: '/architecture' },
      { text: 'MCP Server', link: '/mcp' },
      { text: 'Installation', link: '/installation' }
    ],
    sidebar: [
      {
        text: 'Introduction',
        items: [
          { text: 'Getting Started', link: '/guide/getting-started' },
          { text: 'Multi-Barrel Workflow', link: '/guide/workflow' },
          { text: 'Installation Guide', link: '/installation' }
        ]
      },
      {
        text: 'Deep Dive',
        items: [
          { text: 'Architecture (Barrels & Batteries)', link: '/architecture' },
          { text: 'Model Context Protocol (MCP)', link: '/mcp' }
        ]
      }
    ],
    socialLinks: [
      { icon: 'github', link: 'https://github.com/twoBoots/battery' }
    ],
    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2026 twoBoots'
    },
    search: {
      provider: 'local'
    }
  }
});

