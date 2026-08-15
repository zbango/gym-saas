import { spawnSync } from "node:child_process";

const target = process.argv[2];

function run(command, args, options = {}) {
  const result = spawnSync(command, args, {
    stdio: "inherit",
    ...options
  });

  if (result.status !== 0) {
    process.exit(result.status ?? 1);
  }
}

run("npm", ["run", "typecheck"]);
run("npm", ["run", "build:web"]);
run("npm", ["run", "build:cloud-api"]);
run("npm", ["run", "build:desktop", ...(target ? ["--", target] : [])]);
