import React, { useState } from "react";
import ReactDOM from "react-dom/client";
import "./tailwind.css";
import { LoginPage } from "./features/auth/LoginPage";
import { DashboardPage } from "./features/dashboard/DashboardPage";
import type { AuthenticatedUser } from "./features/auth/mockAuth";

function App() {
  const [user, setUser] = useState<AuthenticatedUser | null>(null);

  if (user) {
    return <DashboardPage user={user} onSignOut={() => setUser(null)} />;
  }

  return <LoginPage onAuthenticated={setUser} />;
}

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
