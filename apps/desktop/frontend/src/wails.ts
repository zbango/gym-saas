import { buildCloudApiUrl, type UpdateManifest } from "@gym-saas/shared";

export type DownloadedPackage = {
  fileName: string;
  path: string;
};

type DesktopBindings = {
  GetDesktopVersion(): Promise<string>;
  CheckForUpdates(manifestURL?: string): Promise<UpdateManifest | null>;
  OpenExternalURL(url: string): Promise<void>;
  DownloadUpdatePackage(url: string, checksum: string): Promise<DownloadedPackage>;
  InstallDownloadedPackage(path: string): Promise<void>;
  OpenPath(path: string): Promise<void>;
};

declare global {
  interface Window {
    go?: {
      main?: {
        App?: DesktopBindings;
      };
    };
  }
}

function app(): DesktopBindings | null {
  return window.go?.main?.App ?? null;
}

export function isDesktopApp(): boolean {
  return app() !== null;
}

export async function getDesktopVersion(): Promise<string | null> {
  const binding = app();
  if (!binding) {
    return null;
  }

  return binding.GetDesktopVersion();
}

export async function checkForUpdates(): Promise<UpdateManifest | null> {
  const binding = app();
  if (!binding) {
    return null;
  }

  const manifestUrl = buildCloudApiUrl(import.meta.env.VITE_CLOUD_API_BASE_URL, "/update-manifest");
  if (!manifestUrl) {
    return null;
  }

  return binding.CheckForUpdates(manifestUrl);
}

export async function openExternalURL(url: string): Promise<void> {
  const binding = app();
  if (!binding) {
    window.open(url, "_blank", "noopener,noreferrer");
    return;
  }

  await binding.OpenExternalURL(url);
}

export async function downloadUpdatePackage(url: string, checksum: string): Promise<DownloadedPackage | null> {
  const binding = app();
  if (!binding) {
    return null;
  }

  return binding.DownloadUpdatePackage(url, checksum);
}

export async function openPath(path: string): Promise<void> {
  const binding = app();
  if (!binding) {
    return;
  }

  await binding.OpenPath(path);
}

export async function installDownloadedPackage(path: string): Promise<void> {
  const binding = app();
  if (!binding) {
    return;
  }

  await binding.InstallDownloadedPackage(path);
}
