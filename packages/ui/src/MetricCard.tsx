import type { ReactNode } from "react";

export type MetricTone = "neutral" | "success" | "warning" | "danger";

export type MetricCardProps = {
  value: string | number;
  label: string;
  icon: ReactNode;
  tone?: MetricTone;
  selected?: boolean;
  className?: string;
};

export function MetricCard({ value, label, icon, tone = "neutral", selected = false, className }: MetricCardProps) {
  return (
    <article className={`gs-metric-card gs-metric-card--${tone} ${selected ? "is-selected" : ""} ${className ?? ""}`.trim()}>
      <span className="gs-metric-card-icon" aria-hidden="true">{icon}</span>
      <div>
        <strong>{value}</strong>
        <p>{label}</p>
      </div>
    </article>
  );
}
