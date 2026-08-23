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
    <div className="min-h-[430px] overflow-hidden rounded-[20px] border border-[#eff0f2] bg-white px-[22px] py-[30px] shadow-[0_2px_4px_rgba(23,31,44,0.025)] min-[621px]:rounded-[31px] min-[621px]:px-[55px] min-[621px]:py-[55px]">
      <h1 className="m-0 text-[36px] font-bold tracking-[-0.04em] text-[#172131]">{title}</h1>
      <p className="mb-0 mt-5 max-w-[670px] text-[19px] leading-[1.5] text-[#68727e]">
        Esta es una vista de navegación simulada. El contenido real se conectará a su caso de uso de Go después de
        completar la paridad visual.
      </p>
    </div>
  );
}
