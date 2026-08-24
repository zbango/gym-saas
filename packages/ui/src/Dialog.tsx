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
    <div className="fixed inset-0 z-50 grid place-items-center overflow-y-auto bg-[rgba(21,14,8,0.62)] p-3 max-[560px]:p-[14px]" onMouseDown={onClose}>
      <section
        className={`max-h-[calc(100dvh-24px)] w-full overflow-y-auto rounded-[26px] bg-white text-[#172131] shadow-[0_28px_80px_rgba(12,18,28,0.34)] max-[560px]:max-h-[calc(100dvh-28px)] max-[560px]:rounded-[18px] ${className ?? ""}`.trim()}
        role="dialog"
        aria-modal="true"
        aria-labelledby="gs-dialog-title"
        onMouseDown={(event) => event.stopPropagation()}
      >
        <header className="flex min-h-[122px] items-center justify-between border-b border-[#e1e4e8] px-11 max-[560px]:min-h-[90px] max-[560px]:px-[23px]">
          <h2 id="gs-dialog-title" className="m-0 text-[38px] font-extrabold tracking-[-0.04em] max-[560px]:text-[28px]">{title}</h2>
          <button className="grid h-[46px] w-[46px] place-items-center border-0 bg-transparent p-0 text-[#9da5b0] [&>svg]:h-9 [&>svg]:w-9 max-[560px]:[&>svg]:h-[29px] max-[560px]:[&>svg]:w-[29px]" type="button" onClick={onClose} aria-label="Cerrar diálogo"><XIcon /></button>
        </header>
        {children}
      </section>
    </div>
  );
}
