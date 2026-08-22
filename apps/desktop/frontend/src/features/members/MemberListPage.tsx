import { useState } from "react";
import {
  AlertCircleIcon,
  AlertTriangleIcon,
  CalendarIcon,
  CheckCircleIcon,
  DataTable,
  type DataTableColumn,
  MetricCard,
  MoreVerticalIcon,
  PageHeader,
  PlusIcon,
  RefreshIcon,
  SearchIcon,
  UsersIcon
} from "@gym-saas/ui";
import "./member-list-page.css";

type MemberListRow = {
  id: string;
  name: string;
  age: number;
  identification: string;
  biometricID?: string;
  email: string;
  phone: string;
  membership: string;
  expiry: { kind: "visits"; value: string } | { kind: "date"; value: string };
};

const members: MemberListRow[] = [
  { id: "mock-1", name: "Alexandra Vega", age: 31, identification: "1003751680", biometricID: "906", email: "alexandra.vega@example.test", phone: "099 464 6625", membership: "Plan Tarjeta", expiry: { kind: "visits", value: "24 / 25 visitas" } },
  { id: "mock-2", name: "Sebastián Paredes", age: 15, identification: "1003000898", email: "sebastian.paredes@example.test", phone: "098 387 2342", membership: "Plan Mensual MMA", expiry: { kind: "date", value: "5 mar 2026" } },
  { id: "mock-3", name: "Alex Granda", age: 56, identification: "1001670999", biometricID: "893", email: "alex.granda@example.test", phone: "099 665 0191", membership: "Plan Tarjeta", expiry: { kind: "visits", value: "24 / 25 visitas" } },
  { id: "mock-4", name: "Soraida Méndez", age: 52, identification: "0400875936", biometricID: "892", email: "soraida.mendez@example.test", phone: "099 134 3480", membership: "Plan Tarjeta", expiry: { kind: "visits", value: "24 / 25 visitas" } },
  { id: "mock-5", name: "Jordan Loor", age: 22, identification: "1005231814", biometricID: "890", email: "jordan.loor@example.test", phone: "095 973 3250", membership: "Plan Quincenal", expiry: { kind: "date", value: "16 feb 2026" } },
  { id: "mock-6", name: "Santiago Paucar", age: 22, identification: "1050436656", biometricID: "889", email: "santiago.paucar@example.test", phone: "099 526 3290", membership: "Plan Tarjeta", expiry: { kind: "visits", value: "25 visitas" } },
  { id: "mock-7", name: "Heydan Cheza", age: 17, identification: "1050396926", biometricID: "885", email: "heydan.cheza@example.test", phone: "095 961 9188", membership: "Plan Mensual", expiry: { kind: "date", value: "4 mar 2026" } },
  { id: "mock-8", name: "Mateo Fuelagan", age: 20, identification: "0402127765", biometricID: "876", email: "mateo.fuelagan@example.test", phone: "093 924 4357", membership: "Plan Mensual", expiry: { kind: "date", value: "4 mar 2026" } },
  { id: "mock-9", name: "Anthony Valverde", age: 30, identification: "1004131916", biometricID: "875", email: "anthony.valverde@example.test", phone: "096 759 1457", membership: "Plan Mensual", expiry: { kind: "date", value: "4 mar 2026" } },
  { id: "mock-10", name: "Anabela Tates", age: 24, identification: "1005386543", biometricID: "859", email: "anabela.tates@example.test", phone: "096 919 6081", membership: "Plan Tarjeta", expiry: { kind: "visits", value: "23 / 25 visitas" } }
];

