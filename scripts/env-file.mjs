import { existsSync, readFileSync, writeFileSync } from "node:fs";

export function readEnvFile(path) {
  if (!existsSync(path)) {
    return {};
  }

  const content = readFileSync(path, "utf8");
  const entries = {};

  for (const line of content.split(/\r?\n/)) {
    const trimmed = line.trim();
    if (!trimmed || trimmed.startsWith("#")) {
      continue;
    }

    const index = trimmed.indexOf("=");
    if (index === -1) {
      continue;
    }

    const key = trimmed.slice(0, index).trim();
    let value = trimmed.slice(index + 1).trim();

    if (
      (value.startsWith('"') && value.endsWith('"')) ||
      (value.startsWith("'") && value.endsWith("'"))
    ) {
      value = value.slice(1, -1);
    }

    entries[key] = value;
  }

  return entries;
}

export function writeEnvFile(path, values) {
  const keys = Object.keys(values).sort();
  const lines = keys.map((key) => `${key}=${escapeValue(values[key] ?? "")}`);
  writeFileSync(path, `${lines.join("\n")}\n`, "utf8");
}

function escapeValue(value) {
  if (/[\s#"'`]/.test(value)) {
    return JSON.stringify(value);
  }

  return value;
}
