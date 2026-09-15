import { test } from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const repoRoot = path.resolve(__dirname, '..');

test('package.json exists and defines required VitePress documentation scripts and dependencies', () => {
  const pkgPath = path.join(repoRoot, 'package.json');
  assert.ok(fs.existsSync(pkgPath), 'package.json must exist at repository root');

  const pkg = JSON.parse(fs.readFileSync(pkgPath, 'utf8'));

  assert.equal(pkg.name, 'battery-docs', 'package name must be battery-docs');
  assert.equal(pkg.type, 'module', 'package type must be module');
  assert.equal(pkg.private, true, 'package must be marked private');

  // Scripts
  assert.ok(pkg.scripts, 'package.json must define scripts');
  assert.ok(pkg.scripts['docs:dev'], 'scripts must include docs:dev');
  assert.ok(pkg.scripts['docs:build'], 'scripts must include docs:build');
  assert.ok(pkg.scripts['docs:preview'], 'scripts must include docs:preview');
  assert.ok(pkg.scripts['test'], 'scripts must include test');

  // Dependencies
  assert.ok(pkg.devDependencies, 'package.json must define devDependencies');
  assert.ok(pkg.devDependencies['vitepress'], 'devDependencies must include vitepress');
});
