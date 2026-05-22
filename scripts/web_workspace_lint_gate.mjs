import { promises as fs } from 'node:fs';
import path from 'node:path';

const ROOT = process.cwd();
const SOURCE_ROOTS = ['apps', 'packages'];
const IGNORED_SEGMENTS = new Set([
  '.next',
  '.turbo',
  'build',
  'dist',
  'node_modules',
  'coverage',
  'generated',
]);
const SOURCE_EXTENSIONS = new Set(['.ts', '.tsx']);
const CONFLICT_MARKERS = ['<<<<<<< ', '=======', '>>>>>>> '];

async function pathExists(filePath) {
  try {
    await fs.access(filePath);
    return true;
  } catch {
    return false;
  }
}

function isIgnored(entryPath) {
  return entryPath
    .split(path.sep)
    .some((segment) => IGNORED_SEGMENTS.has(segment));
}

async function walk(dirPath, files = []) {
  if (!(await pathExists(dirPath))) {
    return files;
  }

  const entries = await fs.readdir(dirPath, { withFileTypes: true });
  for (const entry of entries) {
    const fullPath = path.join(dirPath, entry.name);
    if (isIgnored(path.relative(ROOT, fullPath))) {
      continue;
    }
    if (entry.isDirectory()) {
      await walk(fullPath, files);
      continue;
    }
    if (SOURCE_EXTENSIONS.has(path.extname(entry.name))) {
      files.push(fullPath);
    }
  }
  return files;
}

async function assertPackageTypeConfigs() {
  const failures = [];
  for (const packageRoot of ['packages/store', 'packages/api-client']) {
    const packageJsonPath = path.join(ROOT, packageRoot, 'package.json');
    const tsconfigPath = path.join(ROOT, packageRoot, 'tsconfig.json');
    if ((await pathExists(packageJsonPath)) && !(await pathExists(tsconfigPath))) {
      failures.push(`${packageRoot} has a typecheck script but no tsconfig.json`);
    }
  }
  return failures;
}

async function assertNoConflictMarkers(files) {
  const failures = [];
  for (const file of files) {
    const contents = await fs.readFile(file, 'utf8');
    const lines = contents.split('\n');
    lines.forEach((line, index) => {
      if (CONFLICT_MARKERS.some((marker) => line.startsWith(marker))) {
        failures.push(`${path.relative(ROOT, file)}:${index + 1} contains a merge conflict marker`);
      }
    });
  }
  return failures;
}

const sourceFiles = (
  await Promise.all(SOURCE_ROOTS.map((sourceRoot) => walk(path.join(ROOT, sourceRoot))))
).flat();

const failures = [
  ...(await assertPackageTypeConfigs()),
  ...(await assertNoConflictMarkers(sourceFiles)),
];

if (failures.length > 0) {
  console.error('WEB_WORKSPACE_LINT_GATE_FAIL');
  failures.forEach((failure) => console.error(`- ${failure}`));
  process.exit(1);
}

console.log('WEB_WORKSPACE_LINT_GATE_PASS');
console.log(`WEB_WORKSPACE_LINT_GATE_SCANNED_FILES=${sourceFiles.length}`);
