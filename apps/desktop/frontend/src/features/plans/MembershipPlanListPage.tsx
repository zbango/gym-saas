import { useState, type ReactNode } from "react";
import {
  DocumentIcon,
  EditIcon,
  PageHeader,
  PlusIcon,
  RefreshIcon,
  TrashIcon
} from "@gym-saas/ui";
import type { MembershipPlan as Plan } from "./api";
import { useMembershipPlanList } from "./useMembershipPlanList";
import { MembershipPlanCreateDialog } from "./MembershipPlanCreateDialog";

const previewPlans: Plan[] = [
  { id: "preview-daily-gym", name: "Plan Diario GYM", validityKind: "time", durationValue: 1, durationUnit: "days", visitLimit: 0, priceCents: 200, currency: "USD", status: "active" },
  { id: "preview-daily-mma", name: "Plan Diario MMA", validityKind: "time", durationValue: 1, durationUnit: "days", visitLimit: 0, priceCents: 300, currency: "USD", status: "active" },
  { id: "preview-biweekly", name: "Plan Quincenal", validityKind: "time", durationValue: 2, durationUnit: "weeks", visitLimit: 0, priceCents: 1300, currency: "USD", status: "active" },
  { id: "preview-monthly", name: "Plan Mensual", validityKind: "time", durationValue: 1, durationUnit: "months", visitLimit: 0, priceCents: 2500, currency: "USD", status: "active" },
  { id: "preview-monthly-mma", name: "Plan Mensual MMA", validityKind: "time", durationValue: 1, durationUnit: "months", visitLimit: 0, priceCents: 2500, currency: "USD", status: "active" },
  { id: "preview-visits", name: "Plan Tarjeta", validityKind: "visits", durationValue: 1, durationUnit: "months", visitLimit: 25, priceCents: 3000, currency: "USD", status: "active" }
];

export function MembershipPlanListPage() {
  const { plans, loading, error, refresh, desktopRuntime } = useMembershipPlanList();
  const displayedPlans = desktopRuntime ? plans : previewPlans;
  const [createOpen, setCreateOpen] = useState(false);

  return (
    <div className="px-16 py-14 text-[#172131] max-[1260px]:px-[43px] max-[1260px]:py-[46px] max-[920px]:px-[22px] max-[920px]:py-8 max-[560px]:px-[14px] max-[560px]:py-6">
      <PageHeader icon={<DocumentIcon />} iconClassName="h-[66px] w-[66px] rounded-[13px] bg-[#242e3e] pt-0 text-white" title="Planes de Membresía" description="Crea y gestiona planes de membresía, establece precios y configura beneficios para miembros." actions={<button className="inline-flex min-h-[58px] items-center gap-3 rounded-xl border-0 bg-[linear-gradient(100deg,var(--color-brand-button-start)_0%,var(--color-brand-button-middle)_42%,var(--color-brand-button-end)_100%)] px-[29px] text-xl font-extrabold whitespace-nowrap text-[var(--color-brand-button-text)] shadow-[0_12px_22px_rgba(178,103,33,0.18)] max-[1260px]:min-h-[49px] max-[1260px]:px-5 max-[1260px]:text-[17px] max-[920px]:w-fit max-[560px]:w-full max-[560px]:justify-center" type="button" onClick={() => setCreateOpen(true)}><PlusIcon className="h-[25px] w-[25px]" /> Agregar Nuevo Plan</button>} />

      <MembershipPlanCreateDialog open={createOpen} onClose={() => setCreateOpen(false)} onCreated={refresh} />

      {loading ? <PlanCardSkeletons /> : null}
      {!loading && error ? <PlanMessage tone="error" message={error} action={<button className="inline-flex min-h-[45px] items-center gap-2 rounded-[10px] border-0 bg-[var(--color-brand-danger)] px-[18px] font-extrabold text-white [&>svg]:h-5 [&>svg]:w-5" type="button" onClick={() => void refresh()}><RefreshIcon /> Reintentar</button>} /> : null}
      {!loading && !error && displayedPlans.length === 0 ? <PlanMessage message="No hay planes de membresía activos todavía." /> : null}
      {!loading && !error && displayedPlans.length > 0 ? (
        <section className="grid grid-cols-3 gap-10 max-[1260px]:grid-cols-2 max-[1260px]:gap-[27px] max-[920px]:grid-cols-1" aria-label="Planes de membresía">
          {displayedPlans.map((plan) => <MembershipPlanCard key={plan.id} plan={plan} />)}
        </section>
      ) : null}
    </div>
  );
}

