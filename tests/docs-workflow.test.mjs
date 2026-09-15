import { test } from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const repoRoot = path.resolve(__dirname, '..');

test('GitHub Actions workflow .github/workflows/pages.yml exists and has proper configuration', () => {
  const workflowPath = path.join(repoRoot, '.github', 'workflows', 'pages.yml');
  assert.ok(fs.existsSync(workflowPath), '.github/workflows/pages.yml must exist');

  const content = fs.readFileSync(workflowPath, 'utf8');

  // Triggers
  assert.ok(content.includes('branches: [main]') || content.includes('- main'), 'workflow must trigger on push to main');
  assert.ok(content.includes('workflow_dispatch'), 'workflow must allow manual dispatch');

  // Permissions
  assert.ok(content.includes('pages: write'), 'must have pages: write permission');
  assert.ok(content.includes('id-token: write'), 'must have id-token: write permission');

  // Build & Deploy steps
  assert.ok(content.includes('actions/configure-pages'), 'must use actions/configure-pages');
  assert.ok(content.includes('actions/upload-pages-artifact'), 'must use actions/upload-pages-artifact');
  assert.ok(content.includes('actions/deploy-pages'), 'must use actions/deploy-pages');
  assert.ok(content.includes('docs/.vitepress/dist'), 'artifact path must be docs/.vitepress/dist');
  assert.ok(content.includes('npm ci'), 'must use npm ci for reproducible build');
  assert.ok(content.includes('npm run docs:build'), 'must execute npm run docs:build');
});
