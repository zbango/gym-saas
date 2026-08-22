import type { ComponentType } from "react";
import { MemberListPage } from "../members/MemberListPage";
import { MembershipPlanListPage } from "../plans/MembershipPlanListPage";
import { DashboardHomePage } from "./DashboardHomePage";

const routeViews: Readonly<Record<string, ComponentType>> = {
  dashboard: DashboardHomePage,
  members: MemberListPage,
  plans: MembershipPlanListPage
};

export function DashboardRouteOutlet({ route, title }: { route: string; title: string }) {
  const RouteView = routeViews[route];
  return RouteView ? <RouteView /> : <MockRoute title={title} />;
}

function MockRoute({ title }: { title: string }) {
  return (
    <div className="mock-route">
      <h1>{title}</h1>
      <p>Esta es una vista de navegación simulada. El contenido real se conectará a su caso de uso de Go después de completar la paridad visual.</p>
    </div>
  );
}
