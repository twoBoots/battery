import { test } from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const repoRoot = path.resolve(__dirname, '..');

test('VitePress config file exists and contains required site structure', async () => {
  const configPath = path.join(repoRoot, 'docs', '.vitepress', 'config.mts');
  assert.ok(fs.existsSync(configPath), 'docs/.vitepress/config.mts must exist');

  const content = fs.readFileSync(configPath, 'utf8');
  assert.ok(content.includes(`title: 'Battery'`), 'config must set title to Battery');
  assert.ok(content.includes('base:'), 'config must specify base URL');
  assert.ok(content.includes('/battery/'), 'base must include /battery/');
  assert.ok(content.includes('https://github.com/twoBoots/battery'), 'socialLinks must link to twoBoots/battery repository');
  assert.ok(content.includes('nav:'), 'themeConfig must define navigation bar');
  assert.ok(content.includes('sidebar:'), 'themeConfig must define sidebar');
  assert.ok(content.includes('provider: \'local\''), 'themeConfig must define local search provider');
});
