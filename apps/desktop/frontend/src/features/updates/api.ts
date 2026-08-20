import { buildCloudApiUrl, type UpdateManifest } from "@gym-saas/shared";
import { desktopApp } from "../../platform/wails";

export type DownloadedPackage = {
  fileName: string;
  path: string;
};

type UpdateBindings = {
  GetDesktopVersion(): Promise<string>;
  CheckForUpdates(manifestURL?: string): Promise<UpdateManifest | null>;
  OpenExternalURL(url: string): Promise<void>;
  DownloadUpdatePackage(url: string, checksum: string): Promise<DownloadedPackage>;
  InstallDownloadedPackage(path: string): Promise<void>;
  OpenPath(path: string): Promise<void>;
};

function updates(): UpdateBindings | null {
  return desktopApp() as UpdateBindings | null;
}

export async function getDesktopVersion(): Promise<string | null> {
  return updates()?.GetDesktopVersion() ?? null;
}

export async function checkForUpdates(): Promise<UpdateManifest | null> {
  const manifestURL = buildCloudApiUrl(import.meta.env.VITE_CLOUD_API_BASE_URL, "/update-manifest");
  if (!manifestURL) {
    return null;
  }
  return updates()?.CheckForUpdates(manifestURL) ?? null;
}

export async function openExternalURL(url: string): Promise<void> {
  const bridge = updates();
  if (!bridge) {
    window.open(url, "_blank", "noopener,noreferrer");
    return;
  }
  await bridge.OpenExternalURL(url);
}

export async function downloadUpdatePackage(url: string, checksum: string): Promise<DownloadedPackage | null> {
  return updates()?.DownloadUpdatePackage(url, checksum) ?? null;
}

export async function openPath(path: string): Promise<void> {
  await updates()?.OpenPath(path);
}

export async function installDownloadedPackage(path: string): Promise<void> {
  await updates()?.InstallDownloadedPackage(path);
}
