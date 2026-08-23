import { useState, type ComponentType, type ReactNode } from "react";
import {
  ActivityIcon,
  AlertTriangleIcon,
  CalendarIcon,
  ChevronDownIcon,
  DocumentIcon,
  DollarSignIcon,
  RefreshIcon
} from "@gym-saas/ui";

type IconComponent = ComponentType<{ className?: string }>;

export function DashboardHomePage() {
  const [period, setPeriod] = useState<"Diario" | "Mensual">("Diario");
  const [alertsOpen, setAlertsOpen] = useState(false);

  return (
    <div className="mx-auto grid max-w-[1550px] gap-5 px-[14px] py-[22px] min-[621px]:gap-[31px] min-[621px]:px-[22px] min-[621px]:py-[30px] min-[921px]:px-[42px] min-[921px]:py-[42px] min-[1261px]:px-[64px] min-[1261px]:py-[64px]">
      <section className="grid items-start gap-7 min-[921px]:flex min-[921px]:justify-between">
        <div>
          <h1 className="m-0 text-[30px] font-bold tracking-[-0.04em] text-[#111827] min-[921px]:text-[clamp(31px,2.35vw,42px)]">
            Bienvenido al panel de gestión de tu gimnasio Zeus Gym
          </h1>
          <p className="mt-[15px] text-[17px] text-[#66717e] min-[921px]:text-[21px]">
            Hoy: sábado, 22 de agosto de 2026
          </p>
        </div>
        <button
          className="inline-flex min-h-[49px] items-center gap-2 self-start rounded-[11px] bg-[linear-gradient(100deg,var(--color-brand-button-start)_0%,var(--color-brand-button-middle)_42%,var(--color-brand-button-end)_100%)] px-[19px] text-[16px] font-extrabold whitespace-nowrap text-[var(--color-brand-button-text)]"
          type="button"
        >
          <RefreshIcon className="h-[22px] w-[22px] shrink-0" /> Actualizar Datos
        </button>
      </section>

      <section className={cardShellClassName}>
        <CardHeading
          Icon={DollarSignIcon}
          title="Resumen Financiero"
          subtitle="sábado, 22 de agosto de 2026"
          action={<RefreshIcon className="h-[25px] w-[25px]" />}
        />
        <div className="flex min-h-[108px] flex-col items-start justify-between gap-[17px] border-b border-b-[#e7e8eb] px-[18px] py-[18px] min-[621px]:px-[34px] min-[621px]:py-[22px] min-[921px]:flex-row min-[921px]:items-center">
          <div className="flex rounded-[11px] border border-[#d8dbe0] p-[5px]" aria-label="Periodo financiero">
            {(["Diario", "Mensual"] as const).map((option) => (
              <button
                className={[
                  "min-h-[44px] rounded-[7px] px-[19px] text-[17px] font-bold text-[#495462]",
                  period === option ? "bg-[#4e83e9] text-white" : "bg-transparent"
                ].join(" ")}
                type="button"
                key={option}
                onClick={() => setPeriod(option)}
              >
                {option}
              </button>
            ))}
          </div>
          <label className="flex min-h-[45px] items-center gap-[11px] rounded-[10px] border border-[#d8dbe0] px-[13px] min-[621px]:min-h-[54px]">
            <CalendarIcon className="h-[22px] w-[22px]" />
            <input
              className="border-0 bg-transparent text-[17px] font-semibold text-[#20242a] outline-none"
              aria-label="Fecha del resumen financiero"
              type="date"
              defaultValue="2026-08-22"
            />
          </label>
        </div>
        <EmptyState>No hay actividad financiera registrada en este {period === "Diario" ? "día" : "mes"}</EmptyState>
      </section>

      <section className={cardShellClassName}>
        <CardHeading Icon={AlertTriangleIcon} title="Alertas" iconClassName="text-[#66717e]" />
        <button
          className="flex min-h-[76px] w-full items-center gap-[17px] border-l-[6px] border-l-[#ef5b55] bg-[var(--color-brand-danger-soft)] px-[30px] text-left text-[19px] text-[#a13e35]"
          type="button"
          onClick={() => setAlertsOpen((open) => !open)}
          aria-expanded={alertsOpen}
        >
          <DocumentIcon className="h-6 w-6 shrink-0" />
          <strong>20 clientes tienen pagos pendientes</strong>
          <ChevronDownIcon className="ml-auto h-5 w-5 text-[#2d323a]" />
        </button>
        {alertsOpen ? (
          <div className="bg-white px-[35px] py-5 text-[#67717e]">
            Este bloque mostrará el detalle de pagos pendientes cuando el backend esté integrado.
          </div>
        ) : null}
      </section>

      <section className={cardShellClassName}>
        <CardHeading Icon={DocumentIcon} title="Pagos Pendientes" action="Ver más" />
        <EmptyState>No hay pagos pendientes del último mes</EmptyState>
      </section>

      <section className={`${cardShellClassName} mb-5`}>
        <CardHeading Icon={ActivityIcon} title="Check-ins Recientes" />
        <EmptyState>Los próximos ingresos aparecerán aquí</EmptyState>
      </section>
    </div>
  );
}

const cardShellClassName =
  "overflow-hidden rounded-[20px] border border-[#eff0f2] bg-white shadow-[0_2px_4px_rgba(23,31,44,0.025)] min-[621px]:rounded-[31px]";

function CardHeading(props: {
  Icon: IconComponent;
  title: string;
  subtitle?: string;
  action?: ReactNode;
  iconClassName?: string;
}) {
  return (
    <header className="flex min-h-[106px] items-center gap-[18px] border-b border-b-[#e7e8eb] px-[18px] py-[18px] min-[621px]:px-[34px] min-[621px]:py-[22px]">
      <props.Icon className={["h-[30px] w-[30px] shrink-0 text-[var(--color-brand-success)]", props.iconClassName ?? ""].join(" ")} />
      <div>
        <h2 className="m-0 text-[23px] font-bold text-[#172131]">{props.title}</h2>
        {props.subtitle ? <p className="mt-[5px] text-[17px] text-[#6e7783]">{props.subtitle}</p> : null}
      </div>
      {props.action ? (
        <button className="ml-auto bg-transparent text-[17px] font-extrabold text-[#4d78df]" type="button">
          {props.action}
        </button>
      ) : null}
    </header>
  );
}

function EmptyState({ children }: { children: ReactNode }) {
  return (
    <div className="grid min-h-[96px] place-items-center px-[22px] py-[22px] text-center text-[16px] text-[#747d88] min-[621px]:min-h-[115px] min-[621px]:text-[19px]">
      {children}
    </div>
  );
}
