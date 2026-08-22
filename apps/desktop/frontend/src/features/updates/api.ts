import { buildCloudApiUrl, type UpdateManifest } from "@gym-saas/shared";
import {
  CheckForUpdates,
  DownloadUpdatePackage,
  GetDesktopVersion,
  InstallDownloadedPackage,
  OpenExternalURL,
  OpenPath
} from "../../../wailsjs/go/main/DesktopAPI";
import type { updater } from "../../../wailsjs/go/models";

export type DownloadedPackage = updater.DownloadedPackage;

export async function getDesktopVersion(): Promise<string | null> {
  if (!isDesktopRuntime()) {
    return null;
  }
  return GetDesktopVersion();
}

export async function checkForUpdates(): Promise<UpdateManifest | null> {
  const manifestURL = buildCloudApiUrl(import.meta.env.VITE_CLOUD_API_BASE_URL, "/update-manifest");
  if (!manifestURL) {
    return null;
  }
  if (!isDesktopRuntime()) {
    return null;
  }
  return CheckForUpdates(manifestURL) as Promise<UpdateManifest>;
}

export async function openExternalURL(url: string): Promise<void> {
  if (!isDesktopRuntime()) {
    window.open(url, "_blank", "noopener,noreferrer");
    return;
  }
  await OpenExternalURL(url);
}

export async function downloadUpdatePackage(url: string, checksum: string): Promise<DownloadedPackage | null> {
  if (!isDesktopRuntime()) {
    return null;
  }
  return DownloadUpdatePackage(url, checksum);
}

export async function openPath(path: string): Promise<void> {
  if (isDesktopRuntime()) {
    await OpenPath(path);
  }
}

export async function installDownloadedPackage(path: string): Promise<void> {
  if (isDesktopRuntime()) {
    await InstallDownloadedPackage(path);
  }
}

function isDesktopRuntime(): boolean {
  return typeof window !== "undefined" && "go" in window;
}
