import React, { useEffect, useState } from "react";
import ReactDOM from "react-dom/client";
import {
  appName,
  desktopVersion,
  type HelloRecord,
  type UpdateAsset,
  type UpdateManifest
} from "@gym-saas/shared";
import { Panel } from "@gym-saas/ui";
import {
  checkForUpdates,
  createHelloRecord,
  deleteHelloRecord,
  downloadUpdatePackage,
  getDesktopVersion,
  installDownloadedPackage,
  isDesktopApp,
  listHelloRecords,
  openPath,
  openExternalURL,
  type DownloadedPackage,
  updateHelloRecord
} from "./wails";

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
  const [name, setName] = useState("");
  const [records, setRecords] = useState<HelloRecord[]>([]);
  const [update, setUpdate] = useState<UpdateManifest | null>(null);
  const [runtimeVersion, setRuntimeVersion] = useState<string>(desktopVersion);
  const [updateState, setUpdateState] = useState<UpdateState>("idle");
  const [updateError, setUpdateError] = useState<string | null>(null);
  const [showUpdateModal, setShowUpdateModal] = useState(false);
  const [downloadedPackage, setDownloadedPackage] = useState<DownloadedPackage | null>(null);
  const desktopRuntime = isDesktopApp();

  const releaseAsset: UpdateAsset | null = update?.assets[0] ?? null;

  async function refresh() {
    setRecords(await listHelloRecords());
  }

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
    void refresh();
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
          "linear-gradient(135deg, #efe3c6 0%, #f8f4ee 45%, #e9eef0 100%)",
        fontFamily: "Avenir Next, Futura, sans-serif",
        display: "grid",
        gap: 20
      }}
    >
      <Panel title="Desktop Host" eyebrow={`${appName} ${desktopVersion}`}>
        <p>Wails desktop shell with embedded SQLite and updater check.</p>
        <p>Runtime version: {runtimeVersion}</p>
        <div style={{ display: "flex", gap: 8 }}>
          <input
            value={name}
            onChange={(event) => setName(event.target.value)}
            placeholder="Hello record name"
            style={{ padding: 10, flex: 1 }}
          />
          <button
            onClick={async () => {
              if (!name.trim()) return;
              await createHelloRecord(name.trim());
              setName("");
              await refresh();
            }}
          >
            Create
          </button>
          <button
            onClick={async () => {
              await runUpdateCheck(true);
            }}
            disabled={updateState === "checking"}
          >
            {updateState === "checking" ? "Checking..." : "Check updates"}
          </button>
        </div>
      </Panel>

      <Panel title="SQLite Hello Records" eyebrow="Embedded DB">
        {records.length === 0 ? <p>No records yet.</p> : null}
        <ul style={{ listStyle: "none", padding: 0, margin: 0, display: "grid", gap: 12 }}>
          {records.map((record) => (
            <li
              key={record.id}
              style={{
                display: "flex",
                justifyContent: "space-between",
                gap: 12,
                alignItems: "center"
              }}
            >
              <div>
                <strong>{record.name}</strong>
                <div style={{ fontSize: 12, color: "#745d3f" }}>{record.updatedAt}</div>
              </div>
              <div style={{ display: "flex", gap: 8 }}>
                <button
                  onClick={async () => {
                    await updateHelloRecord(record.id, `${record.name} updated`);
                    await refresh();
                  }}
                >
                  Update
                </button>
                <button
                  onClick={async () => {
                    await deleteHelloRecord(record.id);
                    await refresh();
                  }}
                >
                  Delete
                </button>
              </div>
            </li>
          ))}
        </ul>
      </Panel>

      <Panel title="Updater" eyebrow="Cross-OS manifest">
        <p>{updaterDescription}</p>
        {update ? <p>Release notes: {update.notes}</p> : null}
        {releaseAsset ? (
          <div style={{ display: "flex", gap: 8, flexWrap: "wrap" }}>
            <button
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
            </button>
            <button onClick={() => setShowUpdateModal(true)}>Show update details</button>
          </div>
        ) : null}
        {downloadedPackage ? (
          <div style={{ display: "flex", gap: 8, flexWrap: "wrap" }}>
            <button
              onClick={async () => {
                await openPath(downloadedPackage.path);
              }}
            >
              Open downloaded package
            </button>
            <button
              onClick={async () => {
                await openExternalURL(releaseAsset?.url ?? "");
              }}
            >
              Open release URL
            </button>
          </div>
        ) : null}
      </Panel>

      {showUpdateModal && (updateState === "available" || updateState === "downloaded") ? (
        <div
          style={{
            position: "fixed",
            inset: 0,
            background: "rgba(38, 23, 8, 0.45)",
            display: "grid",
            placeItems: "center",
            padding: 24
          }}
        >
          <div
            style={{
              width: "min(560px, 100%)",
              borderRadius: 20,
              background: "#fff8ef",
              border: "1px solid #d8c4a7",
              boxShadow: "0 24px 60px rgba(47, 27, 8, 0.24)",
              padding: 24,
              display: "grid",
              gap: 14
            }}
          >
            <div>
              <div
                style={{
                  fontSize: 12,
                  fontWeight: 700,
                  letterSpacing: "0.14em",
                  textTransform: "uppercase",
                  color: "#8a5a18",
                  marginBottom: 8
                }}
              >
                Update available
              </div>
              <h2 style={{ margin: 0, color: "#2d1b05" }}>
                gym-saas {update?.version ?? "unknown"} is ready
              </h2>
            </div>

            <p style={{ margin: 0, color: "#4c3923" }}>
              A new desktop version was published. After download, the app hands off to the
              installer/update flow immediately.
            </p>

            <div
              style={{
                borderRadius: 14,
                padding: 16,
                background: "#f5eadc",
                color: "#3b2a17"
              }}
            >
              <strong>Release notes</strong>
              <p style={{ margin: "8px 0 0" }}>{update?.notes}</p>
            </div>

            <div style={{ fontSize: 13, color: "#6f5634" }}>
              <div>Current version: {desktopVersion}</div>
              <div>Runtime version: {runtimeVersion}</div>
              <div>Published version: {update?.version ?? "unknown"}</div>
              <div>Cloud URL: {import.meta.env.VITE_CLOUD_API_BASE_URL || "not configured"}</div>
              {downloadedPackage ? <div>Downloaded package: {downloadedPackage.path}</div> : null}
            </div>

            <div style={{ display: "flex", gap: 10, justifyContent: "flex-end", flexWrap: "wrap" }}>
              <button onClick={() => setShowUpdateModal(false)}>Later</button>
              <button
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
              </button>
            </div>
          </div>
        </div>
      ) : null}

      {showUpdateModal && updateState === "error" ? (
        <div
          style={{
            position: "fixed",
            inset: 0,
            background: "rgba(38, 23, 8, 0.45)",
            display: "grid",
            placeItems: "center",
            padding: 24
          }}
        >
          <div
            style={{
              width: "min(520px, 100%)",
              borderRadius: 20,
              background: "#fff8ef",
              border: "1px solid #d8c4a7",
              boxShadow: "0 24px 60px rgba(47, 27, 8, 0.24)",
              padding: 24,
              display: "grid",
              gap: 14
            }}
          >
            <h2 style={{ margin: 0, color: "#2d1b05" }}>Update check failed</h2>
            <p style={{ margin: 0, color: "#4c3923" }}>{updateError}</p>
            <div style={{ display: "flex", gap: 10, justifyContent: "flex-end" }}>
              <button onClick={() => setShowUpdateModal(false)}>Close</button>
              <button
                onClick={async () => {
                  await runUpdateCheck(true);
                }}
              >
                Retry
              </button>
            </div>
          </div>
        </div>
      ) : null}
    </main>
  );
}

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
