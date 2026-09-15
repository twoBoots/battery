import { test } from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const repoRoot = path.resolve(__dirname, '..');

test('Landing page docs/index.md exists with required layout and sections', () => {
  const indexPath = path.join(repoRoot, 'docs', 'index.md');
  assert.ok(fs.existsSync(indexPath), 'docs/index.md must exist');

  const content = fs.readFileSync(indexPath, 'utf8');

  // Home layout & hero
  assert.ok(content.includes('layout: home'), 'docs/index.md frontmatter must declare layout: home');
  assert.ok(content.includes('name: Battery'), 'hero must specify name: Battery');
  assert.ok(content.includes('tagline:'), 'hero must specify tagline');

  // Quickstart command
  assert.ok(
    content.includes('curl -fsSL') && content.includes('install.sh'),
    'landing page must contain official curl quickstart install command'
  );

  // Core pillars
  assert.ok(content.includes('Multi-Barrel') || content.includes('Multi-Repository'), 'must feature multi-barrel/repo coordination');
  assert.ok(content.includes('https://github.com/twoBoots/cooper'), 'must link Cooper to https://github.com/twoBoots/cooper');
  assert.ok(content.includes('https://github.com/twoBoots/troop'), 'must link Troop to https://github.com/twoBoots/troop');
  assert.ok(content.includes('Model Context Protocol') || content.includes('MCP'), 'must feature MCP server');
  assert.ok(content.includes('.batteryrc'), 'must mention .batteryrc topology');

  // Workflow section
  assert.ok(content.includes('Workflow') || content.includes('Lifecycle'), 'must include workflow or lifecycle section');
});
