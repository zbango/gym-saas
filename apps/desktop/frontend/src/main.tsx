import React, { useState } from "react";
import ReactDOM from "react-dom/client";
import { ThemeProvider } from "@gym-saas/ui";
import { LoginPage } from "./features/auth/LoginPage";
import { DashboardPage } from "./features/dashboard/DashboardPage";

function App() {
  const [screen, setScreen] = useState<"login" | "dashboard">("login");

  if (screen === "dashboard") {
    return <DashboardPage onSignOut={() => setScreen("login")} />;
  }

  return <LoginPage onSubmit={() => setScreen("dashboard")} />;
}

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <ThemeProvider fixedTheme="zeus">
      <App />
    </ThemeProvider>
  </React.StrictMode>
);
