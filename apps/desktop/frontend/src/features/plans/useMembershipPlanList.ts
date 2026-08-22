import { useCallback, useEffect, useState } from "react";
import { isDesktopApp } from "../../platform/wails";
import { listMembershipPlans, type MembershipPlan } from "./api";

export function useMembershipPlanList() {
  const desktopRuntime = isDesktopApp();
  const [plans, setPlans] = useState<MembershipPlan[]>([]);
  const [loading, setLoading] = useState(desktopRuntime);
  const [error, setError] = useState<string | null>(null);

  const refresh = useCallback(async () => {
    if (!desktopRuntime) {
      setLoading(false);
      return;
    }

    setLoading(true);
    try {
      setPlans(await listMembershipPlans());
      setError(null);
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : "No se pudieron cargar los planes de membresía.");
    } finally {
      setLoading(false);
    }
  }, [desktopRuntime]);

  useEffect(() => {
    void refresh();
  }, [refresh]);

  return { plans, loading, error, refresh, desktopRuntime };
}
