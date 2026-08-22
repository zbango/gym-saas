import React from "react";
import ReactDOM from "react-dom/client";
import { ThemeProvider } from "@gym-saas/ui";
import { LoginPage } from "./features/auth/LoginPage";

function App() {
  return <LoginPage />;
}

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <ThemeProvider fixedTheme="zeus">
      <App />
    </ThemeProvider>
  </React.StrictMode>
);
