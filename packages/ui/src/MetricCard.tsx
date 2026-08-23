import type { ReactNode } from "react";

export type MetricTone = "neutral" | "success" | "warning" | "danger" | "info";

export type MetricCardProps = {
  value: string | number;
  label: string;
  icon: ReactNode;
  tone?: MetricTone;
  selected?: boolean;
  className?: string;
};

export function MetricCard({ value, label, icon, tone = "neutral", selected = false, className }: MetricCardProps) {
  const toneClasses = {
    neutral: "bg-[#f7f1ff] text-[#954fee]",
    success: "bg-[var(--color-brand-success-soft)] text-[var(--color-brand-success)]",
    warning: "bg-[#fffbea] text-[#c59327]",
    danger: "bg-[var(--color-brand-danger-soft)] text-[var(--color-brand-danger)]",
    info: "bg-[#edf4ff] text-[#4b78dd]"
  }[tone];
  return (
    <article className={`flex min-h-[185px] items-center gap-6 rounded-[15px] border border-[#e0e3e8] bg-white px-[28px] py-[31px] ${selected ? "border-3 border-[#a46df2] shadow-[0_0_0_4px_rgba(164,109,242,0.12)]" : ""} ${className ?? ""}`.trim()}>
      <span className={`grid h-[76px] w-[76px] shrink-0 place-items-center rounded-[16px] ${toneClasses} [&>svg]:h-[39px] [&>svg]:w-[39px]`} aria-hidden="true">{icon}</span>
      <div>
        <strong className="block text-[40px] leading-none">{value}</strong>
        <p className="mt-[10px] text-[20px] leading-[1.45] text-[#59636e]">{label}</p>
      </div>
    </article>
  );
}
