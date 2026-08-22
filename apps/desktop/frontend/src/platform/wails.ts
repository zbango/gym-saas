// Feature clients import Wails-generated modules directly. This tiny helper
// only lets feature hooks distinguish the desktop runtime from a browser.
export function isDesktopApp(): boolean {
  return typeof window !== "undefined" && "go" in window;
}
