import type { ComponentType } from "react";
import {
  ActivityIcon,
  AwardIcon,
  BarChartIcon,
  CakeIcon,
  CalendarIcon,
  CameraIcon,
  CashIcon,
  CreditCardIcon,
  FilterIcon,
  HistoryIcon,
  ListIcon,
  ShoppingBagIcon,
  ShoppingCartIcon,
  UsersGroupIcon,
  UsersIcon
} from "@gym-saas/ui";
import type { UserRole } from "../auth/mockAuth";

export type DashboardRoute =
  | "dashboard"
  | "members"
  | "attendance"
  | "birthdays"
  | "memberships"
  | "plans"
  | "membership-payments"
  | "store"
  | "products"
  | "sales"
  | "users"
  | "marketing"
  | "settings"
  | "monitoring"
  | "system-health"
  | "logs";

export type IconComponent = ComponentType<{ className?: string }>;

type NavigationItem = {
  id: DashboardRoute;
  label: string;
  icon: IconComponent;
  allowedRoles: readonly UserRole[];
  children?: readonly NavigationItem[];
};

const allRoles: readonly UserRole[] = ["super_admin", "gym_admin", "receptionist"];
const managementRoles: readonly UserRole[] = ["super_admin", "gym_admin"];
const superAdminOnly: readonly UserRole[] = ["super_admin"];

export const navigation: readonly NavigationItem[] = [
  { id: "dashboard", label: "Panel de Control", icon: BarChartIcon, allowedRoles: allRoles },
  { id: "members", label: "Clientes", icon: UsersIcon, allowedRoles: allRoles },
  { id: "attendance", label: "Asistencias", icon: CalendarIcon, allowedRoles: allRoles },
  { id: "birthdays", label: "Cumpleañeros", icon: CakeIcon, allowedRoles: allRoles },
  {
    id: "memberships",
    label: "Membresías",
    icon: CreditCardIcon,
    allowedRoles: allRoles,
    children: [
      { id: "plans", label: "Planes", icon: ListIcon, allowedRoles: managementRoles },
      { id: "membership-payments", label: "Pagos", icon: CashIcon, allowedRoles: allRoles }
    ]
  },
  {
    id: "store",
    label: "Tienda",
    icon: ShoppingBagIcon,
    allowedRoles: managementRoles,
    children: [
      { id: "products", label: "Productos", icon: ShoppingBagIcon, allowedRoles: managementRoles },
      { id: "sales", label: "Ventas", icon: ShoppingCartIcon, allowedRoles: managementRoles }
    ]
  },
  { id: "users", label: "Gestión de Usuarios", icon: UsersGroupIcon, allowedRoles: managementRoles },
  { id: "marketing", label: "Marketing", icon: AwardIcon, allowedRoles: managementRoles },
  { id: "settings", label: "Configuración", icon: FilterIcon, allowedRoles: managementRoles },
  { id: "monitoring", label: "Monitoreo", icon: CameraIcon, allowedRoles: managementRoles },
  { id: "system-health", label: "Salud del Sistema", icon: ActivityIcon, allowedRoles: superAdminOnly },
  { id: "logs", label: "Logs", icon: HistoryIcon, allowedRoles: superAdminOnly }
];

const routeItems = navigation.flatMap((item) => [item, ...(item.children ?? [])]);

export const routeLabels: Readonly<Record<DashboardRoute, string>> = Object.fromEntries(
  routeItems.map(({ id, label }) => [id, label])
) as Record<DashboardRoute, string>;

export function canAccessRoute(role: UserRole, route: DashboardRoute): boolean {
  return routeItems.some((item) => item.id === route && item.allowedRoles.includes(role));
}

export function visibleNavigation(role: UserRole): readonly NavigationItem[] {
  return navigation.flatMap((item) => {
    if (!item.allowedRoles.includes(role)) {
      return [];
    }

    const children = item.children?.filter((child) => child.allowedRoles.includes(role));
    if (item.children && !children?.length) {
      return [];
    }

    return [{ ...item, children }];
  });
}

export function firstAccessibleRoute(role: UserRole): DashboardRoute {
  const route = routeItems.find((item) => item.allowedRoles.includes(role));
  if (!route) {
    throw new Error(`No dashboard route is configured for role ${role}`);
  }

  return route.id;
}
