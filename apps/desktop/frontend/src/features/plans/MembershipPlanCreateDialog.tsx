import { useEffect, useState, type FormEvent, type ReactNode, type SelectHTMLAttributes } from "react";
import { ChevronDownIcon, XIcon } from "@gym-saas/ui";
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
  const { open, onClose, onCreated } = props;
  const [form, setForm] = useState<FormState>(initialForm);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const desktopRuntime = isDesktopApp();
  const selectedTemplate = planTemplates[form.type];
  const visitPlan = selectedTemplate.validityKind === "visits";

  useEffect(() => {
    if (open) {
      setForm(initialForm);
      setError(null);
    }
  }, [open]);

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
      await onCreated();
      onClose();
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : "No se pudo crear el plan.");
    } finally {
      setSubmitting(false);
    }
  }

  if (!open) {
    return null;
  }

  return (
    <div
      className="fixed inset-0 z-20 grid place-items-center overflow-y-auto bg-[rgba(21,14,8,0.62)] p-7 max-[560px]:p-[14px]"
      onMouseDown={onClose}
    >
      <section
        className="max-h-[calc(100dvh-56px)] w-full max-w-[760px] overflow-y-auto rounded-[26px] bg-white text-[#172131] shadow-[0_28px_80px_rgba(12,18,28,0.34)] max-[560px]:max-h-[calc(100dvh-28px)] max-[560px]:rounded-[18px]"
        role="dialog"
        aria-modal="true"
        aria-labelledby="membership-plan-create-title"
        onMouseDown={(event) => event.stopPropagation()}
      >
        <header className="flex min-h-[122px] items-center justify-between border-b border-[#e1e4e8] px-11 max-[560px]:min-h-[90px] max-[560px]:px-[23px]">
          <h2 id="membership-plan-create-title" className="m-0 text-[38px] font-extrabold tracking-[-0.04em] max-[560px]:text-[28px]">
            Agregar Nuevo Plan
          </h2>
          <button
            className="grid h-[46px] w-[46px] place-items-center border-0 bg-transparent p-0 text-[#9da5b0] [&>svg]:h-9 [&>svg]:w-9 max-[560px]:[&>svg]:h-[29px] max-[560px]:[&>svg]:w-[29px]"
            type="button"
            onClick={onClose}
            aria-label="Cerrar diálogo"
          >
            <XIcon />
          </button>
        </header>

        <form className="grid gap-[27px] px-11 pb-10 pt-12 max-[560px]:gap-5 max-[560px]:px-[23px] max-[560px]:py-7" onSubmit={(event) => void submit(event)}>
          <FormField label="Nombre del Plan" htmlFor="plan-name">
            <input
              id="plan-name"
              required
              value={form.name}
              onChange={(event) => setForm((current) => ({ ...current, name: event.target.value }))}
              placeholder="Ej: Plan Mensual Básico"
              className={fieldClassName}
            />
          </FormField>

          <FormField label="Tipo de Plan" htmlFor="plan-type">
            <SelectField
              id="plan-type"
              value={form.type}
              onChange={(event) => setForm((current) => ({ ...current, type: event.target.value as PlanType }))}
            >
              {(Object.keys(planTemplates) as PlanType[]).map((type) => (
                <option key={type} value={type}>{planTemplates[type].label}</option>
              ))}
            </SelectField>
          </FormField>

          <FormField label="Precio (USD)" htmlFor="plan-price">
            <input
              id="plan-price"
              required
              type="number"
              min="0"
              step="0.01"
              inputMode="decimal"
              value={form.price}
              onChange={(event) => setForm((current) => ({ ...current, price: event.target.value }))}
              className={fieldClassName}
            />
          </FormField>

          {visitPlan ? (
            <div className="grid grid-cols-[1fr_1.35fr] gap-[21px] max-[560px]:grid-cols-1">
              <FormField label="Cantidad de Visitas" htmlFor="plan-visit-limit">
                <input
                  id="plan-visit-limit"
                  required
                  type="number"
                  min="1"
                  value={form.visitLimit}
                  onChange={(event) => setForm((current) => ({ ...current, visitLimit: event.target.value }))}
                  className={fieldClassName}
                />
              </FormField>
              <FormField label="Vigencia" htmlFor="plan-visit-duration-value">
                <div className="grid grid-cols-[110px_1fr] gap-[10px]">
                  <input
                    id="plan-visit-duration-value"
                    required
                    type="number"
                    min="1"
                    value={form.visitDurationValue}
                    onChange={(event) => setForm((current) => ({ ...current, visitDurationValue: event.target.value }))}
                    className={fieldClassName}
                  />
                  <SelectField
                    aria-label="Unidad de vigencia"
                    value={form.visitDurationUnit}
                    onChange={(event) => setForm((current) => ({ ...current, visitDurationUnit: event.target.value as MembershipPlanInput["durationUnit"] }))}
                  >
                    <option value="days">días</option>
                    <option value="weeks">semanas</option>
                    <option value="months">meses</option>
                    <option value="years">años</option>
                  </SelectField>
                </div>
              </FormField>
            </div>
          ) : (
            <FormField label="Duración" htmlFor="plan-duration">
              <input id="plan-duration" disabled value={fixedDurationLabel(selectedTemplate)} className={`${fieldClassName} bg-[#f0f1f3] text-[#7e8792]`} />
            </FormField>
          )}

          <FormField label="Estado" htmlFor="plan-status">
            <SelectField
              id="plan-status"
              value={form.status}
              onChange={(event) => setForm((current) => ({ ...current, status: event.target.value as MembershipPlanInput["status"] }))}
            >
              <option value="active">Activo</option>
              <option value="inactive">Inactivo</option>
            </SelectField>
          </FormField>

          {error ? <p className="mt-[-7px] text-[17px] font-bold text-[var(--color-brand-danger)]" role="alert">{error}</p> : null}

          <footer className="mt-1 grid grid-cols-2 gap-5 border-t border-[#e1e4e8] pt-8 max-[560px]:grid-cols-1 max-[560px]:gap-3">
            <button
              type="button"
              className="min-h-[71px] rounded-xl border-0 bg-[#374357] text-[23px] font-extrabold text-white disabled:cursor-wait disabled:opacity-65 max-[560px]:min-h-[59px] max-[560px]:text-[19px]"
              onClick={onClose}
              disabled={submitting}
            >
              Cancelar
            </button>
            <button
              type="submit"
              className="min-h-[71px] rounded-xl border-0 bg-[linear-gradient(100deg,var(--color-brand-button-start)_0%,var(--color-brand-button-middle)_42%,var(--color-brand-button-end)_100%)] text-[23px] font-extrabold text-[var(--color-brand-button-text)] disabled:cursor-wait disabled:opacity-65 max-[560px]:min-h-[59px] max-[560px]:text-[19px]"
              disabled={submitting}
            >
              {submitting ? "Creando..." : "Crear Plan"}
            </button>
          </footer>
        </form>
      </section>
    </div>
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

const fieldClassName = "min-h-[76px] w-full rounded-xl border border-[#cbd0d7] bg-white px-[23px] text-[22px] text-[#172131] outline-0 focus:border-[var(--color-brand-accent)] focus:outline-3 focus:outline-[color-mix(in_srgb,var(--color-brand-accent)_32%,transparent)] max-[560px]:min-h-[58px] max-[560px]:text-[18px]";

function FormField(props: {
  label: string;
  htmlFor: string;
  children: ReactNode;
}) {
  return (
    <label className="grid gap-3 text-[21px] font-bold text-[#394351] max-[560px]:text-[18px]" htmlFor={props.htmlFor}>
      <span>{props.label}</span>
      {props.children}
    </label>
  );
}

function SelectField(props: SelectHTMLAttributes<HTMLSelectElement>) {
  const { children, className, ...rest } = props;
  return (
    <span className="relative block">
      <select
        {...rest}
        className={`${fieldClassName} appearance-none pr-16 ${className ?? ""}`.trim()}
      >
        {children}
      </select>
      <ChevronDownIcon className="pointer-events-none absolute right-5 top-1/2 h-6 w-6 -translate-y-1/2 text-[#a4acb6]" />
    </span>
  );
}
