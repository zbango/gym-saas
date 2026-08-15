import { readFileSync, writeFileSync } from "node:fs";
import { resolve } from "node:path";

const repoRoot = process.cwd();
const versionFilePath = resolve(repoRoot, "VERSION");

export function readCurrentVersion() {
  return readFileSync(versionFilePath, "utf8").trim();
}

export function bumpPatch(version) {
  const parts = version.split(".").map((part) => Number.parseInt(part, 10));
  if (parts.length !== 3 || parts.some((part) => Number.isNaN(part) || part < 0)) {
    throw new Error(`Invalid semantic version: ${version}`);
  }

  return `${parts[0]}.${parts[1]}.${parts[2] + 1}`;
}

export function syncVersionFiles(nextVersion) {
  writeFile(versionFilePath, nextVersion);
  updateByPattern("packages/shared/src/index.ts", /desktopVersion = "[^"]+"/, `desktopVersion = "${nextVersion}"`);
  updateByPattern("go/core/platform/version.go", /const Version = "[^"]+"/, `const Version = "${nextVersion}"`);
  updateByPattern("apps/desktop/wails.json", /"productVersion": "[^"]+"/, `"productVersion": "${nextVersion}"`);
}

function writeFile(path, value) {
  writeFileSync(resolve(repoRoot, path), `${value}\n`, "utf8");
}

function updateByPattern(path, pattern, replacement) {
  const absolute = resolve(repoRoot, path);
  const current = readFileSync(absolute, "utf8");
  const next = current.replace(pattern, replacement);
  if (next === current) {
    throw new Error(`Pattern not found when updating ${path}`);
  }
  writeFileSync(absolute, next, "utf8");
}
