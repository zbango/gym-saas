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
    <main className="grid min-h-screen w-full grid-cols-1 bg-[var(--color-brand-workspace)] text-[#1d2430] min-[921px]:[zoom:0.75] min-[921px]:grid-cols-[276px_minmax(0,1fr)] min-[1261px]:grid-cols-[338px_minmax(0,1fr)]">
      <aside className="flex h-auto flex-col overflow-y-auto border-r border-r-[color-mix(in_srgb,var(--color-brand-accent)_28%,transparent)] bg-[var(--color-brand-sidebar)] text-[var(--color-brand-sidebar-text)] min-[921px]:sticky min-[921px]:top-0 min-[921px]:h-[133.333vh]">
        <div className="flex min-h-[110px] items-center justify-center gap-[17px] border-b border-b-[color-mix(in_srgb,var(--color-brand-accent)_28%,transparent)] px-[22px] py-[22px] min-[921px]:min-h-[164px] min-[1261px]:px-[29px]">
          <strong className="w-[130px] text-[22px] font-bold leading-[1.3] text-[var(--color-brand-accent)] min-[1261px]:text-[26px]">
            Control de Gimnasio
          </strong>
          <img className="h-[84px] w-[84px] object-contain min-[1261px]:h-[105px] min-[1261px]:w-[105px]" src={logoSrc} alt="Zeus Gym" />
        </div>

        <nav
          className="grid max-h-[255px] flex-1 content-start gap-0.5 overflow-y-auto px-4 py-4 min-[921px]:max-h-none"
          aria-label="Navegación principal"
        >
          {navigation.map((item) => {
            const isGroup = Boolean(item.children);
            const expanded = Boolean(openGroups[item.id]);
            const selected = activeRoute === item.id || item.children?.some((child) => child.id === activeRoute);
            const navItemClass = [
              "flex min-h-[47px] w-full items-center gap-4 rounded-[11px] bg-transparent px-[22px] text-left text-[16px] font-bold text-[color-mix(in_srgb,var(--color-brand-sidebar-text)_86%,transparent)] transition-colors hover:bg-[color-mix(in_srgb,var(--color-brand-accent)_9%,transparent)] hover:text-[var(--color-brand-accent)] min-[621px]:text-[18px]",
              selected && !isGroup
                ? "bg-[color-mix(in_srgb,var(--color-brand-accent)_9%,transparent)] text-[var(--color-brand-accent)] outline outline-3 -outline-offset-3 outline-[var(--color-brand-accent)]"
                : ""
            ].join(" ");
            const navChildClass = (isActive: boolean) =>
              [
                "flex min-h-10 w-full items-center gap-[15px] bg-transparent text-left text-[17px] font-bold text-[color-mix(in_srgb,var(--color-brand-sidebar-text)_70%,transparent)] transition-colors hover:text-[var(--color-brand-accent)]",
                isActive ? "text-[var(--color-brand-accent)]" : ""
              ].join(" ");

            return (
              <div className="grid gap-0.5" key={item.id}>
                <button
                  className={navItemClass}
                  type="button"
                  onClick={() => (isGroup ? toggleGroup(item.id) : activate(item.id))}
                  aria-expanded={isGroup ? expanded : undefined}
                >
                  <MenuGlyph Icon={item.icon} />
                  <span>{item.label}</span>
                  {isGroup ? (
                    <ChevronDownIcon
                      className={[
                        "ml-auto h-5 w-5 transition-transform duration-150",
                        expanded ? "rotate-180" : ""
                      ].join(" ")}
                    />
                  ) : null}
                </button>
                {isGroup && expanded ? (
                  <div className="mb-[5px] ml-[34px] grid gap-[3px] border-l-4 border-l-[color-mix(in_srgb,var(--color-brand-accent)_25%,transparent)] py-1 pl-[22px]">
                    {item.children?.map((child) => (
                      <button
                        className={navChildClass(activeRoute === child.id)}
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

        <div className="mt-auto hidden gap-[14px] border-t border-t-[color-mix(in_srgb,var(--color-brand-accent)_28%,transparent)] px-4 py-5 min-[921px]:grid">
          <div className="flex items-center gap-3">
            <span className="grid h-[43px] w-[43px] shrink-0 place-items-center rounded-full border border-[var(--color-brand-accent)] bg-[color-mix(in_srgb,var(--color-brand-accent)_18%,transparent)] text-[20px] font-extrabold text-[var(--color-brand-accent)]">
              S
            </span>
            <span className="grid min-w-0 gap-1">
              <strong className="text-[16px]">Super Administrador</strong>
              <small className="text-[13px] text-[color-mix(in_srgb,var(--color-brand-sidebar-text)_66%,transparent)]">
                superadmin@zeus.gym
              </small>
            </span>
          </div>
          <span className="pl-[55px] text-[13px] text-[color-mix(in_srgb,var(--color-brand-sidebar-text)_66%,transparent)]">
            Versión {desktopVersion}
          </span>
          <button
            className="inline-flex min-h-[49px] items-center justify-center gap-2 rounded-[10px] border border-[color-mix(in_srgb,var(--color-brand-accent)_30%,transparent)] bg-transparent text-[16px] font-bold text-[var(--color-brand-sidebar-text)] transition-colors hover:bg-[color-mix(in_srgb,var(--color-brand-accent)_12%,transparent)]"
            type="button"
            onClick={onSignOut}
          >
            <ArrowLeftIcon className="h-[21px] w-[21px]" /> Cerrar sesión
          </button>
        </div>
      </aside>

      <section className="min-w-0 bg-[var(--color-brand-workspace)]">
        <DeviceBanner />
        <header className="flex flex-wrap items-center gap-7 border-b border-b-[#e6e7ea] bg-white px-[22px] py-4 min-[921px]:min-h-[98px] min-[921px]:px-[26px] min-[1261px]:px-[46px] min-[921px]:py-0">
          <button className="bg-transparent text-[#5b6570]" type="button" aria-label="Abrir menú">
            <ListIcon className="h-[27px] w-[27px]" />
          </button>
          <div className="flex items-center gap-3 text-[15px] min-[621px]:text-[17px]">
            <span aria-hidden="true">⌂</span>
            <span>Inicio</span>
            <b className="text-[#a0a5aa]">/</b>
            <strong>{activeLabel}</strong>
          </div>
          <div className="flex w-full gap-[10px] overflow-x-auto min-[921px]:ml-auto min-[921px]:w-auto min-[1261px]:gap-4">
            <button
              className="inline-flex min-h-[49px] items-center gap-2 rounded-[11px] border border-[var(--color-brand-accent)] bg-[#fff7cc] px-[19px] text-[13px] font-extrabold text-[#4e535d] whitespace-nowrap min-[621px]:text-[16px]"
              type="button"
            >
              <InfoCircleIcon className="h-[22px] w-[22px] shrink-0" /> Ayuda
            </button>
            <button
              className="inline-flex min-h-[49px] items-center gap-2 rounded-[11px] bg-[var(--color-brand-success)] px-[19px] text-[13px] font-extrabold text-white whitespace-nowrap min-[621px]:text-[16px]"
              type="button"
            >
              <PlusIcon className="h-[22px] w-[22px] shrink-0" /> Registrar Diario
            </button>
            <button
              className="inline-flex min-h-[49px] items-center gap-2 rounded-[11px] bg-[#343c4c] px-[19px] text-[13px] font-extrabold text-white whitespace-nowrap min-[621px]:text-[16px] min-[921px]:max-[1260px]:text-[0px]"
              type="button"
            >
              <CameraIcon className="h-[22px] w-[22px] shrink-0" /> Abrir monitor de Reconocimiento
            </button>
          </div>
        </header>

        <DashboardRouteOutlet route={activeRoute} title={activeLabel} />
      </section>
    </main>
  );
}

function DeviceBanner() {
  return (
    <section
      className="grid grid-cols-[auto_minmax(0,1fr)] items-center gap-[18px] border-b border-b-[color-mix(in_srgb,var(--color-brand-danger)_24%,transparent)] bg-[var(--color-brand-danger-soft)] px-[22px] py-4 text-[#7d3732] min-[621px]:grid-cols-[auto_minmax(0,1fr)_auto] min-[921px]:min-h-[115px] min-[921px]:px-[31px] min-[921px]:py-[18px]"
      aria-label="Estado del dispositivo"
    >
      <AlertCircleIcon className="h-[27px] w-[27px] text-current" />
      <div>
        <strong className="text-[19px]">Error de Conexión</strong>
        <p className="mt-1.5 max-w-[1390px] text-[17px] leading-[1.4]">
          No se pudo conectar al dispositivo: Host is down. Device may be offline or unreachable. Verifica la
          configuración en Configuración &gt; Dispositivo.
        </p>
      </div>
      <button
        className="col-start-2 justify-self-start rounded-[11px] bg-[var(--color-brand-danger)] px-[25px] py-[14px] text-[16px] font-extrabold text-white min-[621px]:col-start-auto"
        type="button"
      >
        Reintentar
      </button>
    </section>
  );
}

function MenuGlyph({ Icon }: { Icon: IconComponent }) {
  return <Icon className="h-[26px] w-[26px] shrink-0 text-current" />;
}
