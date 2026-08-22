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
    <div className="dashboard-content">
      <section className="dashboard-welcome">
        <div>
          <h1>Bienvenido al panel de gestión de tu gimnasio Zeus Gym</h1>
          <p>Hoy: sábado, 22 de agosto de 2026</p>
        </div>
        <button className="refresh-data" type="button"><RefreshIcon className="button-icon" /> Actualizar Datos</button>
      </section>

      <section className="dashboard-card financial-card">
        <CardHeading Icon={DollarSignIcon} title="Resumen Financiero" subtitle="sábado, 22 de agosto de 2026" action={<RefreshIcon className="card-action-icon" />} />
        <div className="financial-controls">
          <div className="period-toggle" aria-label="Periodo financiero">
            {(["Diario", "Mensual"] as const).map((option) => (
              <button className={period === option ? "is-selected" : ""} type="button" key={option} onClick={() => setPeriod(option)}>{option}</button>
            ))}
          </div>
          <label className="date-control"><CalendarIcon className="date-icon" /><input aria-label="Fecha del resumen financiero" type="date" defaultValue="2026-08-22" /></label>
        </div>
        <EmptyState>No hay actividad financiera registrada en este {period === "Diario" ? "día" : "mes"}</EmptyState>
      </section>

      <section className="dashboard-card alert-card">
        <CardHeading Icon={AlertTriangleIcon} title="Alertas" />
        <button className="alert-row" type="button" onClick={() => setAlertsOpen((open) => !open)} aria-expanded={alertsOpen}>
          <DocumentIcon className="alert-document-icon" /><strong>20 clientes tienen pagos pendientes</strong><ChevronDownIcon className="alert-chevron" />
        </button>
        {alertsOpen ? <div className="alert-detail">Este bloque mostrará el detalle de pagos pendientes cuando el backend esté integrado.</div> : null}
      </section>

      <section className="dashboard-card pending-card">
        <CardHeading Icon={DocumentIcon} title="Pagos Pendientes" action="Ver más" />
        <EmptyState>No hay pagos pendientes del último mes</EmptyState>
      </section>

      <section className="dashboard-card checkin-card">
        <CardHeading Icon={ActivityIcon} title="Check-ins Recientes" />
        <EmptyState>Los próximos ingresos aparecerán aquí</EmptyState>
      </section>
    </div>
  );
}

function CardHeading(props: { Icon: IconComponent; title: string; subtitle?: string; action?: ReactNode }) {
  return (
    <header className="card-heading">
      <props.Icon className="card-icon" />
      <span><h2>{props.title}</h2>{props.subtitle ? <p>{props.subtitle}</p> : null}</span>
      {props.action ? <button type="button" className="card-action">{props.action}</button> : null}
    </header>
  );
}

function EmptyState({ children }: { children: ReactNode }) {
  return <div className="empty-state">{children}</div>;
}
