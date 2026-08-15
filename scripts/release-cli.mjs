import { existsSync } from "node:fs";
import { spawnSync } from "node:child_process";
import { resolve } from "node:path";
import { createInterface } from "node:readline/promises";
import { stdin as input, stdout as output } from "node:process";
import { readEnvFile } from "./env-file.mjs";
import { bumpPatch, readCurrentVersion, syncVersionFiles } from "./versioning.mjs";

const repoRoot = process.cwd();
const envPath = resolve(repoRoot, ".env");
const dryRun = process.argv.includes("--dry-run");

const defaults = {
  ...readEnvFile(resolve(repoRoot, ".env.example")),
  ...readEnvFile(envPath)
};

const targets = [
  {
    key: "darwin-arm64",
    label: "macOS (Apple Silicon)",
    envUrl: "RELEASE_DARWIN_ARM64_URL",
    envChecksum: "RELEASE_DARWIN_ARM64_CHECKSUM",
    buildScript: "build:desktop:mac",
    publishScript: "publish:desktop:mac",
    available: true
  },
  {
    key: "windows-amd64",
    label: "Windows x64",
    envUrl: "RELEASE_WINDOWS_AMD64_URL",
    envChecksum: "RELEASE_WINDOWS_AMD64_CHECKSUM",
    buildScript: "build:desktop:windows",
    publishScript: null,
    available: false
  },
  {
    key: "linux-amd64",
    label: "Linux x64",
    envUrl: "RELEASE_LINUX_AMD64_URL",
    envChecksum: "RELEASE_LINUX_AMD64_CHECKSUM",
    buildScript: "build:desktop:linux",
    publishScript: null,
    available: false
  }
];

const rl = createInterface({ input, output });
let rlClosed = false;

async function main() {
  console.log("gym-saas release wizard");
  console.log("");

  const target = await chooseTarget();
  const currentVersion = readCurrentVersion();
  const releaseVersion = bumpPatch(currentVersion);
  const minimumSupportedVersion = currentVersion;
  const releaseNotes = await ask(
    `Release notes for ${releaseVersion}`,
    defaults.RELEASE_NOTES || "Release notes"
  );

  const nextEnv = {
    ...defaults,
    RELEASE_VERSION: releaseVersion,
    RELEASE_MINIMUM_SUPPORTED_VERSION: minimumSupportedVersion,
    RELEASE_NOTES: releaseNotes
  };

  ensureReleaseInfra(nextEnv);

  console.log("");
  console.log("Target:", target.label);
  console.log("Current version:", currentVersion);
  console.log("Version:", releaseVersion);
  console.log("Minimum supported:", minimumSupportedVersion);
  console.log("Notes:", releaseNotes);
  console.log("");

  const proceed = await confirm("Apply version bump and continue?", true);
  if (!proceed) {
    console.log("Release cancelled.");
    return;
  }

  if (dryRun) {
    console.log(`[dry-run] VERSION would change: ${currentVersion} -> ${releaseVersion}`);
  } else {
    syncVersionFiles(releaseVersion);
  }

  const runBuild = await confirm("Run desktop build?", true);
  const runPublish = target.available ? await confirm("Publish release artifact?", true) : false;
  const runDeploy = await confirm("Deploy cloud manifest?", true);

  console.log("");
  closeWizardInput();

  if (runBuild) {
    runNpmScript(target.buildScript);
  }

  if (!target.available) {
    console.log(
      `Release publishing for ${target.label} is not implemented yet. Build/deploy scaffolding remains available.`
    );
  } else if (runPublish && target.publishScript) {
    const published = runPublishScript(target.publishScript, nextEnv);
    nextEnv[target.envUrl] = published.publicUrl;
    nextEnv[target.envChecksum] = published.checksum;
  }

  if (runDeploy) {
    runNpmScript("deploy", nextEnv);
  }

  console.log("");
  console.log("Release wizard complete.");
}
 
main()
  .catch((error) => {
    console.error(error);
    process.exitCode = 1;
  })
  .finally(() => {
    closeWizardInput();
  });

async function chooseTarget() {
  console.log("Choose target:");
  targets.forEach((target, index) => {
    const suffix = target.available ? "" : " (publishing not implemented yet)";
    console.log(`  ${index + 1}. ${target.label}${suffix}`);
  });

  while (true) {
    const answer = await ask("Target number", "1");
    const index = Number.parseInt(answer, 10) - 1;
    if (index >= 0 && index < targets.length) {
      return targets[index];
    }
    console.log("Choose a valid target number.");
  }
}

async function ask(label, defaultValue) {
  const suffix = defaultValue ? ` [${defaultValue}]` : "";
  const answer = (await rl.question(`${label}${suffix}: `)).trim();
  return answer || defaultValue;
}

async function confirm(label, defaultYes) {
  const suffix = defaultYes ? " [Y/n]" : " [y/N]";
  const answer = (await rl.question(`${label}${suffix}: `)).trim().toLowerCase();
  if (!answer) {
    return defaultYes;
  }
  return answer === "y" || answer === "yes";
}

function ensureReleaseInfra(env) {
  const required = [
    "R2_AWS_PROFILE",
    "AWS_DEPLOY_PROFILE",
    "R2_ACCOUNT_ID",
    "R2_BUCKET_NAME",
    "R2_PUBLIC_BASE_URL"
  ];
  const missing = required.filter((key) => !env[key]);

  if (missing.length > 0) {
    console.log("Missing release infrastructure config:");
    for (const key of missing) {
      console.log(`  - ${key}`);
    }
    console.log("Fill those values in .env before publishing.");
  }
}

function runNpmScript(scriptName, env = {}) {
  if (dryRun) {
    console.log(`[dry-run] npm run ${scriptName}`);
    return;
  }

  const result = spawnSync("npm", ["run", scriptName], {
    cwd: repoRoot,
    stdio: "inherit",
    env: {
      ...process.env,
      ...env
    }
  });

  if (result.status !== 0) {
    process.exit(result.status ?? 1);
  }
}

function runPublishScript(scriptName, env) {
  if (dryRun) {
    console.log(`[dry-run] npm run ${scriptName} -- --json`);
    return {
      publicUrl: "https://example.com/dry-run.zip",
      checksum: "sha256:dry-run-checksum"
    };
  }

  const result = spawnSync("npm", ["run", scriptName, "--", "--json"], {
    cwd: repoRoot,
    encoding: "utf8",
    env: {
      ...process.env,
      ...env
    }
  });

  if (result.status !== 0) {
    process.stdout.write(result.stdout || "");
    process.stderr.write(result.stderr || "");
    process.exit(result.status ?? 1);
  }

  const lines = (result.stdout || "")
    .split(/\r?\n/)
    .map((line) => line.trim())
    .filter(Boolean);
  const payload = JSON.parse(lines.at(-1));

  if (result.stderr) {
    process.stderr.write(result.stderr);
  }
  if (result.stdout) {
    for (const line of lines.slice(0, -1)) {
      console.log(line);
    }
  }

  return payload;
}

function closeWizardInput() {
  if (rlClosed) {
    return;
  }

  rlClosed = true;
  rl.close();
}
