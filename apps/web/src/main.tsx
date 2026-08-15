import React, { useEffect, useState } from "react";
import ReactDOM from "react-dom/client";
import { appName, buildCloudApiUrl, platformLabel } from "@gym-saas/shared";
import { Panel } from "@gym-saas/ui";

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
          "radial-gradient(circle at top left, #f6dbc8 0%, #f0ece6 44%, #efe7db 100%)",
        fontFamily: "Avenir Next, Futura, sans-serif"
      }}
    >
      <Panel title="Web Host" eyebrow={appName}>
        <p>This is the browser host for gym-saas.</p>
        <p>Client platform: {platformLabel()}</p>
        <p>Cloud API URL: {cloudApiBaseUrl || "not configured"}</p>
        <p>Cloud version: {cloudVersion}</p>
        {cloudError ? <p>Cloud error: {cloudError}</p> : null}
      </Panel>
    </main>
  );
}

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
