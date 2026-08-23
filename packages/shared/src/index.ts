export const appName = "gym-saas";
export const desktopVersion = "0.1.2";


export type UpdateAsset = {
  os: string;
  arch: string;
  url: string;
  checksum: string;
};

export type UpdateManifest = {
  version: string;
  minimumSupportedVersion: string;
  notes: string;
  publishedAt: string;
  assets: UpdateAsset[];
};

export function platformLabel(): string {
  if (typeof navigator === "undefined") {
    return "unknown";
  }

  return navigator.userAgent;
}

export function buildCloudApiUrl(baseUrl: string | undefined, path: string): string | null {
  const trimmedBase = baseUrl?.trim();
  if (!trimmedBase) {
    return null;
  }

  const normalizedBase = trimmedBase.replace(/\/+$/, "");
  const normalizedPath = path.startsWith("/") ? path : `/${path}`;

  return `${normalizedBase}${normalizedPath}`;
}
