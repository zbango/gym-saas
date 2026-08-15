import { spawnSync } from "node:child_process";
import { resolve } from "node:path";

const repoRoot = process.cwd();
const desktopDir = resolve(repoRoot, "apps/desktop");
const args = ["run", "github.com/wailsapp/wails/v2/cmd/wails@v2.10.1", "build"];
const target = process.argv[2];

if (target) {
  args.push("-platform", target);
}

if ((target && target.startsWith("windows")) || (!target && process.platform === "win32")) {
  args.push("-nsis");
}

const result = spawnSync("go", args, {
  cwd: desktopDir,
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
