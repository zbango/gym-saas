import React, { useEffect, useState } from "react";
import ReactDOM from "react-dom/client";
import { appName, buildCloudApiUrl, platformLabel } from "@gym-saas/shared";
import { Panel, ThemeProvider, ThemeSwitcher } from "@gym-saas/ui";

function App() {
  const [cloudVersion, setCloudVersion] = useState<string>("not fetched");
  const [cloudError, setCloudError] = useState<string | null>(null);
  const cloudApiBaseUrl = import.meta.env.VITE_CLOUD_API_BASE_URL;
  const versionUrl = buildCloudApiUrl(cloudApiBaseUrl, "/version");

  useEffect(() => {
    if (!versionUrl) {
      setCloudError("cloud api url not configured");
      return;
    }

    let cancelled = false;

    void fetch(versionUrl)
      .then(async (response) => {
        if (!response.ok) {
          throw new Error(`cloud status ${response.status}`);
        }

        return (await response.json()) as { version?: string };
      })
      .then((payload) => {
        if (!cancelled) {
          setCloudVersion(payload.version ?? "unknown");
          setCloudError(null);
        }
      })
      .catch((error: Error) => {
        if (!cancelled) {
          setCloudError(error.message);
        }
      });

    return () => {
      cancelled = true;
    };
  }, [versionUrl]);

  return (
    <main
      style={{
        minHeight: "100vh",
        padding: 32,
        background:
          "radial-gradient(circle at top left, var(--gs-background-accent) 0%, var(--gs-surface) 44%, var(--gs-background-ambient) 100%)",
        fontFamily: "Avenir Next, Futura, sans-serif"
      }}
    >
      <div style={{ display: "grid", gap: 20 }}>
        <Panel title="Theme Studio" eyebrow={appName}>
          <div style={{ display: "flex", justifyContent: "space-between", gap: 16, flexWrap: "wrap" }}>
            <div>
              <p style={{ marginTop: 0 }}>Desktop and web now share the same theme engine.</p>
              <p style={{ marginBottom: 0, color: "var(--gs-text-muted)" }}>
                Pick a shell look once and let future business screens inherit it automatically.
              </p>
            </div>
            <ThemeSwitcher />
          </div>
        </Panel>

        <Panel title="Web Host" eyebrow={appName}>
          <p>This is the browser host for gym-saas.</p>
          <p>Client platform: {platformLabel()}</p>
          <p>Cloud API URL: {cloudApiBaseUrl || "not configured"}</p>
          <p>Cloud version: {cloudVersion}</p>
          {cloudError ? <p>Cloud error: {cloudError}</p> : null}
        </Panel>
      </div>
    </main>
  );
}

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <ThemeProvider>
      <App />
    </ThemeProvider>
  </React.StrictMode>
);
