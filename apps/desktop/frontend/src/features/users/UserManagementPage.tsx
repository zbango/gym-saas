import { useState } from "react";
import {
  CheckCircleIcon,
  DataTable,
  type DataTableColumn,
  MetricCard,
  PageHeader,
  PlusIcon,
  RefreshIcon,
  SearchIcon,
  ShieldIcon,
  UserIcon,
  UsersIcon
} from "@gym-saas/ui";

type UserRole = "super_admin" | "gym_owner" | "gym_admin" | "trainer" | "receptionist" | "member";

type SystemUser = {
  id: string;
  name: string;
  handle: string;
  email: string;
  role: UserRole;
  active: boolean;
  lastAccess: string;
  createdAt: string;
};

const users: SystemUser[] = [
  { id: "mock-user-1", name: "Carla Mora", handle: "@carla", email: "carla.mora@example.test", role: "receptionist", active: true, lastAccess: "Nunca", createdAt: "6 ene 2026, 13:44" },
  { id: "mock-user-2", name: "Diego Vera", handle: "@diego", email: "diego.vera@example.test", role: "receptionist", active: true, lastAccess: "Nunca", createdAt: "6 ene 2026, 13:23" },
  { id: "mock-user-3", name: "Sofía Paz", handle: "@sofia", email: "sofia.paz@example.test", role: "receptionist", active: true, lastAccess: "Nunca", createdAt: "6 ene 2026, 13:05" },
  { id: "mock-user-4", name: "Mateo León", handle: "@mateo", email: "mateo.leon@example.test", role: "gym_admin", active: true, lastAccess: "Nunca", createdAt: "30 dic 2025, 21:11" },
  { id: "mock-user-5", name: "Paula Navas", handle: "@paula", email: "paula.navas@example.test", role: "gym_admin", active: true, lastAccess: "Nunca", createdAt: "30 dic 2025, 21:11" }
];

const roleLabels: Record<UserRole, string> = {
  super_admin: "Super administrador",
  gym_owner: "Propietario",
  gym_admin: "Administrador",
  trainer: "Entrenador",
  receptionist: "Recepcionista",
  member: "Miembro"
};

const tableHeaderClass = "h-[78px] border-b border-[#e1e4e8] px-[28px] text-left text-[16px] font-extrabold uppercase tracking-[0.08em] text-[#737d89]";
const tableCellClass = "h-[135px] border-b border-[#e1e4e8] px-[28px] py-[20px] text-[19px] text-[#1c2431]";

const columns: Array<DataTableColumn<SystemUser>> = [
  { id: "user", header: <SortHeader label="Usuario" />, headerClassName: `${tableHeaderClass} w-[38%]`, cellClassName: tableCellClass, cell: (user) => <span className="grid gap-[5px]"><strong className="text-[25px]">{user.name}</strong><small className="text-[20px] text-[#7c8590]">{user.handle} · {user.email}</small></span> },
  { id: "role", header: <SortHeader label="Rol" />, headerClassName: `${tableHeaderClass} w-[18%]`, cellClassName: tableCellClass, cell: (user) => <RoleBadge role={user.role} /> },
  { id: "status", header: <SortHeader label="Estado" />, headerClassName: `${tableHeaderClass} w-[14%]`, cellClassName: tableCellClass, cell: (user) => <span className="inline-flex rounded-full bg-[var(--color-brand-success-soft)] px-[18px] py-[7px] text-[18px] font-bold text-[#427a4a]">{user.active ? "Activo" : "Inactivo"}</span> },
  { id: "last-access", header: <SortHeader label="Último acceso" />, headerClassName: `${tableHeaderClass} w-[15%]`, cellClassName: tableCellClass, cell: (user) => user.lastAccess },
  { id: "created", header: <SortHeader label="Creado" />, headerClassName: `${tableHeaderClass} w-[20%]`, cellClassName: tableCellClass, cell: (user) => user.createdAt }
];

