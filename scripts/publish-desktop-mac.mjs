import { createHash } from "node:crypto";
import { createReadStream, existsSync, mkdirSync } from "node:fs";
import { resolve } from "node:path";
import { spawnSync } from "node:child_process";
import { readEnvFile } from "./env-file.mjs";
import { readCurrentVersion } from "./versioning.mjs";

const repoRoot = process.cwd();
const envPath = resolve(repoRoot, ".env");
const dryRun = process.argv.includes("--dry-run");
const jsonOutput = process.argv.includes("--json");
const rootEnv = {
  ...readEnvFile(resolve(repoRoot, ".env.example")),
  ...readEnvFile(envPath)
};

const required = ["R2_AWS_PROFILE", "R2_ACCOUNT_ID", "R2_BUCKET_NAME", "R2_PUBLIC_BASE_URL"];
for (const key of required) {
  if (!valueFor(key)) {
    console.error(`Missing required release config: ${key}`);
    process.exit(1);
  }
}

const appBundlePath = resolve(repoRoot, "apps/desktop/build/bin/gym-saas desktop.app");
const releaseDir = resolve(repoRoot, "release");
const version = valueFor("RELEASE_VERSION") || readCurrentVersion();
const archiveName = `gym-saas-desktop-darwin-arm64-v${version}.zip`;
const archivePath = resolve(releaseDir, archiveName);
const objectKey = `releases/${archiveName}`;
const endpoint = `https://${valueFor("R2_ACCOUNT_ID")}.r2.cloudflarestorage.com`;
const publicUrl = `${valueFor("R2_PUBLIC_BASE_URL").replace(/\/+$/, "")}/${objectKey}`;

mkdirSync(releaseDir, { recursive: true });

if (!existsSync(appBundlePath)) {
  console.error("Missing desktop app bundle:", appBundlePath);
  console.error("Run npm run build:desktop:mac first.");
  process.exit(1);
}

run("ditto", ["-c", "-k", "--sequesterRsrc", "--keepParent", appBundlePath, archivePath], {
  dryRun
});

const checksum = dryRun ? "dry-run-checksum" : await sha256ForFile(archivePath);

run(
  "aws",
  [
    "s3",
    "cp",
    archivePath,
    `s3://${valueFor("R2_BUCKET_NAME")}/${objectKey}`,
    "--endpoint-url",
    endpoint
  ],
  { dryRun }
);

const payload = {
  archivePath,
  objectKey,
  publicUrl,
  checksum: `sha256:${checksum}`,
  envUrlKey: "RELEASE_DARWIN_ARM64_URL",
  envChecksumKey: "RELEASE_DARWIN_ARM64_CHECKSUM"
};

if (jsonOutput) {
  console.log(JSON.stringify(payload));
} else {
  console.log(`Archive: ${archivePath}`);
  console.log(`R2 object: s3://${valueFor("R2_BUCKET_NAME")}/${objectKey}`);
  console.log(`Public URL: ${publicUrl}`);
  console.log(`Checksum: sha256:${checksum}`);
  console.log("Use the release wizard or pass env vars into deploy to publish these values.");
}

function run(command, args, { dryRun }) {
  if (dryRun) {
    console.log(`[dry-run] ${command} ${args.join(" ")}`);
    return;
  }

  const result = spawnSync(command, args, {
    cwd: repoRoot,
    stdio: "inherit",
    env: {
      ...process.env,
      AWS_PROFILE: valueFor("R2_AWS_PROFILE")
    }
  });
  if (result.status !== 0) {
    process.exit(result.status ?? 1);
  }
}

async function sha256ForFile(path) {
  const hash = createHash("sha256");
  const stream = createReadStream(path);
  for await (const chunk of stream) {
    hash.update(chunk);
  }
  return hash.digest("hex");
}

function valueFor(key) {
  return process.env[key] || rootEnv[key];
}
