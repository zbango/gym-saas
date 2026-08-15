import { existsSync } from "node:fs";
import { spawnSync } from "node:child_process";
import { resolve } from "node:path";
import { readEnvFile } from "./env-file.mjs";

const repoRoot = process.cwd();
const cloudDir = resolve(repoRoot, "apps/cloud-api");
const dryRun = process.argv.includes("--dry-run");
const guided = process.argv.includes("--guided");
const hasConfig = existsSync(resolve(cloudDir, "samconfig.toml"));
const configEnv = "main";
const rootEnv = {
  ...readEnvFile(resolve(repoRoot, ".env.example")),
  ...readEnvFile(resolve(repoRoot, ".env"))
};
const parameterOverrides = buildParameterOverrides();

const steps = [
  ["sam", ["build", "--template-file", "template.yaml"]],
  [
    "sam",
    [
      "deploy",
      "--config-env",
      configEnv,
      ...(parameterOverrides.length > 0 ? ["--parameter-overrides", ...parameterOverrides] : []),
      ...(guided || !hasConfig ? ["--guided"] : [])
    ]
  ]
];

if (dryRun) {
  for (const [command, args] of steps) {
    console.log(`[dry-run] (cwd=${cloudDir}) ${command} ${args.join(" ")}`);
  }
  process.exit(0);
}

for (const [command, args] of steps) {
  const result = spawnSync(command, args, {
    cwd: cloudDir,
    stdio: "inherit",
    env: {
      ...process.env,
      AWS_PROFILE:
        process.env.AWS_DEPLOY_PROFILE || rootEnv.AWS_DEPLOY_PROFILE || process.env.AWS_PROFILE
    }
  });

  if (result.error?.code === "ENOENT") {
    console.error("Missing required command:", command);
    console.error("Install AWS SAM CLI before running deploy.");
    process.exit(1);
  }

  if (result.status !== 0) {
    process.exit(result.status ?? 1);
  }
}

function buildParameterOverrides(env) {
  const mappings = [
    ["RELEASE_VERSION", "ReleaseVersion"],
    ["RELEASE_MINIMUM_SUPPORTED_VERSION", "ReleaseMinimumSupportedVersion"],
    ["RELEASE_NOTES", "ReleaseNotes"],
    ["RELEASE_DARWIN_ARM64_URL", "DarwinArm64Url"],
    ["RELEASE_DARWIN_ARM64_CHECKSUM", "DarwinArm64Checksum"],
    ["RELEASE_WINDOWS_AMD64_URL", "WindowsAmd64Url"],
    ["RELEASE_WINDOWS_AMD64_CHECKSUM", "WindowsAmd64Checksum"],
    ["RELEASE_LINUX_AMD64_URL", "LinuxAmd64Url"],
    ["RELEASE_LINUX_AMD64_CHECKSUM", "LinuxAmd64Checksum"]
  ];

  return mappings.flatMap(([envKey, parameterKey]) => {
    const value = process.env[envKey] || rootEnv[envKey];
    if (!value) {
      return [];
    }

    return [`${parameterKey}=${value}`];
  });
}
