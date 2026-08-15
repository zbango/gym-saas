import { mkdirSync } from "node:fs";
import { spawnSync } from "node:child_process";
import { resolve } from "node:path";

const repoRoot = process.cwd();
const outDir = resolve(repoRoot, "apps/cloud-api/build");

mkdirSync(outDir, { recursive: true });

function run(args) {
  const result = spawnSync("go", args, {
    cwd: repoRoot,
    env: {
      ...process.env,
      GOWORK: resolve(repoRoot, "go.work"),
      GOCACHE: resolve(repoRoot, ".gocache")
    },
    stdio: "inherit"
  });

  if (result.status !== 0) {
    process.exit(result.status ?? 1);
  }
}

run(["build", "-o", resolve(outDir, "cloud-api-lambda"), "./apps/cloud-api/lambda"]);
