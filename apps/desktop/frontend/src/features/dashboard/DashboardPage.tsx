import { useState, type ComponentType } from "react";
import { desktopVersion } from "@gym-saas/shared";
import {
  ActivityIcon,
  AlertCircleIcon,
  ArrowLeftIcon,
  AwardIcon,
  BarChartIcon,
  CakeIcon,
  CalendarIcon,
  CameraIcon,
  CashIcon,
  ChevronDownIcon,
  CreditCardIcon,
  FilterIcon,
  HistoryIcon,
  InfoCircleIcon,
  ListIcon,
  PlusIcon,
  ShoppingBagIcon,
  ShoppingCartIcon,
  UsersGroupIcon,
  UsersIcon
} from "@gym-saas/ui";
import logoSrc from "../../assets/logo.png";
import { DashboardRouteOutlet } from "./DashboardRouteOutlet";
import "./dashboard-page.css";

type DashboardPageProps = {
  onSignOut: () => void;
};

type NavigationItem = {
  id: string;
  label: string;
  icon: IconComponent;
  children?: Array<{ id: string; label: string; icon: IconComponent }>;
};

type IconComponent = ComponentType<{ className?: string }>;

const navigation: NavigationItem[] = [
  { id: "dashboard", label: "Panel de Control", icon: BarChartIcon },
  { id: "members", label: "Clientes", icon: UsersIcon },
  { id: "attendance", label: "Asistencias", icon: CalendarIcon },
  { id: "birthdays", label: "Cumpleañeros", icon: CakeIcon },
  {
    id: "memberships",
    label: "Membresías",
    icon: CreditCardIcon,
    children: [
      { id: "plans", label: "Planes", icon: ListIcon },
      { id: "membership-payments", label: "Pagos", icon: CashIcon }
    ]
  },
  {
    id: "store",
    label: "Tienda",
    icon: ShoppingBagIcon,
    children: [
      { id: "products", label: "Productos", icon: ShoppingBagIcon },
      { id: "sales", label: "Ventas", icon: ShoppingCartIcon }
    ]
  },
  { id: "users", label: "Gestión de Usuarios", icon: UsersGroupIcon },
  { id: "marketing", label: "Marketing", icon: AwardIcon },
  { id: "settings", label: "Configuración", icon: FilterIcon },
  { id: "monitoring", label: "Monitoreo", icon: CameraIcon },
  { id: "system-health", label: "Salud del Sistema", icon: ActivityIcon },
  { id: "logs", label: "Logs", icon: HistoryIcon }
];

const routeLabels = Object.fromEntries(
  navigation.flatMap((item) => [
    [item.id, item.label],
    ...(item.children?.map((child) => [child.id, child.label]) ?? [])
  ])
);

export function DashboardPage({ onSignOut }: DashboardPageProps) {
  const [activeRoute, setActiveRoute] = useState("dashboard");
  const [openGroups, setOpenGroups] = useState<Record<string, boolean>>({ memberships: true, store: true });

  function activate(route: string) {
    setActiveRoute(route);
  }

  function toggleGroup(group: string) {
    setOpenGroups((current) => ({ ...current, [group]: !current[group] }));
  }

  const activeLabel = routeLabels[activeRoute] ?? "Panel de Control";

  return (
    <main className="dashboard-page">
      <aside className="dashboard-sidebar">
        <div className="sidebar-brand">
          <strong>Control de Gimnasio</strong>
          <img src={logoSrc} alt="Zeus Gym" />
        </div>

        <nav className="sidebar-nav" aria-label="Navegación principal">
          {navigation.map((item) => {
            const isGroup = Boolean(item.children);
            const expanded = Boolean(openGroups[item.id]);
            const selected = activeRoute === item.id || item.children?.some((child) => child.id === activeRoute);

            return (
              <div className="nav-group" key={item.id}>
                <button
                  className={`nav-item ${selected && !isGroup ? "is-active" : ""}`}
                  type="button"
                  onClick={() => (isGroup ? toggleGroup(item.id) : activate(item.id))}
                  aria-expanded={isGroup ? expanded : undefined}
                >
                  <MenuGlyph Icon={item.icon} />
                  <span>{item.label}</span>
                  {isGroup ? <ChevronDownIcon className={`nav-chevron ${expanded ? "is-expanded" : ""}`} /> : null}
                </button>
                {isGroup && expanded ? (
                  <div className="nav-children">
                    {item.children?.map((child) => (
                      <button
                        className={`nav-child ${activeRoute === child.id ? "is-active" : ""}`}
                        key={child.id}
                        type="button"
                        onClick={() => activate(child.id)}
                      >
                        <MenuGlyph Icon={child.icon} />
                        <span>{child.label}</span>
                      </button>
                    ))}
                  </div>
                ) : null}
              </div>
            );
          })}
        </nav>

        <div className="sidebar-footer">
          <div className="profile-summary">
            <span className="profile-avatar">S</span>
            <span>
              <strong>Super Administrador</strong>
              <small>superadmin@zeus.gym</small>
            </span>
          </div>
          <span className="app-version">Versión {desktopVersion}</span>
          <button className="sign-out" type="button" onClick={onSignOut}>
            <ArrowLeftIcon className="sign-out-icon" /> Cerrar sesión
          </button>
        </div>
      </aside>

      <section className="dashboard-workspace">
        <DeviceBanner />
        <header className="workspace-toolbar">
          <button className="menu-button" type="button" aria-label="Abrir menú"><ListIcon className="toolbar-icon" /></button>
          <div className="breadcrumbs"><span aria-hidden="true">⌂</span><span>Inicio</span><b>/</b><strong>{activeLabel}</strong></div>
          <div className="toolbar-actions">
            <button className="help-button" type="button"><InfoCircleIcon className="button-icon" /> Ayuda</button>
            <button className="daily-button" type="button"><PlusIcon className="button-icon" /> Registrar Diario</button>
            <button className="monitor-button" type="button"><CameraIcon className="button-icon" /> Abrir monitor de Reconocimiento</button>
          </div>
        </header>

        <DashboardRouteOutlet route={activeRoute} title={activeLabel} />
      </section>
    </main>
  );
}

function DeviceBanner() {
  return (
    <section className="device-banner" aria-label="Estado del dispositivo">
      <AlertCircleIcon className="device-error-icon" />
      <div>
        <strong>Error de Conexión</strong>
        <p>No se pudo conectar al dispositivo: Host is down. Device may be offline or unreachable. Verifica la configuración en Configuración &gt; Dispositivo.</p>
      </div>
      <button type="button">Reintentar</button>
    </section>
  );
}

function MenuGlyph({ Icon }: { Icon: IconComponent }) {
  return <Icon className="menu-glyph" />;
}
