import { useState } from "react";
import { Panel, ShellButton } from "@gym-saas/ui";
import type { MembershipPlan } from "./api";
import { useMembershipPlans } from "./useMembershipPlans";

const inputStyle = {
  padding: "10px 12px",
  background: "var(--gs-input-background)",
  color: "var(--gs-input-text)",
  border: "1px solid var(--gs-border)",
  borderRadius: 12,
  font: "inherit"
};

export function MembershipPlanPanel() {
  const { plans, form, editingPlanID, error, saving, setField, setValidityKind, startEditing, cancelEditing, save, archive } = useMembershipPlans();
  const [archiveCandidate, setArchiveCandidate] = useState<MembershipPlan | null>(null);
  const [archiving, setArchiving] = useState(false);

  return (
    <Panel title="Membership Plans" eyebrow="Plan catalog">
      <p style={{ marginTop: 0, color: "var(--gs-text-muted)" }}>
        {form.validityKind === "visits"
          ? "Visit plans expire when all visits are used or their calendar window ends—whichever happens first."
          : "Period plans expire at the end of their configured calendar duration."}
      </p>
      <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fit, minmax(210px, 1fr))", gap: 12 }}>
        <PlanField label="Plan name" value={form.name} onChange={(value) => setField("name", value)} />
        <label style={{ display: "grid", gap: 6, fontWeight: 600 }}>
          Plan type
          <select value={form.validityKind} onChange={(event) => setValidityKind(event.target.value as "time" | "visits")} style={inputStyle}>
            <option value="time">Period</option>
            <option value="visits">Visits with expiry</option>
          </select>
        </label>
        <PlanField label="Validity duration" type="number" min={1} value={String(form.durationValue)} onChange={(value) => setField("durationValue", Number(value))} />
        <label style={{ display: "grid", gap: 6, fontWeight: 600 }}>
          Duration unit
          <select value={form.durationUnit} onChange={(event) => setField("durationUnit", event.target.value as typeof form.durationUnit)} style={inputStyle}>
            <option value="days">Days</option>
            <option value="weeks">Weeks</option>
            <option value="months">Months</option>
            <option value="years">Years</option>
          </select>
        </label>
        {form.validityKind === "visits" ? <PlanField label="Included visits" type="number" min={1} value={String(form.visitLimit)} onChange={(value) => setField("visitLimit", Number(value))} /> : null}
        <PlanField label="Price (cents)" type="number" min={0} value={String(form.priceCents)} onChange={(value) => setField("priceCents", Number(value))} />
        <label style={{ display: "grid", gap: 6, fontWeight: 600 }}>
          Status
          <select value={form.status} onChange={(event) => setField("status", event.target.value as typeof form.status)} style={inputStyle}>
            <option value="active">Active</option>
            <option value="inactive">Inactive</option>
          </select>
        </label>
      </div>
      {error ? <p style={{ color: "var(--gs-danger, #b42318)", marginBottom: 0 }}>{error}</p> : null}
      <div style={{ display: "flex", gap: 8, marginTop: 16, flexWrap: "wrap" }}>
        <ShellButton onClick={() => void save()} disabled={saving}>
          {saving ? "Saving..." : editingPlanID ? "Save changes" : "Create plan"}
        </ShellButton>
        {editingPlanID ? <ShellButton variant="secondary" onClick={cancelEditing}>Cancel</ShellButton> : null}
      </div>
      <div style={{ display: "grid", gap: 8, marginTop: 20 }}>
        {plans.length === 0 ? (
          <p style={{ color: "var(--gs-text-muted)", margin: 0 }}>No active membership plans yet.</p>
        ) : plans.map((plan) => (
          <div key={plan.id} style={{ display: "flex", gap: 12, justifyContent: "space-between", alignItems: "center", flexWrap: "wrap", borderTop: "1px solid var(--gs-border)", paddingTop: 12 }}>
            <div>
              <strong>{plan.name}</strong>
              <div style={{ color: "var(--gs-text-muted)", fontSize: 14 }}>
                {plan.validityKind === "visits" ? `${plan.visitLimit} visits · ` : ""}{plan.durationValue} {plan.durationUnit} · {formatMoney(plan.priceCents, plan.currency)} · {plan.status}
              </div>
            </div>
            <div style={{ display: "flex", gap: 8 }}>
              <ShellButton variant="secondary" onClick={() => startEditing(plan)}>Edit</ShellButton>
              <ShellButton variant="secondary" onClick={() => setArchiveCandidate(plan)}>Archive</ShellButton>
            </div>
          </div>
        ))}
      </div>
      {archiveCandidate ? (
        <div style={{ marginTop: 20, padding: 16, borderRadius: 12, background: "var(--gs-accent-soft)", display: "grid", gap: 10 }}>
          <strong>Archive {archiveCandidate.name}?</strong>
          <span>This removes the plan from the active catalog without deleting historical records.</span>
          <div style={{ display: "flex", gap: 8 }}>
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
    <label style={{ display: "grid", gap: 6, fontWeight: 600 }}>
      {props.label}
      <input type={props.type ?? "text"} min={props.min} value={props.value} onChange={(event) => props.onChange(event.target.value)} style={inputStyle} />
    </label>
  );
}

function formatMoney(cents: number, currency: string): string {
  return new Intl.NumberFormat(undefined, { style: "currency", currency }).format(cents / 100);
}
