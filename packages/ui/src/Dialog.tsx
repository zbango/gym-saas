import { useEffect, type ReactNode } from "react";
import { XIcon } from "./icons";

export type DialogProps = {
  open: boolean;
  title: string;
  onClose: () => void;
  children: ReactNode;
  className?: string;
};

export function Dialog({ open, title, onClose, children, className }: DialogProps) {
  useEffect(() => {
    if (!open) {
      return;
    }

    function closeOnEscape(event: KeyboardEvent) {
      if (event.key === "Escape") {
        onClose();
      }
    }

    window.addEventListener("keydown", closeOnEscape);
    return () => window.removeEventListener("keydown", closeOnEscape);
  }, [open, onClose]);

  if (!open) {
    return null;
  }

  return (
    <div className="gs-dialog-backdrop" onMouseDown={onClose}>
      <section
        className={`gs-dialog ${className ?? ""}`.trim()}
        role="dialog"
        aria-modal="true"
        aria-labelledby="gs-dialog-title"
        onMouseDown={(event) => event.stopPropagation()}
      >
        <header className="gs-dialog-header">
          <h2 id="gs-dialog-title">{title}</h2>
          <button type="button" onClick={onClose} aria-label="Cerrar diálogo"><XIcon /></button>
        </header>
        {children}
      </section>
    </div>
  );
}
