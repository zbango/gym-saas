import type { PropsWithChildren } from "react";

export function Panel(props: PropsWithChildren<{ title: string; eyebrow?: string }>) {
  return (
    <section
      style={{
        border: "1px solid #d9d2c5",
        borderRadius: 16,
        padding: 20,
        background: "linear-gradient(180deg, #fff8f0 0%, #f7efe5 100%)",
        boxShadow: "0 10px 30px rgba(64, 42, 17, 0.08)"
      }}
    >
      {props.eyebrow ? (
        <div
          style={{
            fontSize: 12,
            fontWeight: 700,
            letterSpacing: "0.14em",
            textTransform: "uppercase",
            color: "#8a5a18",
            marginBottom: 10
          }}
        >
          {props.eyebrow}
        </div>
      ) : null}
      <h2 style={{ margin: "0 0 12px", color: "#2d1b05" }}>{props.title}</h2>
      <div>{props.children}</div>
    </section>
  );
}

