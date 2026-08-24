import { useEffect, useState, type FormEvent, type ReactNode } from "react";
import { ChevronDownIcon, Dialog, UserIcon } from "@gym-saas/ui";

export type CreatableUserRole = "gym_admin" | "receptionist";

export type NewUserInput = {
  name: string;
  username: string;
  email: string;
  phone: string;
  role: CreatableUserRole;
};

type FormState = NewUserInput;

const initialForm: FormState = {
  name: "",
  username: "",
  email: "",
  phone: "",
  role: "receptionist"
};

export function UserCreateDialog({ open, onClose, onCreated }: {
  open: boolean;
  onClose: () => void;
  onCreated: (user: NewUserInput) => void;
}) {
  const [form, setForm] = useState<FormState>(initialForm);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (open) {
      setForm(initialForm);
      setError(null);
    }
  }, [open]);

  function update<Key extends keyof FormState>(key: Key, value: FormState[Key]) {
    setForm((current) => ({ ...current, [key]: value }));
  }

  function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const username = form.username.trim().toLowerCase();
    if (!/^[a-z0-9._-]+$/.test(username)) {
      setError("El nombre de usuario solo puede incluir letras, números, puntos, guiones y guiones bajos.");
      return;
    }

    onCreated({
      ...form,
      name: form.name.trim(),
      username,
      email: form.email.trim(),
      phone: form.phone.trim()
    });
    onClose();
  }

  return (
    <Dialog open={open} title="Agregar Nuevo Usuario" onClose={onClose} className="min-h-[calc(100dvh-24px)] max-w-[1344px] [zoom:1.333333] [&>header]:min-h-[154px]">
      <form className="flex min-h-[calc(100dvh-178px)] flex-col gap-10 px-11 pb-12 pt-12 max-[560px]:min-h-0 max-[560px]:gap-5 max-[560px]:px-[23px] max-[560px]:py-7" onSubmit={submit}>
        <div className="grid grid-cols-2 gap-x-8 gap-y-7 max-[720px]:grid-cols-1">
          <FormField label="Nombre Completo" htmlFor="user-name">
            <input id="user-name" required value={form.name} onChange={(event) => update("name", event.target.value)} placeholder="Ej: Juan Pérez" autoComplete="name" className={fieldClassName} />
          </FormField>
          <FormField label="Nombre de Usuario" htmlFor="user-username">
            <input id="user-username" required value={form.username} onChange={(event) => update("username", event.target.value)} placeholder="Ej: juanperez" autoComplete="username" className={fieldClassName} />
          </FormField>
          <FormField label="Dirección de Email" htmlFor="user-email">
            <input id="user-email" required type="email" value={form.email} onChange={(event) => update("email", event.target.value)} placeholder="ejemplo@correo.com" autoComplete="email" className={fieldClassName} />
          </FormField>
          <FormField label="Número de Teléfono" htmlFor="user-phone">
            <input id="user-phone" required type="tel" value={form.phone} onChange={(event) => update("phone", event.target.value)} placeholder="Ej: 0999999999" autoComplete="tel" className={fieldClassName} />
          </FormField>
          <FormField label="Rol" htmlFor="user-role" className="col-span-2 max-[720px]:col-span-1">
            <span className="relative block">
              <select id="user-role" value={form.role} onChange={(event) => update("role", event.target.value as CreatableUserRole)} className={`${fieldClassName} appearance-none pr-16`}>
                <option value="receptionist">Recepcionista</option>
                <option value="gym_admin">Administrador</option>
              </select>
              <ChevronDownIcon className="pointer-events-none absolute right-5 top-1/2 h-7 w-7 -translate-y-1/2 text-[#a4acb6]" />
            </span>
          </FormField>
        </div>

        <p className="m-0 rounded-[15px] border border-[#a8edf3] bg-[#effdff] px-8 py-7 text-[19px] leading-[1.45] text-[#315e71] max-[560px]:px-5 max-[560px]:py-5 max-[560px]:text-[16px]">
          <strong>Nota:</strong> El usuario recibirá una contraseña temporal por defecto. Podrá cambiarla después de iniciar sesión.
        </p>
        {error ? <p className="m-0 text-[17px] font-bold text-[var(--color-brand-danger)]" role="alert">{error}</p> : null}

        <footer className="mt-auto flex flex-wrap justify-end gap-6 pt-2 max-[560px]:grid max-[560px]:grid-cols-1 max-[560px]:gap-3">
          <button type="button" className="min-h-[71px] rounded-xl border-0 bg-[#374357] px-10 text-[23px] font-extrabold text-white shadow-[0_12px_22px_rgba(33,44,62,0.18)] max-[560px]:min-h-[59px] max-[560px]:text-[19px]" onClick={onClose}>Cancelar</button>
          <button type="submit" className="inline-flex min-h-[71px] items-center justify-center gap-4 rounded-xl border-0 bg-[linear-gradient(100deg,var(--color-brand-button-start)_0%,var(--color-brand-button-middle)_42%,var(--color-brand-button-end)_100%)] px-10 text-[23px] font-extrabold text-[var(--color-brand-button-text)] shadow-[0_12px_22px_rgba(178,103,33,0.18)] max-[560px]:min-h-[59px] max-[560px]:text-[19px]"><UserIcon className="h-7 w-7" /> Crear Usuario</button>
        </footer>
      </form>
    </Dialog>
  );
}

function FormField({ label, htmlFor, className, children }: { label: string; htmlFor: string; className?: string; children: ReactNode }) {
  return <label className={`grid gap-3 text-[21px] font-bold text-[#394351] max-[560px]:text-[18px] ${className ?? ""}`.trim()} htmlFor={htmlFor}><span>{label}</span>{children}</label>;
}

const fieldClassName = "min-h-[100px] w-full rounded-xl border border-[#cbd0d7] bg-white px-8 text-[22px] text-[#172131] outline-0 placeholder:text-[#707782] focus:border-[var(--color-brand-accent)] focus:outline-3 focus:outline-[color-mix(in_srgb,var(--color-brand-accent)_32%,transparent)] max-[560px]:min-h-[58px] max-[560px]:px-5 max-[560px]:text-[18px]";
