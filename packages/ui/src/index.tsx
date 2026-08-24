import type { PropsWithChildren } from "react";

export * from "./icons";
export * from "./DataTable";
export * from "./Dialog";
export * from "./MetricCard";
export * from "./PageHeader";
export * from "./Tabs";

export function Panel(props: PropsWithChildren<{ title: string; eyebrow?: string }>) {
  return (
    <section
      className="rounded-2xl border border-[var(--color-brand-border)] bg-[linear-gradient(180deg,var(--color-brand-surface)_0%,var(--color-brand-surface-elevated)_100%)] p-5 text-[var(--color-brand-text)] shadow-[0_28px_72px_rgba(28,14,4,0.42)]"
    >
      {props.eyebrow ? (
        <div className="mb-2.5 text-xs font-bold uppercase tracking-[0.14em] text-[var(--color-brand-accent)]">
          {props.eyebrow}
        </div>
      ) : null}
      <h2 className="mb-3 text-[var(--color-brand-text)]">{props.title}</h2>
      <div>{props.children}</div>
    </section>
  );
}

export function ShellButton(
  props: PropsWithChildren<{
    onClick?: () => void | Promise<void>;
    disabled?: boolean;
    variant?: "primary" | "secondary";
  }>
) {
  const variant = props.variant ?? "primary";
  const toneClasses =
    variant === "primary"
      ? "bg-[linear-gradient(100deg,var(--color-brand-button-start)_0%,var(--color-brand-button-middle)_42%,var(--color-brand-button-end)_100%)] text-[var(--color-brand-button-text)]"
      : "bg-[#373126] text-[var(--color-brand-text)]";
  const stateClasses = props.disabled ? "cursor-not-allowed opacity-60" : "cursor-pointer";

  return (
    <button
      onClick={props.onClick}
      disabled={props.disabled}
      className={`rounded-xl border border-[var(--color-brand-border)] px-3.5 py-2.5 font-semibold ${toneClasses} ${stateClasses}`}
    >
      {props.children}
    </button>
  );
}

export function ShellInput(props: {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
}) {
  return (
    <input
      value={props.value}
      onChange={(event) => props.onChange(event.target.value)}
      placeholder={props.placeholder}
      className="flex-1 rounded-xl border border-[var(--color-brand-border)] bg-[var(--color-brand-input)] px-3 py-2.5 text-[var(--color-brand-text)]"
    />
  );
}
