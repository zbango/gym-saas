import { useState, type FormEvent } from "react";
import { PageHeader, SettingsIcon, Tabs, type TabItem } from "@gym-saas/ui";

type SettingsTab = "profile" | "security" | "notifications" | "device" | "surveys" | "system";

const settingsTabs = [
  { id: "profile", label: "Perfil" },
  { id: "security", label: "Seguridad" },
  { id: "notifications", label: "Notificaciones" },
  { id: "device", label: "Configuración del Dispositivo" },
  { id: "surveys", label: "Encuestas" },
  { id: "system", label: "Sistema" }
] as const satisfies readonly TabItem<SettingsTab>[];

type ProfileForm = {
  name: string;
  email: string;
  phone: string;
};

const initialProfile: ProfileForm = {
  name: "Super Administrator",
  email: "steven.tabango@gmail.com",
  phone: ""
};

export function SettingsPage() {
  const [activeTab, setActiveTab] = useState<SettingsTab>("profile");
  const [profile, setProfile] = useState(initialProfile);
  const [saved, setSaved] = useState(false);

  function saveProfile(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSaved(true);
  }

  return (
    <div className="px-16 py-14 text-[#1c2431] max-[1260px]:px-[43px] max-[1260px]:py-[46px] max-[920px]:px-[22px] max-[920px]:py-8 max-[560px]:px-[14px] max-[560px]:py-6">
      <PageHeader icon={<SettingsIcon />} title="Configuración del Sistema" description="Gestiona tu perfil, gimnasio y preferencias" className="mb-[43px]" />
      <div className="overflow-x-auto">
        <Tabs tabs={settingsTabs} activeTab={activeTab} onChange={setActiveTab} ariaLabel="Secciones de configuración" className="gap-[57px]" />
      </div>

      <section id={`tabpanel-${activeTab}`} role="tabpanel" aria-labelledby={`tab-${activeTab}`} className="mt-10 min-h-[705px] rounded-[14px] border border-[#e1e4e8] bg-white px-10 py-11 shadow-[0_2px_4px_rgba(23,31,44,0.025)] max-[920px]:min-h-0 max-[560px]:mt-6 max-[560px]:px-5 max-[560px]:py-7">
        {activeTab === "profile" ? (
          <form className="max-w-[1075px]" onSubmit={saveProfile}>
            <h2 className="m-0 text-[30px] font-extrabold tracking-[-0.04em] text-[#172131] max-[560px]:text-[25px]">Información de Perfil</h2>
            <div className="mt-11 grid gap-7">
              <SettingsField label="Nombre Completo" htmlFor="profile-name">
                <input id="profile-name" required value={profile.name} onChange={(event) => { setProfile((current) => ({ ...current, name: event.target.value })); setSaved(false); }} className={fieldClassName} />
              </SettingsField>
              <SettingsField label="Correo Electrónico" htmlFor="profile-email">
                <input id="profile-email" type="email" required value={profile.email} onChange={(event) => { setProfile((current) => ({ ...current, email: event.target.value })); setSaved(false); }} className={fieldClassName} />
              </SettingsField>
              <SettingsField label="Teléfono" htmlFor="profile-phone">
                <input id="profile-phone" type="tel" value={profile.phone} onChange={(event) => { setProfile((current) => ({ ...current, phone: event.target.value })); setSaved(false); }} placeholder="Ingrese su número de teléfono" className={fieldClassName} />
              </SettingsField>
            </div>
            <button type="submit" className="mt-12 min-h-[70px] rounded-xl border-0 bg-[linear-gradient(100deg,var(--color-brand-button-start)_0%,var(--color-brand-button-middle)_42%,var(--color-brand-button-end)_100%)] px-8 text-[23px] font-extrabold text-white shadow-[0_12px_22px_rgba(178,103,33,0.18)] max-[560px]:min-h-[57px] max-[560px]:w-full max-[560px]:text-[19px]">Guardar Cambios</button>
            {saved ? <p className="mt-4 text-[17px] font-semibold text-[var(--color-brand-success)]" role="status">Cambios guardados para esta sesión.</p> : null}
          </form>
        ) : <SettingsPlaceholder tab={settingsTabs.find((tab) => tab.id === activeTab)!} />}
      </section>
    </div>
  );
}

function SettingsField({ label, htmlFor, children }: { label: string; htmlFor: string; children: React.ReactNode }) {
  return <label className="grid gap-[13px] text-[22px] font-bold text-[#394351] max-[560px]:text-[18px]" htmlFor={htmlFor}><span>{label}</span>{children}</label>;
}

function SettingsPlaceholder({ tab }: { tab: TabItem<SettingsTab> }) {
  return <div className="grid min-h-[490px] place-items-center text-center max-[920px]:min-h-[280px]"><div><h2 className="m-0 text-[30px] font-extrabold tracking-[-0.04em] text-[#172131]">{tab.label}</h2><p className="mt-4 text-[19px] text-[#68727e]">La configuración de esta sección estará disponible próximamente.</p></div></div>;
}

const fieldClassName = "min-h-[79px] w-full rounded-xl border border-[#cbd0d7] bg-white px-7 text-[22px] text-[#172131] outline-0 placeholder:text-[#757d88] focus:border-[var(--color-brand-accent)] focus:outline-3 focus:outline-[color-mix(in_srgb,var(--color-brand-accent)_32%,transparent)] max-[560px]:min-h-[58px] max-[560px]:px-5 max-[560px]:text-[18px]";
