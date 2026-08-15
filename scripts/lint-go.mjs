import { spawnSync } from "node:child_process";
import { resolve } from "node:path";

const repoRoot = process.cwd();
const goTargets = ["./apps/cloud-api/...", "./apps/desktop/...", "./go/core/..."];

function run(command, args, options = {}) {
  const result = spawnSync(command, args, {
    cwd: repoRoot,
    env: {
      ...process.env,
      GOWORK: resolve(repoRoot, "go.work"),
      GOCACHE: resolve(repoRoot, ".gocache")
    },
    encoding: "utf8",
    ...options
  });

  if (result.error) {
    throw result.error;
  }

  return result;
}

const gofmt = run("gofmt", ["-l", "apps/cloud-api", "apps/desktop", "go/core"]);
const dirty = gofmt.stdout.trim();

if (dirty) {
  console.error("gofmt found unformatted files:");
  console.error(dirty);
  process.exit(1);
}

const vet = run("go", ["vet", ...goTargets], { stdio: "inherit" });
if (vet.status !== 0) {
  process.exit(vet.status ?? 1);
}
