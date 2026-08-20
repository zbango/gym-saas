import React, { useEffect, useState } from "react";
import type { PropsWithChildren } from "react";
import ReactDOM from "react-dom/client";
import {
  appName,
  desktopVersion,
  type UpdateAsset,
  type UpdateManifest
} from "@gym-saas/shared";
import { Panel, ShellButton, ThemeProvider, ThemeSwitcher } from "@gym-saas/ui";
import {
  checkForUpdates,
  downloadUpdatePackage,
  getDesktopVersion,
  installDownloadedPackage,
  openPath,
  openExternalURL,
  type DownloadedPackage
} from "./features/updates/api";
import { MemberPanel } from "./features/members/MemberPanel";
import { isDesktopApp } from "./platform/wails";

type UpdateState =
  | "idle"
  | "checking"
  | "available"
  | "current"
  | "unsupported"
  | "downloading"
  | "installing"
  | "downloaded"
  | "error";

function App() {
  const [update, setUpdate] = useState<UpdateManifest | null>(null);
  const [runtimeVersion, setRuntimeVersion] = useState<string>(desktopVersion);
  const [updateState, setUpdateState] = useState<UpdateState>("idle");
  const [updateError, setUpdateError] = useState<string | null>(null);
  const [showUpdateModal, setShowUpdateModal] = useState(false);
  const [downloadedPackage, setDownloadedPackage] = useState<DownloadedPackage | null>(null);
  const desktopRuntime = isDesktopApp();

  const releaseAsset: UpdateAsset | null = update?.assets[0] ?? null;

  async function runUpdateCheck(userInitiated = false) {
    if (!desktopRuntime) {
      setUpdateState("unsupported");
      setUpdateError(null);
      return;
    }

    setUpdateState("checking");
    setUpdateError(null);
    setDownloadedPackage(null);

    try {
      const manifest = await checkForUpdates();
      if (manifest) {
        setUpdate(manifest);
        setUpdateState("available");
        setShowUpdateModal(true);
        return;
      }

      setUpdate(null);
      setUpdateState("current");
      if (userInitiated) {
        setShowUpdateModal(false);
      }
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown updater error";
      setUpdate(null);
      setUpdateState("error");
      setUpdateError(message);
      setShowUpdateModal(userInitiated);
    }
  }

  useEffect(() => {
    void runUpdateCheck();
    void (async () => {
      const value = await getDesktopVersion();
      if (value) {
        setRuntimeVersion(value);
      }
    })();
  }, []);

  const updaterDescription = (() => {
    switch (updateState) {
      case "checking":
        return "Checking for updates...";
      case "available":
        return `Update ${update?.version ?? ""} is available.`;
      case "downloading":
        return "Downloading update package...";
      case "installing":
        return "Installing update package...";
      case "downloaded":
        return `Update package downloaded: ${downloadedPackage?.fileName ?? "unknown file"}`;
      case "current":
        return `You are on the latest version (${runtimeVersion}).`;
      case "unsupported":
        return "Update checks only work inside the desktop app.";
      case "error":
        return `Update check failed: ${updateError}`;
      default:
        return `Cloud URL: ${import.meta.env.VITE_CLOUD_API_BASE_URL || "not configured"}`;
    }
  })();

  return (
    <main
      style={{
        minHeight: "100vh",
        padding: 24,
        background:
          "linear-gradient(135deg, var(--gs-background) 0%, var(--gs-background-accent) 45%, var(--gs-background-ambient) 100%)",
        fontFamily: "Avenir Next, Futura, sans-serif",
        display: "grid",
        gap: 20
      }}
    >
      <Panel title="Theme Studio" eyebrow={appName}>
        <div style={{ display: "flex", justifyContent: "space-between", gap: 16, flexWrap: "wrap" }}>
          <div>
            <p style={{ marginTop: 0 }}>Desktop and web share the same theme engine and persistence model.</p>
            <p style={{ marginBottom: 0, color: "var(--gs-text-muted)" }}>
              Change the shell look now so future business screens inherit it automatically.
            </p>
          </div>
          <ThemeSwitcher />
        </div>
      </Panel>

      <MemberPanel />

      <Panel title="Desktop Host" eyebrow={`${appName} ${desktopVersion}`}>
        <p>Wails desktop shell with embedded SQLite migration support and updater check.</p>
        <p>Runtime version: {runtimeVersion}</p>
        <div style={{ display: "flex", gap: 8, flexWrap: "wrap" }}>
          <ShellButton
            variant="secondary"
            onClick={async () => {
              await runUpdateCheck(true);
            }}
            disabled={updateState === "checking"}
          >
            {updateState === "checking" ? "Checking..." : "Check updates"}
          </ShellButton>
        </div>
      </Panel>

      <Panel title="Updater" eyebrow="Cross-OS manifest">
        <p>{updaterDescription}</p>
        {update ? <p>Release notes: {update.notes}</p> : null}
        {releaseAsset ? (
          <div style={{ display: "flex", gap: 8, flexWrap: "wrap" }}>
            <ShellButton
              onClick={async () => {
                setUpdateState("downloading");
                setUpdateError(null);
                try {
                  const result = await downloadUpdatePackage(releaseAsset.url, releaseAsset.checksum);
                  if (!result) {
                    throw new Error("download is only available in the desktop app");
                  }
                  setDownloadedPackage(result);
                  setUpdateState("installing");
                  setShowUpdateModal(true);
                  await installDownloadedPackage(result.path);
                  setUpdateState("downloaded");
                } catch (error) {
                  const message = error instanceof Error ? error.message : "Unknown download error";
                  setUpdateState("error");
                  setUpdateError(message);
                  setShowUpdateModal(true);
                }
              }}
              disabled={updateState === "downloading" || updateState === "installing"}
            >
              {updateState === "downloading"
                ? "Downloading..."
                : updateState === "installing"
                  ? "Installing..."
                  : "Download update"}
            </ShellButton>
            <ShellButton variant="secondary" onClick={() => setShowUpdateModal(true)}>
              Show update details
            </ShellButton>
          </div>
        ) : null}
        {downloadedPackage ? (
          <div style={{ display: "flex", gap: 8, flexWrap: "wrap" }}>
            <ShellButton
              variant="secondary"
              onClick={async () => {
                await openPath(downloadedPackage.path);
              }}
            >
              Open downloaded package
            </ShellButton>
            <ShellButton
              variant="secondary"
              onClick={async () => {
                await openExternalURL(releaseAsset?.url ?? "");
              }}
            >
              Open release URL
            </ShellButton>
          </div>
        ) : null}
      </Panel>

      {showUpdateModal && (updateState === "available" || updateState === "downloaded") ? (
        <Modal>
          <div>
            <div
              style={{
                fontSize: 12,
                fontWeight: 700,
                letterSpacing: "0.14em",
                textTransform: "uppercase",
                color: "var(--gs-accent)",
                marginBottom: 8
              }}
            >
              Update available
            </div>
            <h2 style={{ margin: 0, color: "var(--gs-text)" }}>
              gym-saas {update?.version ?? "unknown"} is ready
            </h2>
          </div>

          <p style={{ margin: 0, color: "var(--gs-text)" }}>
            A new desktop version was published. After download, the app hands off to the installer/update flow immediately.
          </p>

          <div
            style={{
              borderRadius: 14,
              padding: 16,
              background: "var(--gs-accent-soft)",
              color: "var(--gs-text)"
            }}
          >
            <strong>Release notes</strong>
            <p style={{ margin: "8px 0 0" }}>{update?.notes}</p>
          </div>

          <div style={{ fontSize: 13, color: "var(--gs-text-muted)" }}>
            <div>Current version: {desktopVersion}</div>
            <div>Runtime version: {runtimeVersion}</div>
            <div>Published version: {update?.version ?? "unknown"}</div>
            <div>Cloud URL: {import.meta.env.VITE_CLOUD_API_BASE_URL || "not configured"}</div>
            {downloadedPackage ? <div>Downloaded package: {downloadedPackage.path}</div> : null}
          </div>

          <div style={{ display: "flex", gap: 10, justifyContent: "flex-end", flexWrap: "wrap" }}>
            <ShellButton variant="secondary" onClick={() => setShowUpdateModal(false)}>
              Later
            </ShellButton>
            <ShellButton
              onClick={async () => {
                if (downloadedPackage) {
                  setUpdateState("installing");
                  await installDownloadedPackage(downloadedPackage.path);
                  return;
                }
                if (!releaseAsset) {
                  return;
                }
                setUpdateState("downloading");
                setUpdateError(null);
                try {
                  const result = await downloadUpdatePackage(releaseAsset.url, releaseAsset.checksum);
                  if (!result) {
                    throw new Error("download is only available in the desktop app");
                  }
                  setDownloadedPackage(result);
                  setUpdateState("installing");
                  await installDownloadedPackage(result.path);
                  setUpdateState("downloaded");
                } catch (error) {
                  const message = error instanceof Error ? error.message : "Unknown download error";
                  setUpdateState("error");
                  setUpdateError(message);
                }
              }}
              disabled={!releaseAsset}
            >
              {downloadedPackage
                ? "Install downloaded package"
                : releaseAsset
                  ? "Download and install update"
                  : "No package for this platform"}
            </ShellButton>
          </div>
        </Modal>
      ) : null}

      {showUpdateModal && updateState === "error" ? (
        <Modal>
          <h2 style={{ margin: 0, color: "var(--gs-text)" }}>Update check failed</h2>
          <p style={{ margin: 0, color: "var(--gs-text)" }}>{updateError}</p>
          <div style={{ display: "flex", gap: 10, justifyContent: "flex-end" }}>
            <ShellButton variant="secondary" onClick={() => setShowUpdateModal(false)}>
              Close
            </ShellButton>
            <ShellButton
              onClick={async () => {
                await runUpdateCheck(true);
              }}
            >
              Retry
            </ShellButton>
          </div>
        </Modal>
      ) : null}
    </main>
  );
}

function Modal(props: PropsWithChildren) {
  return (
    <div
      style={{
        position: "fixed",
        inset: 0,
        background: "var(--gs-overlay)",
        display: "grid",
        placeItems: "center",
        padding: 24
      }}
    >
      <div
        style={{
          width: "min(560px, 100%)",
          borderRadius: 20,
          background: "var(--gs-surface)",
          border: "1px solid var(--gs-border)",
          boxShadow: "var(--gs-shadow)",
          padding: 24,
          display: "grid",
          gap: 14,
          color: "var(--gs-text)"
        }}
      >
        {props.children}
      </div>
    </div>
  );
}

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <ThemeProvider>
      <App />
    </ThemeProvider>
  </React.StrictMode>
);