function MembershipPlanCard({ plan }: { plan: Plan }) {
  const active = plan.status === "active";
  return (
    <article className="flex min-h-[425px] flex-col rounded-[28px] border border-[#e9ebef] bg-white px-[38px] pb-[34px] pt-[38px] shadow-[0_3px_5px_rgba(20,31,48,0.04)] max-[1260px]:min-h-[385px] max-[1260px]:p-[30px] max-[920px]:min-h-[360px] max-[560px]:p-6">
      <header className="flex items-center justify-between gap-4">
        <span className="rounded-full bg-[#24303f] px-[21px] py-[10px] text-[17px] font-extrabold text-white">
          {planKindLabel(plan)}
        </span>
        <span className={`rounded-full px-[21px] py-[9px] text-[17px] font-extrabold ${active ? "bg-[var(--color-brand-success-soft)] text-[#427a4a]" : "bg-[#edf0f3] text-[#65707b]"}`}>
          {active ? "Activo" : "Inactivo"}
        </span>
      </header>
      <h2 className="mt-[39px] text-[33px] font-extrabold tracking-[-0.04em] max-[560px]:text-[28px]">{plan.name}</h2>
      <strong className="mt-[25px] block text-[57px] leading-none font-extrabold tracking-[-0.04em] max-[560px]:text-5xl">{formatMoney(plan.priceCents, plan.currency)}</strong>
      <p className="mt-5 text-[21px] text-[#59636e]">{durationLabel(plan)}</p>
      <div className="mt-auto grid grid-cols-2 gap-[15px] pt-[37px]">
        <button className="inline-flex min-h-[70px] items-center justify-center gap-3 rounded-xl border-0 bg-[linear-gradient(100deg,var(--color-brand-button-start)_0%,var(--color-brand-button-middle)_42%,var(--color-brand-button-end)_100%)] text-[22px] font-extrabold text-[var(--color-brand-button-text)] max-[560px]:min-h-[57px] max-[560px]:text-lg" type="button"><EditIcon className="h-[25px] w-[25px]" /> Editar</button>
        <button className="inline-flex min-h-[70px] items-center justify-center gap-3 rounded-xl border-0 bg-[#8d4d1e] text-[22px] font-extrabold text-white max-[560px]:min-h-[57px] max-[560px]:text-lg" type="button"><TrashIcon className="h-[25px] w-[25px]" /> Eliminar</button>
      </div>
    </article>
  );
}

function PlanCardSkeletons() {
  return (
    <section className="grid grid-cols-3 gap-10 max-[1260px]:grid-cols-2 max-[1260px]:gap-[27px] max-[920px]:grid-cols-1" aria-label="Cargando planes de membresía">
      {Array.from({ length: 6 }, (_, index) => (
        <div
          className="min-h-[425px] animate-pulse rounded-[28px] bg-[linear-gradient(100deg,#f1f2f4_22%,#fafafa_38%,#f1f2f4_55%)] bg-[length:200%_100%] max-[1260px]:min-h-[385px] max-[920px]:min-h-[360px]"
          key={index}
          aria-hidden="true"
        />
      ))}
    </section>
  );
}

function PlanMessage(props: { message: string; tone?: "error"; action?: ReactNode }) {
  return (
    <section
      className={`grid min-h-[210px] place-items-center gap-[18px] rounded-3xl border bg-white p-8 text-center text-xl ${
        props.tone === "error"
          ? "border-[color-mix(in_srgb,var(--color-brand-danger)_35%,transparent)] text-[var(--color-brand-danger)]"
          : "border-[#e7e9ed] text-[#64707b]"
      }`}
    >
      <p className="m-0">{props.message}</p>
      {props.action}
    </section>
  );
}

function planKindLabel(plan: Plan): string {
  if (plan.validityKind === "visits") {
    return "Plan por visitas";
  }
  if (plan.durationUnit === "days" && plan.durationValue === 1) {
    return "Plan Diario";
  }
  if (plan.durationUnit === "weeks" && plan.durationValue === 2) {
    return "Plan Quincenal";
  }
  if (plan.durationUnit === "months" && plan.durationValue === 1) {
    return "Plan Mensual";
  }
  return "Plan por período";
}

function durationLabel(plan: Plan): string {
  if (plan.validityKind === "visits") {
    return `por ${plan.visitLimit} ${plan.visitLimit === 1 ? "visita" : "visitas"}`;
  }
  const unit = {
    days: plan.durationValue === 1 ? "día" : "días",
    weeks: plan.durationValue === 1 ? "semana" : "semanas",
    months: plan.durationValue === 1 ? "mes" : "meses",
    years: plan.durationValue === 1 ? "año" : "años"
  }[plan.durationUnit];
  return `por ${plan.durationValue} ${unit}`;
}

function formatMoney(priceCents: number, currency: string): string {
  return new Intl.NumberFormat("en-US", { style: "currency", currency }).format(priceCents / 100);
}
