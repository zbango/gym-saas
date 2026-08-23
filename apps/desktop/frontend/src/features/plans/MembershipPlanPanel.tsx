import { useState } from "react";
import { Panel, ShellButton } from "@gym-saas/ui";
import type { MembershipPlan } from "./api";
import { useMembershipPlans } from "./useMembershipPlans";

const fieldClassName = "min-h-11 rounded-xl border border-[var(--color-brand-border)] bg-[var(--color-brand-input)] px-3 py-2.5 text-[var(--color-brand-text)] outline-0";
const labelClassName = "grid gap-1.5 font-semibold";

export function MembershipPlanPanel() {
  const { plans, form, editingPlanID, error, saving, setField, setValidityKind, startEditing, cancelEditing, save, archive } = useMembershipPlans();
  const [archiveCandidate, setArchiveCandidate] = useState<MembershipPlan | null>(null);
  const [archiving, setArchiving] = useState(false);

  return (
    <Panel title="Membership Plans" eyebrow="Plan catalog">
      <p className="mt-0 text-[var(--color-brand-muted)]">
        {form.validityKind === "visits"
          ? "Visit plans expire when all visits are used or their calendar window ends—whichever happens first."
          : "Period plans expire at the end of their configured calendar duration."}
      </p>
      <div className="grid gap-3 [grid-template-columns:repeat(auto-fit,minmax(210px,1fr))]">
        <PlanField label="Plan name" value={form.name} onChange={(value) => setField("name", value)} />
        <label className={labelClassName}>
          Plan type
          <select className={fieldClassName} value={form.validityKind} onChange={(event) => setValidityKind(event.target.value as "time" | "visits")}>
            <option value="time">Period</option>
            <option value="visits">Visits with expiry</option>
          </select>
        </label>
        <PlanField label="Validity duration" type="number" min={1} value={String(form.durationValue)} onChange={(value) => setField("durationValue", Number(value))} />
        <label className={labelClassName}>
          Duration unit
          <select className={fieldClassName} value={form.durationUnit} onChange={(event) => setField("durationUnit", event.target.value as typeof form.durationUnit)}>
            <option value="days">Days</option>
            <option value="weeks">Weeks</option>
            <option value="months">Months</option>
            <option value="years">Years</option>
          </select>
        </label>
        {form.validityKind === "visits" ? <PlanField label="Included visits" type="number" min={1} value={String(form.visitLimit)} onChange={(value) => setField("visitLimit", Number(value))} /> : null}
        <PlanField label="Price (cents)" type="number" min={0} value={String(form.priceCents)} onChange={(value) => setField("priceCents", Number(value))} />
        <label className={labelClassName}>
          Status
          <select className={fieldClassName} value={form.status} onChange={(event) => setField("status", event.target.value as typeof form.status)}>
            <option value="active">Active</option>
            <option value="inactive">Inactive</option>
          </select>
        </label>
      </div>
      {error ? <p className="mb-0 text-[var(--color-brand-danger)]">{error}</p> : null}
      <div className="mt-4 flex flex-wrap gap-2">
        <ShellButton onClick={() => void save()} disabled={saving}>
          {saving ? "Saving..." : editingPlanID ? "Save changes" : "Create plan"}
        </ShellButton>
        {editingPlanID ? <ShellButton variant="secondary" onClick={cancelEditing}>Cancel</ShellButton> : null}
      </div>
      <div className="mt-5 grid gap-2">
        {plans.length === 0 ? (
          <p className="m-0 text-[var(--color-brand-muted)]">No active membership plans yet.</p>
        ) : plans.map((plan) => (
          <div key={plan.id} className="flex flex-wrap items-center justify-between gap-3 border-t border-[var(--color-brand-border)] pt-3">
            <div>
              <strong>{plan.name}</strong>
              <div className="text-sm text-[var(--color-brand-muted)]">
                {plan.validityKind === "visits" ? `${plan.visitLimit} visits · ` : ""}{plan.durationValue} {plan.durationUnit} · {formatMoney(plan.priceCents, plan.currency)} · {plan.status}
              </div>
            </div>
            <div className="flex gap-2">
              <ShellButton variant="secondary" onClick={() => startEditing(plan)}>Edit</ShellButton>
              <ShellButton variant="secondary" onClick={() => setArchiveCandidate(plan)}>Archive</ShellButton>
            </div>
          </div>
        ))}
      </div>
      {archiveCandidate ? (
        <div className="mt-5 grid gap-2.5 rounded-xl bg-[var(--color-brand-accent-soft)] p-4">
          <strong>Archive {archiveCandidate.name}?</strong>
          <span>This removes the plan from the active catalog without deleting historical records.</span>
          <div className="flex gap-2">
            <ShellButton variant="secondary" onClick={() => setArchiveCandidate(null)} disabled={archiving}>Cancel</ShellButton>
            <ShellButton
              disabled={archiving}
              onClick={async () => {
                setArchiving(true);
                try {
                  if (await archive(archiveCandidate)) {
                    setArchiveCandidate(null);
                  }
                } finally {
                  setArchiving(false);
                }
              }}
            >
              {archiving ? "Archiving..." : "Archive plan"}
            </ShellButton>
          </div>
        </div>
      ) : null}
    </Panel>
  );
}

function PlanField(props: { label: string; value: string; onChange: (value: string) => void; type?: string; min?: number }) {
  return (
    <label className={labelClassName}>
      {props.label}
      <input className={fieldClassName} type={props.type ?? "text"} min={props.min} value={props.value} onChange={(event) => props.onChange(event.target.value)} />
    </label>
  );
}

function formatMoney(cents: number, currency: string): string {
  return new Intl.NumberFormat(undefined, { style: "currency", currency }).format(cents / 100);
}
