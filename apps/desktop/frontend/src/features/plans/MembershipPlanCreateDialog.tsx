import { useEffect, useState, type FormEvent } from "react";
import { Dialog } from "@gym-saas/ui";
import { isDesktopApp } from "../../platform/wails";
import { createMembershipPlan, type MembershipPlanInput } from "./api";

type PlanType = "daily" | "weekly" | "biweekly" | "monthly" | "quarterly" | "semiannual" | "annual" | "visits";

type PlanTemplate = {
  label: string;
  validityKind: "time" | "visits";
  durationValue: number;
  durationUnit: MembershipPlanInput["durationUnit"];
};

const planTemplates: Record<PlanType, PlanTemplate> = {
  daily: { label: "Plan Diario", validityKind: "time", durationValue: 1, durationUnit: "days" },
  weekly: { label: "Plan Semanal", validityKind: "time", durationValue: 1, durationUnit: "weeks" },
  biweekly: { label: "Plan Quincenal", validityKind: "time", durationValue: 2, durationUnit: "weeks" },
  monthly: { label: "Plan Mensual", validityKind: "time", durationValue: 1, durationUnit: "months" },
  quarterly: { label: "Plan Trimestral", validityKind: "time", durationValue: 3, durationUnit: "months" },
  semiannual: { label: "Plan Semestral", validityKind: "time", durationValue: 6, durationUnit: "months" },
  annual: { label: "Plan Anual", validityKind: "time", durationValue: 1, durationUnit: "years" },
  visits: { label: "Plan por Visitas", validityKind: "visits", durationValue: 1, durationUnit: "months" }
};

type FormState = {
  name: string;
  type: PlanType;
  price: string;
  status: MembershipPlanInput["status"];
  visitLimit: string;
  visitDurationValue: string;
  visitDurationUnit: MembershipPlanInput["durationUnit"];
};

const initialForm: FormState = {
  name: "",
  type: "daily",
  price: "0",
  status: "active",
  visitLimit: "25",
  visitDurationValue: "1",
  visitDurationUnit: "months"
};

export function MembershipPlanCreateDialog(props: {
  open: boolean;
  onClose: () => void;
  onCreated: () => Promise<void>;
}) {
  const [form, setForm] = useState<FormState>(initialForm);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const desktopRuntime = isDesktopApp();
  const selectedTemplate = planTemplates[form.type];
  const visitPlan = selectedTemplate.validityKind === "visits";

  useEffect(() => {
    if (props.open) {
      setForm(initialForm);
      setError(null);
    }
  }, [props.open]);

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!desktopRuntime) {
      setError("La creación de planes requiere la aplicación de escritorio.");
      return;
    }

    const price = Number(form.price);
    const durationValue = visitPlan ? Number(form.visitDurationValue) : selectedTemplate.durationValue;
    const durationUnit = visitPlan ? form.visitDurationUnit : selectedTemplate.durationUnit;
    const visitLimit = visitPlan ? Number(form.visitLimit) : 0;

    if (!Number.isFinite(price) || price < 0 || !Number.isInteger(durationValue) || durationValue < 1 || (visitPlan && (!Number.isInteger(visitLimit) || visitLimit < 1))) {
      setError("Revisa el precio, duración y visitas antes de crear el plan.");
      return;
    }

    setSubmitting(true);
    setError(null);
    try {
      await createMembershipPlan({
        name: form.name,
        validityKind: selectedTemplate.validityKind,
        durationValue,
        durationUnit,
        visitLimit,
        priceCents: Math.round(price * 100),
        currency: "USD",
        status: form.status
      });
      await props.onCreated();
      props.onClose();
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : "No se pudo crear el plan.");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <Dialog open={props.open} onClose={props.onClose} title="Agregar Nuevo Plan" className="membership-plan-create-dialog">
      <form onSubmit={(event) => void submit(event)}>
        <label htmlFor="plan-name">
          <span>Nombre del Plan</span>
          <input id="plan-name" required value={form.name} onChange={(event) => setForm((current) => ({ ...current, name: event.target.value }))} placeholder="Ej: Plan Mensual Básico" />
        </label>

        <label htmlFor="plan-type">
          <span>Tipo de Plan</span>
          <select id="plan-type" value={form.type} onChange={(event) => setForm((current) => ({ ...current, type: event.target.value as PlanType }))}>
            {(Object.keys(planTemplates) as PlanType[]).map((type) => <option key={type} value={type}>{planTemplates[type].label}</option>)}
          </select>
        </label>

        <label htmlFor="plan-price">
          <span>Precio (USD)</span>
          <input id="plan-price" required type="number" min="0" step="0.01" inputMode="decimal" value={form.price} onChange={(event) => setForm((current) => ({ ...current, price: event.target.value }))} />
        </label>

        {visitPlan ? (
          <div className="visit-plan-fields">
            <label htmlFor="plan-visit-limit">
              <span>Cantidad de Visitas</span>
              <input id="plan-visit-limit" required type="number" min="1" value={form.visitLimit} onChange={(event) => setForm((current) => ({ ...current, visitLimit: event.target.value }))} />
            </label>
            <label htmlFor="plan-visit-duration-value">
              <span>Vigencia</span>
              <div className="duration-controls">
                <input id="plan-visit-duration-value" required type="number" min="1" value={form.visitDurationValue} onChange={(event) => setForm((current) => ({ ...current, visitDurationValue: event.target.value }))} />
                <select aria-label="Unidad de vigencia" value={form.visitDurationUnit} onChange={(event) => setForm((current) => ({ ...current, visitDurationUnit: event.target.value as MembershipPlanInput["durationUnit"] }))}>
                  <option value="days">días</option><option value="weeks">semanas</option><option value="months">meses</option><option value="years">años</option>
                </select>
              </div>
            </label>
          </div>
        ) : (
          <label htmlFor="plan-duration">
            <span>Duración</span>
            <input id="plan-duration" disabled value={fixedDurationLabel(selectedTemplate)} />
          </label>
        )}

        <label htmlFor="plan-status">
          <span>Estado</span>
          <select id="plan-status" value={form.status} onChange={(event) => setForm((current) => ({ ...current, status: event.target.value as MembershipPlanInput["status"] }))}>
            <option value="active">Activo</option><option value="inactive">Inactivo</option>
          </select>
        </label>

        {error ? <p className="create-plan-error" role="alert">{error}</p> : null}
        <footer>
          <button type="button" className="create-plan-cancel" onClick={props.onClose} disabled={submitting}>Cancelar</button>
          <button type="submit" className="create-plan-submit" disabled={submitting}>{submitting ? "Creando..." : "Crear Plan"}</button>
        </footer>
      </form>
    </Dialog>
  );
}

function fixedDurationLabel(template: PlanTemplate): string {
  const unit = {
    days: template.durationValue === 1 ? "día" : "días",
    weeks: template.durationValue === 1 ? "semana" : "semanas",
    months: template.durationValue === 1 ? "mes" : "meses",
    years: template.durationValue === 1 ? "año" : "años"
  }[template.durationUnit];
  return `${template.durationValue} ${unit}`;
}