const columns: Array<DataTableColumn<MemberListRow>> = [
  {
    id: "actions",
    header: <span className="visually-hidden">Acciones</span>,
    cell: () => <button className="row-action" type="button" aria-label="Acciones del miembro"><MoreVerticalIcon /></button>,
    width: "55px",
    align: "center"
  },
  {
    id: "member",
    header: <button className="table-sort" type="button">Cliente <span aria-hidden="true">↕</span></button>,
    cell: (row) => <span className="member-name"><strong>{row.name}</strong><small>{row.age} años</small></span>,
    width: "235px"
  },
  {
    id: "identification",
    header: <button className="table-sort" type="button">Identificación <span aria-hidden="true">↕</span></button>,
    cell: (row) => row.identification,
    width: "185px"
  },
  {
    id: "biometric",
    header: "Biométrico ID",
    cell: (row) => <span className="muted-cell">{row.biometricID ?? ""}</span>,
    width: "150px"
  },
  {
    id: "contact",
    header: "Contacto",
    cell: (row) => <span className="contact-cell"><strong>{row.email}</strong><small>{row.phone}</small></span>,
    width: "325px"
  },
  { id: "membership", header: "Membresía", cell: (row) => row.membership, width: "210px" },
  {
    id: "expiry",
    header: "Expiración",
    cell: (row) => row.expiry.kind === "visits"
      ? <span className="expiry expiry-visits"><span aria-hidden="true">▥</span>{row.expiry.value}</span>
      : <span className="expiry expiry-date"><CalendarIcon />{row.expiry.value}</span>,
    width: "185px"
  }
];

export function MemberListPage() {
  const [query, setQuery] = useState("");
  const [status, setStatus] = useState("all");
  const [page, setPage] = useState(1);

  return (
    <div className="member-list-page">
      <PageHeader
        className="member-page-header"
        icon={<UsersIcon />}
        title="Lista de Clientes"
        description="Ver y gestionar todos los clientes del gimnasio"
        actions={
          <>
            <button className="member-secondary-action" type="button"><RefreshIcon /> Actualizar</button>
            <button className="member-primary-action" type="button"><PlusIcon /> Agregar Nuevo Cliente</button>
          </>
        }
      />

      <section className="member-metrics" aria-label="Resumen de miembros">
        <MetricCard value="291" label="Total General de usuarios registrados" icon={<UsersIcon />} selected />
        <MetricCard value="76" label="Total de Clientes con membresía activa pagada por completo" icon={<CheckCircleIcon />} tone="success" />
        <MetricCard value="16" label="Total de Clientes con membresía activa pendiente de pago" icon={<AlertTriangleIcon />} tone="warning" />
        <MetricCard value="197" label="Total de Clientes con membresía expirada" icon={<AlertCircleIcon />} tone="danger" />
        <MetricCard value="6" label="Total de Clientes con membresía expirada y pendiente de pago" icon={<AlertCircleIcon />} tone="danger" />
      </section>

      <section className="member-controls" aria-label="Controles de miembros">
        <label className="member-search"><SearchIcon /><input value={query} onChange={(event) => setQuery(event.target.value)} placeholder="Buscar usuarios..." /></label>
        <label className="member-status"><span className="visually-hidden">Estado</span><select value={status} onChange={(event) => setStatus(event.target.value)}><option value="all">Todos los Estados</option><option value="active">Activos</option><option value="pending">Pendientes</option><option value="expired">Expirados</option></select></label>
      </section>

      <DataTable className="member-table" caption="Listado simulado de miembros" columns={columns} rows={members} rowKey={(row) => row.id} />

      <footer className="member-pagination">
        <label>Mostrar <select defaultValue="10"><option value="10">10</option><option value="25">25</option><option value="50">50</option></select> entradas</label>
        <p>Mostrando 1 a 10 de 291 entradas</p>
        <nav aria-label="Paginación de miembros">
          <button type="button" disabled>««</button><button type="button" disabled>«</button>
          <button type="button" className={page === 1 ? "is-current" : ""} onClick={() => setPage(1)}>1</button>
          <button type="button" className={page === 2 ? "is-current" : ""} onClick={() => setPage(2)}>2</button>
          <span>…</span><button type="button" onClick={() => setPage(30)}>30</button><button type="button" onClick={() => setPage(Math.min(30, page + 1))}>»</button><button type="button" onClick={() => setPage(30)}>»»</button>
        </nav>
      </footer>
    </div>
  );
}
