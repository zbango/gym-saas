import type { CSSProperties, PropsWithChildren } from "react";
import { createContext, useContext, useEffect, useMemo, useState } from "react";
import {
  defaultThemeName,
  getThemeDefinition,
  listThemes,
  themeStorageKey,
  type ThemeDefinition,
  type ThemeName
} from "@gym-saas/shared";

type ThemeContextValue = {
  theme: ThemeDefinition;
  themeName: ThemeName;
  themes: ThemeDefinition[];
  setThemeName: (name: ThemeName) => void;
};

const ThemeContext = createContext<ThemeContextValue | null>(null);

export function ThemeProvider(props: PropsWithChildren<{ initialTheme?: ThemeName; fixedTheme?: ThemeName }>) {
  const [themeName, setThemeName] = useState<ThemeName>(() => {
    if (typeof window === "undefined") {
      return props.fixedTheme ?? props.initialTheme ?? defaultThemeName;
    }

    const stored = window.localStorage.getItem(themeStorageKey);
    return getThemeDefinition(props.fixedTheme ?? stored ?? props.initialTheme ?? defaultThemeName).name;
  });

  const activeThemeName = props.fixedTheme ?? themeName;

  useEffect(() => {
    if (typeof window === "undefined") {
      return;
    }

    window.localStorage.setItem(themeStorageKey, activeThemeName);
  }, [activeThemeName]);

  const value = useMemo<ThemeContextValue>(() => {
    const theme = getThemeDefinition(activeThemeName);
    return {
      theme,
      themeName: theme.name,
      themes: listThemes(),
      setThemeName
    };
  }, [activeThemeName]);

  return (
    <ThemeContext.Provider value={value}>
      <div data-theme={value.themeName} style={themeVariables(value.theme)}>
        {props.children}
      </div>
    </ThemeContext.Provider>
  );
}

export function useTheme() {
  const value = useContext(ThemeContext);
  if (!value) {
    throw new Error("useTheme must be used inside ThemeProvider");
  }

  return value;
}

export function ThemeSwitcher() {
  const { themeName, themes, setThemeName } = useTheme();

  return (
    <label
      style={{
        display: "grid",
        gap: 6,
        minWidth: 220
      }}
    >
      <span
        style={{
          color: "var(--gs-text-muted)",
          fontSize: 12,
          fontWeight: 700,
          letterSpacing: "0.12em",
          textTransform: "uppercase"
        }}
      >
        Theme
      </span>
      <select
        value={themeName}
        onChange={(event) => setThemeName(event.target.value as ThemeName)}
        style={{
          background: "var(--gs-input-background)",
          color: "var(--gs-input-text)",
          border: "1px solid var(--gs-border)",
          borderRadius: 12,
          padding: "10px 12px",
          font: "inherit"
        }}
      >
        {themes.map((theme) => (
          <option key={theme.name} value={theme.name}>
            {theme.label}
          </option>
        ))}
      </select>
    </label>
  );
}

function themeVariables(theme: ThemeDefinition): CSSProperties {
  const tokens = theme.tokens;
  return {
    "--gs-background": tokens.background,
    "--gs-background-accent": tokens.backgroundAccent,
    "--gs-background-ambient": tokens.backgroundAmbient,
    "--gs-workspace-background": tokens.workspaceBackground,
    "--gs-sidebar": tokens.sidebar,
    "--gs-sidebar-text": tokens.sidebarText,
    "--gs-surface": tokens.surface,
    "--gs-surface-elevated": tokens.surfaceElevated,
    "--gs-border": tokens.border,
    "--gs-border-strong": tokens.borderStrong,
    "--gs-text": tokens.text,
    "--gs-text-muted": tokens.textMuted,
    "--gs-accent": tokens.accent,
    "--gs-accent-strong": tokens.accentStrong,
    "--gs-accent-soft": tokens.accentSoft,
    "--gs-shadow": tokens.shadow,
    "--gs-input-background": tokens.inputBackground,
    "--gs-input-text": tokens.inputText,
    "--gs-button-primary-background": tokens.buttonPrimaryBackground,
    "--gs-button-primary-text": tokens.buttonPrimaryText,
    "--gs-button-secondary-background": tokens.buttonSecondaryBackground,
    "--gs-button-secondary-text": tokens.buttonSecondaryText,
    "--gs-overlay": tokens.overlay,
    "--gs-success": tokens.success,
    "--gs-success-soft": tokens.successSoft,
    "--gs-danger": tokens.danger,
    "--gs-danger-soft": tokens.dangerSoft
  } as CSSProperties;
}
