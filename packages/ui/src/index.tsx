import type { CSSProperties, PropsWithChildren } from "react";

export * from "./theme";
export * from "./icons";
export * from "./DataTable";
export * from "./Dialog";
export * from "./MetricCard";
export * from "./PageHeader";

export function Panel(props: PropsWithChildren<{ title: string; eyebrow?: string }>) {
  return (
    <section
      style={{
        border: "1px solid var(--gs-border)",
        borderRadius: 16,
        padding: 20,
        background: "linear-gradient(180deg, var(--gs-surface) 0%, var(--gs-surface-elevated) 100%)",
        boxShadow: "var(--gs-shadow)",
        color: "var(--gs-text)"
      }}
    >
      {props.eyebrow ? (
        <div
          style={{
            fontSize: 12,
            fontWeight: 700,
            letterSpacing: "0.14em",
            textTransform: "uppercase",
            color: "var(--gs-accent)",
            marginBottom: 10
          }}
        >
          {props.eyebrow}
        </div>
      ) : null}
      <h2 style={{ margin: "0 0 12px", color: "var(--gs-text)" }}>{props.title}</h2>
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
  const style: CSSProperties =
    variant === "primary"
      ? {
          background: "var(--gs-button-primary-background)",
          color: "var(--gs-button-primary-text)"
        }
      : {
          background: "var(--gs-button-secondary-background)",
          color: "var(--gs-button-secondary-text)"
        };

  return (
    <button
      onClick={props.onClick}
      disabled={props.disabled}
      style={{
        ...style,
        border: "1px solid var(--gs-border)",
        borderRadius: 12,
        padding: "10px 14px",
        font: "inherit",
        fontWeight: 600,
        cursor: props.disabled ? "not-allowed" : "pointer",
        opacity: props.disabled ? 0.6 : 1
      }}
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
      style={{
        padding: "10px 12px",
        flex: 1,
        background: "var(--gs-input-background)",
        color: "var(--gs-input-text)",
        border: "1px solid var(--gs-border)",
        borderRadius: 12,
        font: "inherit"
      }}
    />
  );
}
