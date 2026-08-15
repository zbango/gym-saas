import { spawn } from "node:child_process";
import { resolve } from "node:path";
import { readEnvFile } from "./env-file.mjs";

const repoRoot = process.cwd();
const dryRun = process.argv.includes("--dry-run");
const modes = process.argv.slice(2).filter((arg) => !arg.startsWith("--"));
const envFromFiles = {
  ...readEnvFile(resolve(repoRoot, ".env.example")),
  ...readEnvFile(resolve(repoRoot, ".env"))
};

if (modes.length === 0) {
  console.error("Usage: node ./scripts/dev-stack.mjs <desktop|web> [desktop|web] [--dry-run]");
  process.exit(1);
}

const processes = [];

if (modes.includes("desktop")) {
  processes.push({
    name: "desktop-frontend",
    command: "npm",
    args: ["run", "dev", "-w", "@gym-saas/desktop-frontend"],
    cwd: repoRoot,
    env: {
      ...process.env,
      ...envFromFiles
    }
  });

  processes.push({
    name: "desktop",
    command: "go",
    args: ["run", "github.com/wailsapp/wails/v2/cmd/wails@v2.10.1", "dev"],
    cwd: resolve(repoRoot, "apps/desktop"),
    env: {
      ...process.env,
      ...envFromFiles,
      GOWORK: resolve(repoRoot, "go.work"),
      GOCACHE: resolve(repoRoot, ".gocache")
    }
  });
}

if (modes.includes("web")) {
  processes.push({
    name: "web",
    command: "npm",
    args: ["run", "dev", "-w", "@gym-saas/web"],
    cwd: repoRoot,
    env: {
      ...process.env,
      ...envFromFiles
    }
  });
}

if (dryRun) {
  for (const processConfig of processes) {
    console.log(
      `[dry-run] ${processConfig.name}: ${processConfig.command} ${processConfig.args.join(" ")}`
    );
  }
  process.exit(0);
}

if (processes.length === 0) {
  console.error("No runnable dev processes selected.");
  process.exit(1);
}

const children = [];
let shuttingDown = false;

function shutdown(signal = "SIGTERM") {
  if (shuttingDown) {
    return;
  }
  shuttingDown = true;

  for (const child of children) {
    if (!child.killed) {
      child.kill(signal);
    }
  }
}

process.on("SIGINT", () => {
  shutdown("SIGINT");
  process.exit(130);
});

process.on("SIGTERM", () => {
  shutdown("SIGTERM");
  process.exit(143);
});

for (const processConfig of processes) {
  const child = spawn(processConfig.command, processConfig.args, {
    cwd: processConfig.cwd,
    env: processConfig.env,
    stdio: "inherit"
  });

  child.on("exit", (code, signal) => {
    if (!shuttingDown) {
      shutdown();
      if (signal) {
        process.exit(1);
      }
      process.exit(code ?? 1);
    }
  });

  children.push(child);
}
