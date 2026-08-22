import { useEffect, useState } from "react";
import { isDesktopApp } from "../../platform/wails";
import {
  archiveMembershipPlan,
  createMembershipPlan,
  listMembershipPlans,
  type MembershipPlan,
  type MembershipPlanInput,
  updateMembershipPlan
} from "./api";

const emptyPlan: MembershipPlanInput = {
  name: "",
  validityKind: "time",
  durationValue: 1,
  durationUnit: "months",
  visitLimit: 0,
  priceCents: 0,
  currency: "USD",
  status: "active"
};

export function useMembershipPlans() {
  const [plans, setPlans] = useState<MembershipPlan[]>([]);
  const [form, setForm] = useState<MembershipPlanInput>(emptyPlan);
  const [editingPlanID, setEditingPlanID] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);
  const desktopRuntime = isDesktopApp();

  async function refresh() {
    if (!desktopRuntime) {
      return;
    }
    try {
      setPlans(await listMembershipPlans());
      setError(null);
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : "Unable to load membership plans");
    }
  }

  useEffect(() => {
    void refresh();
  }, []);

  function setField<Field extends keyof MembershipPlanInput>(field: Field, value: MembershipPlanInput[Field]) {
    setForm((current) => ({ ...current, [field]: value }));
  }

  function setValidityKind(validityKind: MembershipPlanInput["validityKind"]) {
    setForm((current) => ({ ...current, validityKind, visitLimit: validityKind === "time" ? 0 : current.visitLimit }));
  }

  function startEditing(plan: MembershipPlan) {
    setEditingPlanID(plan.id);
    setForm({
      name: plan.name,
      validityKind: plan.validityKind,
      durationValue: plan.durationValue,
      durationUnit: plan.durationUnit,
      visitLimit: plan.visitLimit,
      priceCents: plan.priceCents,
      currency: plan.currency,
      status: plan.status
    });
    setError(null);
  }

  function cancelEditing() {
    setEditingPlanID(null);
    setForm(emptyPlan);
    setError(null);
  }

  async function save() {
    if (!desktopRuntime) {
      setError("Membership plans are only available in the desktop app");
      return false;
    }
    setSaving(true);
    setError(null);
    try {
      if (editingPlanID) {
        await updateMembershipPlan(editingPlanID, form);
      } else {
        await createMembershipPlan(form);
      }
      cancelEditing();
      await refresh();
      return true;
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : "Unable to save membership plan");
      return false;
    } finally {
      setSaving(false);
    }
  }

  async function archive(plan: MembershipPlan) {
    setError(null);
    try {
      await archiveMembershipPlan(plan.id);
      if (editingPlanID === plan.id) {
        cancelEditing();
      }
      await refresh();
      return true;
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : "Unable to archive membership plan");
      return false;
    }
  }

  return { plans, form, editingPlanID, error, saving, setField, setValidityKind, startEditing, cancelEditing, save, archive };
}
