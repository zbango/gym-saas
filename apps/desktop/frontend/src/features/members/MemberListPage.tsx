import { useState, type ReactNode } from "react";
import {
  AlertCircleIcon,
  AlertTriangleIcon,
  CalendarIcon,
  CheckCircleIcon,
  MoreVerticalIcon,
  PlusIcon,
  RefreshIcon,
  SearchIcon,
  UsersIcon
} from "@gym-saas/ui";

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

type MemberMetric = {
  value: string;
  label: string;
  icon: ReactNode;
  tone?: "default" | "success" | "warning" | "danger";
  selected?: boolean;
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

const metrics: MemberMetric[] = [
  { value: "291", label: "Total general de miembros registrados", icon: <UsersIcon />, selected: true },
  { value: "76", label: "Miembros con membresía activa pagada por completo", icon: <CheckCircleIcon />, tone: "success" },
  { value: "16", label: "Miembros con membresía activa pendiente de pago", icon: <AlertTriangleIcon />, tone: "warning" },
  { value: "197", label: "Miembros con membresía expirada", icon: <AlertCircleIcon />, tone: "danger" },
  { value: "6", label: "Miembros con membresía expirada y pendiente de pago", icon: <AlertCircleIcon />, tone: "danger" }
];

const tableHeaderClassName = "h-[78px] border-b border-[#e1e4e8] px-[22px] text-left text-base font-extrabold uppercase tracking-[0.08em] text-[#737d89]";
const cellClassName = "h-[121px] border-b border-[#e1e4e8] px-[22px] py-[17px] text-lg text-[#1c2431]";
const pageButtonClassName = "flex h-11 min-w-11 items-center justify-center rounded-[10px] border border-[#f0df9e] bg-white px-3 text-[17px] font-bold text-[#c7cbd1] transition-colors disabled:cursor-not-allowed";

export function MemberListPage() {
  const [query, setQuery] = useState("");
  const [status, setStatus] = useState("all");
  const [page, setPage] = useState(1);

  return (
    <div className="px-16 py-14 text-[#1c2431] max-[1260px]:px-[43px] max-[1260px]:py-[46px] max-[920px]:px-[22px] max-[920px]:py-8 max-[560px]:px-[14px] max-[560px]:py-6">
      <section className="mb-[46px] flex items-start justify-between gap-7 max-[920px]:grid">
        <div className="flex items-start gap-[18px]">
          <div className="grid place-items-center pt-1 text-[#56616d] [&>svg]:h-[37px] [&>svg]:w-[37px]">
            <UsersIcon />
          </div>
          <div>
            <h1 className="m-0 text-[38px] leading-[1.15] font-extrabold tracking-[-0.04em] max-[1260px]:text-[33px] max-[560px]:text-[28px]">
              Lista de Miembros
            </h1>
            <p className="mt-2.5 text-[21px] text-[#5f6974] max-[560px]:text-[17px]">
              Ver y gestionar todos los miembros del gimnasio
            </p>
          </div>
        </div>
        <div className="flex gap-[14px] max-[920px]:flex-wrap">
          <button
            className="inline-flex min-h-[58px] items-center gap-3 rounded-xl border-0 bg-[#374357] px-7 text-xl font-extrabold whitespace-nowrap text-white shadow-[0_12px_22px_rgba(33,44,62,0.18)] max-[1260px]:min-h-[49px] max-[1260px]:px-[18px] max-[1260px]:text-[17px] max-[560px]:w-full max-[560px]:justify-center"
            type="button"
          >
            <RefreshIcon className="h-6 w-6" /> Actualizar
          </button>
          <button
            className="inline-flex min-h-[58px] items-center gap-3 rounded-xl border-0 bg-[linear-gradient(100deg,var(--color-brand-button-start)_0%,var(--color-brand-button-middle)_42%,var(--color-brand-button-end)_100%)] px-7 text-xl font-extrabold whitespace-nowrap text-[var(--color-brand-button-text)] shadow-[0_12px_22px_rgba(178,103,33,0.18)] max-[1260px]:min-h-[49px] max-[1260px]:px-[18px] max-[1260px]:text-[17px] max-[560px]:w-full max-[560px]:justify-center"
            type="button"
          >
            <PlusIcon className="h-6 w-6" /> Agregar Nuevo Miembro
          </button>
        </div>
      </section>

      <section
        className="mb-[39px] grid grid-cols-5 gap-6 max-[1260px]:grid-cols-3 max-[920px]:grid-cols-2 max-[560px]:grid-cols-1"
        aria-label="Resumen de miembros"
      >
        {metrics.map((metric) => (
          <MetricSummaryCard key={metric.label} metric={metric} />
        ))}
      </section>

      <section className="mb-[26px] flex items-center gap-[30px] max-[920px]:flex-col max-[920px]:items-stretch max-[920px]:gap-[14px]" aria-label="Controles de miembros">
        <label className="flex h-[82px] w-[520px] items-center gap-4 rounded-xl border border-[#d5d9df] bg-white px-[19px] text-[#9da6af] max-[1260px]:w-[440px] max-[920px]:w-full">
          <SearchIcon className="h-[29px] w-[29px]" />
          <input
            className="w-full border-0 bg-transparent text-xl text-[#1c2431] outline-0 placeholder:text-[#7a8490]"
            value={query}
            onChange={(event) => setQuery(event.target.value)}
            placeholder="Buscar miembros..."
          />
        </label>
        <label className="block">
          <span className="sr-only">Estado</span>
          <select
            className="h-[62px] w-[300px] rounded-[11px] border border-[#d5d9df] bg-white px-[22px] text-[19px] text-[#1c2431] outline-0 max-[920px]:w-full"
            value={status}
            onChange={(event) => setStatus(event.target.value)}
          >
            <option value="all">Todos los Estados</option>
            <option value="active">Activos</option>
            <option value="pending">Pendientes</option>
            <option value="expired">Expirados</option>
          </select>
        </label>
      </section>

      <div className="overflow-hidden rounded-t-[11px] border border-[#e1e4e8] bg-white">
        <div className="overflow-x-auto">
          <table className="w-full min-w-[1290px] border-collapse">
            <caption className="sr-only">Listado simulado de miembros</caption>
            <thead>
              <tr>
                <th className={`${tableHeaderClassName} w-[55px] text-center`}>
                  <span className="sr-only">Acciones</span>
                </th>
                <th className={`${tableHeaderClassName} w-[235px]`}>
                  <SortHeader label="Miembro" />
                </th>
                <th className={`${tableHeaderClassName} w-[185px]`}>
                  <SortHeader label="Identificación" />
                </th>
                <th className={`${tableHeaderClassName} w-[150px]`}>Biométrico ID</th>
                <th className={`${tableHeaderClassName} w-[325px]`}>Contacto</th>
                <th className={`${tableHeaderClassName} w-[210px]`}>Membresía</th>
                <th className={`${tableHeaderClassName} w-[185px]`}>Expiración</th>
              </tr>
            </thead>
            <tbody>
              {members.map((member) => (
                <tr key={member.id}>
                  <td className={`${cellClassName} text-center`}>
                    <button
                      className="border-0 bg-transparent p-0 text-[#7d8691] [&>svg]:h-[27px] [&>svg]:w-[27px]"
                      type="button"
                      aria-label="Acciones del miembro"
                    >
                      <MoreVerticalIcon />
                    </button>
                  </td>
                  <td className={cellClassName}>
                    <span className="grid gap-[5px]">
                      <strong className="text-xl">{member.name}</strong>
                      <small className="text-lg text-[#7c8590]">{member.age} años</small>
                    </span>
                  </td>
                  <td className={cellClassName}>{member.identification}</td>
                  <td className={`${cellClassName} text-[#7c8590]`}>{member.biometricID ?? ""}</td>
                  <td className={cellClassName}>
                    <span className="grid gap-[5px]">
                      <strong className="text-xl">{member.email}</strong>
                      <small className="text-lg text-[#7c8590]">{member.phone}</small>
                    </span>
                  </td>
                  <td className={cellClassName}>{member.membership}</td>
                  <td className={cellClassName}>
                    {member.expiry.kind === "visits" ? (
                      <span className="inline-flex items-center gap-[9px] text-lg font-bold whitespace-nowrap text-[#4b78dd]">
                        <span aria-hidden="true" className="text-[25px]">▥</span>
                        {member.expiry.value}
                      </span>
                    ) : (
                      <span className="inline-flex items-center gap-[9px] text-lg font-bold whitespace-nowrap text-[var(--color-brand-danger)] [&>svg]:h-[22px] [&>svg]:w-[22px]">
                        <CalendarIcon />
                        {member.expiry.value}
                      </span>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      <footer className="flex min-h-[92px] items-center justify-between gap-6 rounded-b-[11px] border border-t-0 border-[#e1e4e8] bg-white px-[31px] text-[17px] text-[#4d5865] max-[920px]:flex-col max-[920px]:items-start max-[920px]:p-[17px]">
        <label className="inline-flex items-center gap-[10px]">
          Mostrar
          <select className="h-[45px] rounded-[10px] border border-[#d3d7dd] bg-white px-[13px] pr-[29px]" defaultValue="10">
            <option value="10">10</option>
            <option value="25">25</option>
            <option value="50">50</option>
          </select>
          entradas
        </label>
        <p className="m-0 max-[560px]:text-sm">Mostrando 1 a 10 de 291 entradas</p>
        <nav className="flex items-center gap-2 max-[920px]:flex-wrap" aria-label="Paginación de miembros">
          <button className={pageButtonClassName} type="button" disabled>««</button>
          <button className={pageButtonClassName} type="button" disabled>«</button>
          <button
            className={`${pageButtonClassName} ${page === 1 ? "border-[#354155] bg-[#354155] text-white" : ""}`}
            type="button"
            onClick={() => setPage(1)}
          >
            1
          </button>
          <button
            className={`${pageButtonClassName} ${page === 2 ? "border-[#354155] bg-[#354155] text-white" : ""}`}
            type="button"
            onClick={() => setPage(2)}
          >
            2
          </button>
          <span className="px-[6px]">…</span>
          <button className={pageButtonClassName} type="button" onClick={() => setPage(30)}>30</button>
          <button className={pageButtonClassName} type="button" onClick={() => setPage(Math.min(30, page + 1))}>»</button>
          <button className={pageButtonClassName} type="button" onClick={() => setPage(30)}>»»</button>
        </nav>
      </footer>
    </div>
  );
}

function MetricSummaryCard({ metric }: { metric: MemberMetric }) {
  const iconToneClasses = {
    default: "bg-[#f7f1ff] text-[#954fee]",
    success: "bg-[var(--color-brand-success-soft)] text-[var(--color-brand-success)]",
    warning: "bg-[#fffbea] text-[#c59327]",
    danger: "bg-[var(--color-brand-danger-soft)] text-[var(--color-brand-danger)]"
  } satisfies Record<NonNullable<MemberMetric["tone"]> | "default", string>;

  return (
    <article
      className={`flex min-h-[230px] items-center gap-[23px] rounded-[15px] border bg-white px-7 py-[31px] max-[1260px]:min-h-[180px] max-[560px]:min-h-[135px] ${
        metric.selected
          ? "border-[3px] border-[#a46df2] px-[26px] py-[29px] shadow-[0_0_0_4px_rgba(164,109,242,0.12)]"
          : "border-[#e0e3e8]"
      }`}
    >
      <div className={`grid h-[76px] w-[76px] place-items-center rounded-2xl ${iconToneClasses[metric.tone ?? "default"]} [&>svg]:h-[39px] [&>svg]:w-[39px]`}>
        {metric.icon}
      </div>
      <div>
        <strong className="block text-[40px] leading-none">{metric.value}</strong>
        <p className="mt-2.5 text-[20px] leading-[1.45] text-[#59636e]">{metric.label}</p>
      </div>
    </article>
  );
}

function SortHeader({ label }: { label: string }) {
  return (
    <button className="border-0 bg-transparent font-inherit font-extrabold uppercase tracking-[0.08em] text-inherit" type="button">
      {label} <span aria-hidden="true" className="text-[20px] text-[#6f94b4]">↕</span>
    </button>
  );
}
