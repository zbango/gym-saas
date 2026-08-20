type DesktopApp = Record<string, unknown>;

declare global {
  interface Window {
    go?: {
      main?: {
        App?: DesktopApp;
      };
    };
  }
}

// The Wails namespace is the only global bridge. Feature modules narrow it to
// the small set of methods they own rather than collecting every binding here.
export function desktopApp(): DesktopApp | null {
  return window.go?.main?.App ?? null;
}

export function isDesktopApp(): boolean {
  return desktopApp() !== null;
}
