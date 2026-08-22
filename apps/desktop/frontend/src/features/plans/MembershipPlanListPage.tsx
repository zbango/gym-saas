import { type ReactNode } from "react";
import {
  DocumentIcon,
  EditIcon,
  PlusIcon,
  TrashIcon,
  PageHeader,
  RefreshIcon
} from "@gym-saas/ui";
import type { MembershipPlan as Plan } from "./api";
import { useMembershipPlanList } from "./useMembershipPlanList";
import "./membership-plan-list-page.css";

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

  return (
    <div className="membership-plan-list-page">
      <PageHeader
        className="membership-plan-page-header"
        icon={<DocumentIcon />}
        title="Planes de Membresía"
        description="Crea y gestiona planes de membresía, establece precios y configura beneficios para miembros."
        actions={<button className="membership-plan-create-action" type="button"><PlusIcon /> Agregar Nuevo Plan</button>}
      />

      {loading ? <PlanCardSkeletons /> : null}
      {!loading && error ? <PlanMessage tone="error" message={error} action={<button type="button" onClick={() => void refresh()}><RefreshIcon /> Reintentar</button>} /> : null}
      {!loading && !error && displayedPlans.length === 0 ? <PlanMessage message="No hay planes de membresía activos todavía." /> : null}
      {!loading && !error && displayedPlans.length > 0 ? (
        <section className="membership-plan-grid" aria-label="Planes de membresía">
          {displayedPlans.map((plan) => <MembershipPlanCard key={plan.id} plan={plan} />)}
        </section>
      ) : null}
    </div>
  );
}

function MembershipPlanCard({ plan }: { plan: Plan }) {
  const active = plan.status === "active";
  return (
    <article className="membership-plan-card">
      <header>
        <span className="plan-kind">{planKindLabel(plan)}</span>
        <span className={`plan-status ${active ? "is-active" : "is-inactive"}`}>{active ? "Activo" : "Inactivo"}</span>
      </header>
      <h2>{plan.name}</h2>
      <strong className="plan-price">{formatMoney(plan.priceCents, plan.currency)}</strong>
      <p className="plan-duration">{durationLabel(plan)}</p>
      <div className="plan-card-actions">
        <button className="plan-edit" type="button"><EditIcon /> Editar</button>
        <button className="plan-delete" type="button"><TrashIcon /> Eliminar</button>
      </div>
    </article>
  );
}

function PlanCardSkeletons() {
  return (
    <section className="membership-plan-grid" aria-label="Cargando planes de membresía">
      {Array.from({ length: 6 }, (_, index) => <div className="membership-plan-skeleton" key={index} aria-hidden="true" />)}
    </section>
  );
}

function PlanMessage(props: { message: string; tone?: "error"; action?: ReactNode }) {
  return <section className={`membership-plan-message ${props.tone === "error" ? "is-error" : ""}`}><p>{props.message}</p>{props.action}</section>;
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
