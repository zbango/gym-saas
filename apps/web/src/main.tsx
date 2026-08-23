import React, { useEffect, useState } from "react";
import ReactDOM from "react-dom/client";
import "./tailwind.css";
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
    <main className="min-h-screen bg-[radial-gradient(circle_at_top_left,var(--color-brand-ambient)_0%,var(--color-brand-surface)_44%,var(--color-brand-gold)_100%)] p-8 [font-family:Avenir_Next,Futura,sans-serif]">
      <div className="grid gap-5">
        <Panel title="Theme Studio" eyebrow={appName}>
          <div className="flex flex-wrap justify-between gap-4">
            <div>
              <p className="mt-0">Desktop and web now share the same theme engine.</p>
              <p className="mb-0 text-[var(--color-brand-muted)]">
                Zeus styling is compiled through Tailwind utilities.
              </p>
            </div>
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
    <App />
  </React.StrictMode>
);
