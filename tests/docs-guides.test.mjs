import { test } from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const repoRoot = path.resolve(__dirname, '..');
const docsDir = path.join(repoRoot, 'docs');

function getMarkdownFiles(dir) {
  let results = [];
  const list = fs.readdirSync(dir);
  for (const file of list) {
    const filePath = path.join(dir, file);
    const stat = fs.statSync(filePath);
    if (stat.isDirectory()) {
      if (file !== '.vitepress') {
        results = results.concat(getMarkdownFiles(filePath));
      }
    } else if (file.endsWith('.md')) {
      results.push(filePath);
    }
  }
  return results;
}

test('Guide markdown files exist', () => {
  const gettingStarted = path.join(docsDir, 'guide', 'getting-started.md');
  const workflow = path.join(docsDir, 'guide', 'workflow.md');
  const architecture = path.join(docsDir, 'architecture.md');
  const mcp = path.join(docsDir, 'mcp.md');
  const installation = path.join(docsDir, 'installation.md');

  assert.ok(fs.existsSync(gettingStarted), 'docs/guide/getting-started.md must exist');
  assert.ok(fs.existsSync(workflow), 'docs/guide/workflow.md must exist');
  assert.ok(fs.existsSync(architecture), 'docs/architecture.md must exist');
  assert.ok(fs.existsSync(mcp), 'docs/mcp.md must exist');
  assert.ok(fs.existsSync(installation), 'docs/installation.md must exist');
});

test('All Cooper and Troop references in docs/ link to upstream repositories', () => {
  const mdFiles = getMarkdownFiles(docsDir);
  assert.ok(mdFiles.length >= 5, 'Must find documentation markdown files');

  for (const file of mdFiles) {
    const content = fs.readFileSync(file, 'utf8');
    const relativePath = path.relative(repoRoot, file);

    if (content.includes('Cooper')) {
      const hasCooperLink = content.includes('https://twoboots.github.io/cooper') || content.includes('https://github.com/twoBoots/cooper');
      assert.ok(
        hasCooperLink,
        `File ${relativePath} mentions Cooper but is missing link to https://twoboots.github.io/cooper or https://github.com/twoBoots/cooper`
      );
    }

    if (content.includes('Troop')) {
      const hasTroopLink = content.includes('https://twoboots.github.io/troop') || content.includes('https://github.com/twoBoots/troop');
      assert.ok(
        hasTroopLink,
        `File ${relativePath} mentions Troop but is missing link to https://twoboots.github.io/troop or https://github.com/twoBoots/troop`
      );
    }
  }
});

test('All local markdown links in docs/ resolve to existing files', () => {
  const mdFiles = getMarkdownFiles(docsDir);
  const linkRegex = /\[([^\]]+)\]\(([^)]+)\)/g;

  for (const file of mdFiles) {
    const content = fs.readFileSync(file, 'utf8');
    const fileDir = path.dirname(file);
    let match;
    while ((match = linkRegex.exec(content)) !== null) {
      const linkTarget = match[2].trim();
      if (
        linkTarget.startsWith('http://') ||
        linkTarget.startsWith('https://') ||
        linkTarget.startsWith('mailto:') ||
        linkTarget.startsWith('#')
      ) {
        continue;
      }
      const cleanTarget = linkTarget.split('#')[0].split('?')[0];
      if (!cleanTarget) continue;

      const resolved = path.resolve(fileDir, cleanTarget);
      assert.ok(
        fs.existsSync(resolved),
        `Link [${match[1]}](${linkTarget}) in ${path.relative(repoRoot, file)} points to non-existent file: ${resolved}`
      );
    }
  }
});
