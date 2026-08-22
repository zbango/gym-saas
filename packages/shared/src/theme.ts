export type ThemeName = "classic" | "midnight" | "atelier" | "zeus";

export type ThemeTokens = {
  background: string;
  backgroundAccent: string;
  backgroundAmbient: string;
  workspaceBackground: string;
  sidebar: string;
  sidebarText: string;
  surface: string;
  surfaceElevated: string;
  border: string;
  borderStrong: string;
  text: string;
  textMuted: string;
  accent: string;
  accentStrong: string;
  accentSoft: string;
  shadow: string;
  inputBackground: string;
  inputText: string;
  buttonPrimaryBackground: string;
  buttonPrimaryText: string;
  buttonSecondaryBackground: string;
  buttonSecondaryText: string;
  overlay: string;
  success: string;
  successSoft: string;
  danger: string;
  dangerSoft: string;
};

export type ThemeDefinition = {
  name: ThemeName;
  label: string;
  description: string;
  tokens: ThemeTokens;
};

export const themeStorageKey = "gym-saas-theme";
export const defaultThemeName: ThemeName = "classic";

export const themes: Record<ThemeName, ThemeDefinition> = {
  classic: {
    name: "classic",
    label: "Classic Ledger",
    description: "Warm paper surfaces and brass accents.",
    tokens: {
      background: "#efe3c6",
      backgroundAccent: "#f8f4ee",
      backgroundAmbient: "#e9eef0",
      workspaceBackground: "#f8f4ee",
      sidebar: "#2d1b05",
      sidebarText: "#fff8ef",
      surface: "#fff8ef",
      surfaceElevated: "#f5eadc",
      border: "#d8c4a7",
      borderStrong: "#b99156",
      text: "#2d1b05",
      textMuted: "#6f5634",
      accent: "#8a5a18",
      accentStrong: "#6b3f0a",
      accentSoft: "#f2d4a9",
      shadow: "0 24px 60px rgba(47, 27, 8, 0.24)",
      inputBackground: "#fffdf9",
      inputText: "#2d1b05",
      buttonPrimaryBackground: "#8a5a18",
      buttonPrimaryText: "#fff9f1",
      buttonSecondaryBackground: "#efe2d0",
      buttonSecondaryText: "#4a3110",
      overlay: "rgba(38, 23, 8, 0.45)",
      success: "#297a43",
      successSoft: "#e9f6eb",
      danger: "#b42318",
      dangerSoft: "#fceae7"
    }
  },
  midnight: {
    name: "midnight",
    label: "Midnight Terminal",
    description: "Deep navy panels and signal-blue accents.",
    tokens: {
      background: "#08111f",
      backgroundAccent: "#13233d",
      backgroundAmbient: "#1f3144",
      workspaceBackground: "#0b1424",
      sidebar: "#08111f",
      sidebarText: "#eef7ff",
      surface: "#0f1c31",
      surfaceElevated: "#162744",
      border: "#26446f",
      borderStrong: "#44d1ff",
      text: "#eef7ff",
      textMuted: "#9bb8d4",
      accent: "#44d1ff",
      accentStrong: "#6ff4ff",
      accentSoft: "#123f58",
      shadow: "0 24px 60px rgba(2, 10, 21, 0.48)",
      inputBackground: "#0b1728",
      inputText: "#eef7ff",
      buttonPrimaryBackground: "#44d1ff",
      buttonPrimaryText: "#071421",
      buttonSecondaryBackground: "#16314d",
      buttonSecondaryText: "#d4edff",
      overlay: "rgba(3, 10, 20, 0.62)",
      success: "#67df91",
      successSoft: "#123827",
      danger: "#ff9994",
      dangerSoft: "#482330"
    }
  },
  atelier: {
    name: "atelier",
    label: "Atelier Bloom",
    description: "Soft stone surfaces with botanical contrast.",
    tokens: {
      background: "#f4efe8",
      backgroundAccent: "#f6d9cf",
      backgroundAmbient: "#dbe9da",
      workspaceBackground: "#f7f5f0",
      sidebar: "#2d2620",
      sidebarText: "#fffdf9",
      surface: "#fffdf9",
      surfaceElevated: "#f2f0ea",
      border: "#c9b7a5",
      borderStrong: "#2d6d54",
      text: "#2d2620",
      textMuted: "#6f645b",
      accent: "#2d6d54",
      accentStrong: "#1d4937",
      accentSoft: "#d5e8dd",
      shadow: "0 24px 60px rgba(66, 49, 34, 0.16)",
      inputBackground: "#ffffff",
      inputText: "#2d2620",
      buttonPrimaryBackground: "#2d6d54",
      buttonPrimaryText: "#f7fffb",
      buttonSecondaryBackground: "#ebe1d4",
      buttonSecondaryText: "#3f3429",
      overlay: "rgba(41, 34, 28, 0.35)",
      success: "#2d6d54",
      successSoft: "#e4f2e8",
      danger: "#bd3e36",
      dangerSoft: "#fbe9e7"
    }
  },
  zeus: {
    name: "zeus",
    label: "Zeus Gold",
    description: "Charcoal surfaces with the gold and amber signal colors of Zeus Gym.",
    tokens: {
      background: "#1c1510",
      backgroundAccent: "#b75b1e",
      backgroundAmbient: "#ffd51b",
      workspaceBackground: "#fafafb",
      sidebar: "#1b1d1c",
      sidebarText: "#eef0ef",
      surface: "#25221e",
      surfaceElevated: "#2a2621",
      border: "#66511d",
      borderStrong: "#f6ca19",
      text: "#f5f2ee",
      textMuted: "#a5a8b0",
      accent: "#f6ca19",
      accentStrong: "#f28a32",
      accentSoft: "#4f3e17",
      shadow: "0 28px 72px rgba(28, 14, 4, 0.42)",
      inputBackground: "#26231f",
      inputText: "#f5f2ee",
      buttonPrimaryBackground: "linear-gradient(100deg, #ffd919 0%, #ffc51d 42%, #ef8139 100%)",
      buttonPrimaryText: "#201c18",
      buttonSecondaryBackground: "#373126",
      buttonSecondaryText: "#f5f2ee",
      overlay: "rgba(21, 14, 8, 0.62)",
      success: "#4ca555",
      successSoft: "#edf8ee",
      danger: "#cf443a",
      dangerSoft: "#fff3f1"
    }
  }
};

export function getThemeDefinition(name: ThemeName | string | undefined): ThemeDefinition {
  if (!name) {
    return themes[defaultThemeName];
  }

  return themes[name as ThemeName] ?? themes[defaultThemeName];
}

export function listThemes(): ThemeDefinition[] {
  return Object.values(themes);
}