export function UserManagementPage() {
  const [query, setQuery] = useState("");

  return (
    <div className="px-16 py-14 text-[#1c2431] max-[1260px]:px-[43px] max-[1260px]:py-[46px] max-[920px]:px-[22px] max-[920px]:py-8 max-[560px]:px-[14px] max-[560px]:py-6">
      <PageHeader
        icon={<UsersIcon />}
        title="Gestión de Usuarios"
        description="Gestiona cuentas de usuarios administradores y personal"
        actions={<><button className="inline-flex min-h-[58px] items-center gap-3 rounded-xl border-0 bg-[#374357] px-7 text-xl font-extrabold whitespace-nowrap text-white shadow-[0_12px_22px_rgba(33,44,62,0.18)] max-[1260px]:min-h-[49px] max-[1260px]:px-[18px] max-[1260px]:text-[17px] max-[560px]:w-full max-[560px]:justify-center" type="button"><RefreshIcon className="h-6 w-6" /> Actualizar</button><button className="inline-flex min-h-[58px] items-center gap-3 rounded-xl border-0 bg-[linear-gradient(100deg,var(--color-brand-button-start)_0%,var(--color-brand-button-middle)_42%,var(--color-brand-button-end)_100%)] px-7 text-xl font-extrabold whitespace-nowrap text-[var(--color-brand-button-text)] shadow-[0_12px_22px_rgba(178,103,33,0.18)] max-[1260px]:min-h-[49px] max-[1260px]:px-[18px] max-[1260px]:text-[17px] max-[560px]:w-full max-[560px]:justify-center" type="button"><PlusIcon className="h-6 w-6" /> Agregar Usuario</button></>}
      />

      <section className="mb-[39px] grid grid-cols-4 gap-7 max-[920px]:grid-cols-2 max-[560px]:grid-cols-1" aria-label="Resumen de usuarios">
        <MetricCard value="5" label="Total de Usuarios" icon={<UsersIcon />} />
        <MetricCard value="2" label="Administradores" icon={<ShieldIcon />} tone="danger" />
        <MetricCard value="0" label="Miembros del Personal" icon={<UserIcon />} tone="info" />
        <MetricCard value="5" label="Usuarios Activos" icon={<CheckCircleIcon />} tone="success" />
      </section>

      <label className="mb-[29px] flex h-[82px] w-[620px] items-center gap-4 rounded-xl border border-[#d5d9df] bg-white px-[19px] text-[#9da6af] max-[920px]:w-full">
        <SearchIcon className="h-[29px] w-[29px]" />
        <input className="w-full border-0 bg-transparent text-xl text-[#1c2431] outline-0 placeholder:text-[#7a8490]" value={query} onChange={(event) => setQuery(event.target.value)} placeholder="Buscar usuarios..." />
      </label>

      <DataTable caption="Listado simulado de usuarios del sistema" columns={columns} rows={users} rowKey={(user) => user.id} tableClassName="min-w-[1230px]" />

      <footer className="flex min-h-[92px] flex-wrap items-center justify-between gap-6 rounded-b-[11px] border border-t-0 border-[#e1e4e8] bg-white px-[31px] py-[17px] text-[17px] text-[#4d5865]">
        <label className="inline-flex items-center gap-[10px]">Mostrar <select className="h-[45px] rounded-[10px] border border-[#d3d7dd] bg-white px-[13px] pr-[29px] font-inherit" defaultValue="10"><option value="10">10</option></select> entradas</label>
        <p className="m-0">Mostrando 1 a 5 de 5 entradas</p>
        <nav className="flex items-center gap-2" aria-label="Paginación de usuarios"><PaginationButton disabled>««</PaginationButton><PaginationButton disabled>«</PaginationButton><PaginationButton current>1</PaginationButton><PaginationButton disabled>»</PaginationButton><PaginationButton disabled>»»</PaginationButton></nav>
      </footer>
    </div>
  );
}

function RoleBadge({ role }: { role: UserRole }) { return <span className={`inline-flex rounded-full px-[18px] py-[7px] text-[18px] font-bold ${role === "receptionist" ? "bg-[#e5ecff] text-[#4642a3]" : "bg-[#fff8c6] text-[#86632d]"}`}>{roleLabels[role]}</span>; }

function SortHeader({ label }: { label: string }) { return <button className="border-0 bg-transparent font-inherit font-extrabold uppercase tracking-inherit text-inherit" type="button">{label} <span className="text-[#6f94b4]">↕</span></button>; }
function PaginationButton({ children, disabled = false, current = false }: { children: string; disabled?: boolean; current?: boolean }) { return <button className={`flex h-11 min-w-11 items-center justify-center rounded-[10px] border border-[#f0df9e] px-3 font-bold ${current ? "border-[#354155] bg-[#354155] text-white" : "bg-white text-[#c7cbd1]"}`} type="button" disabled={disabled}>{children}</button>; }
