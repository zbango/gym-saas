import { useState, type FormEvent, type ReactNode } from "react";
import { EyeIcon, EyeOffIcon, LockIcon, UserIcon } from "@gym-saas/ui";
import logoSrc from "../../assets/logo.png";

type LoginCredentials = {
  username: string;
  password: string;
};

type LoginPageProps = {
  brandName?: string;
  onSubmit?: (credentials: LoginCredentials) => void | Promise<void>;
  onForgotPassword?: () => void;
};

export function LoginPage({
  brandName = "Zeus Gym",
  onSubmit,
  onForgotPassword
}: LoginPageProps) {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [passwordVisible, setPasswordVisible] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!onSubmit) {
      return;
    }

    setSubmitting(true);
    try {
      await onSubmit({ username, password });
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <main
      className="relative isolate grid min-h-dvh place-items-center overflow-hidden bg-[linear-gradient(124deg,var(--color-brand-canvas)_0%,var(--color-brand-ambient)_53%,var(--color-brand-gold)_100%)] px-6 py-10 before:absolute before:inset-0 before:-z-10 before:bg-[radial-gradient(ellipse_72%_88%_at_-7%_-12%,rgba(0,0,0,0.7),transparent_63%)] before:content-[''] sm:px-4 sm:py-5"
    >
      <div
        aria-hidden="true"
        className="pointer-events-none absolute inset-0 -z-10 bg-[radial-gradient(circle_at_49%_46%,rgba(65,34,12,0.28)_0,rgba(65,34,12,0)_33%)]"
      />
      <section
        className="w-full max-w-[450px] rounded-[32px] border border-[rgba(255,228,92,0.08)] bg-[var(--color-brand-surface)] px-8 pb-8 pt-[42px] shadow-[0_28px_72px_rgba(28,14,4,0.42)] sm:rounded-[26px] sm:px-[26px] sm:pb-8 sm:pt-[38px]"
        aria-labelledby="login-title"
      >
        <BrandMark brandName={brandName} />

        <header className="my-[31px] text-center sm:mb-10 sm:mt-8">
          <h1
            id="login-title"
            className="m-0 text-[clamp(26px,3vw,30px)] font-extrabold leading-[1.15] tracking-[-0.04em] text-[var(--color-brand-accent)]"
          >
            Control de Gimnasio
          </h1>
          <p className="mt-[15px] text-lg leading-[1.2] font-medium text-[var(--color-brand-text)] sm:text-[19px]">
            Inicia sesión para continuar
          </p>
        </header>

        <form className="grid gap-5" onSubmit={(event) => void submit(event)}>
          <AuthField
            label="Usuario"
            name="username"
            placeholder="Ingresa tu usuario"
            autoComplete="username"
            value={username}
            onChange={setUsername}
            icon={<UserIcon />}
          />
          <AuthField
            label="Contraseña"
            name="password"
            type={passwordVisible ? "text" : "password"}
            placeholder="Ingresa tu contraseña"
            autoComplete="current-password"
            value={password}
            onChange={setPassword}
            icon={<LockIcon />}
            endAdornment={
              <button
                className="grid h-full w-[42px] flex-none place-items-center rounded-[10px] border-0 bg-transparent text-[var(--color-brand-muted)] transition-colors hover:text-[var(--color-brand-accent)] focus-visible:text-[var(--color-brand-accent)] focus-visible:outline-3 focus-visible:outline-offset-3 focus-visible:outline-[var(--color-brand-accent)]"
                type="button"
                onClick={() => setPasswordVisible((visible) => !visible)}
                aria-label={passwordVisible ? "Ocultar contraseña" : "Mostrar contraseña"}
                aria-pressed={passwordVisible}
              >
                {passwordVisible ? <EyeOffIcon /> : <EyeIcon />}
              </button>
            }
          />

          <button
            className="mt-[-1px] justify-self-end border-0 bg-transparent p-0 text-[13px] font-bold text-[var(--color-brand-accent)] hover:text-[var(--color-brand-accent-strong)] hover:underline focus-visible:outline-3 focus-visible:outline-offset-3 focus-visible:outline-[var(--color-brand-accent)]"
            type="button"
            onClick={onForgotPassword}
          >
            Olvidé mi contraseña
          </button>
          <button
            className="mt-1 min-h-11 rounded-[9px] border border-[rgba(255,221,86,0.4)] bg-[linear-gradient(100deg,var(--color-brand-button-start)_0%,var(--color-brand-button-middle)_42%,var(--color-brand-button-end)_100%)] px-4 text-base font-extrabold tracking-[-0.02em] text-[var(--color-brand-button-text)] transition-[filter,transform] hover:not-disabled:-translate-y-px hover:not-disabled:brightness-105 focus-visible:outline-3 focus-visible:outline-offset-3 focus-visible:outline-[var(--color-brand-accent)] disabled:cursor-wait disabled:opacity-70"
            type="submit"
            disabled={submitting}
          >
            {submitting ? "INICIANDO SESIÓN..." : "INICIAR SESIÓN"}
          </button>
        </form>
      </section>
    </main>
  );
}

function BrandMark({ brandName }: { brandName: string }) {
  return (
    <div className="grid min-h-40 place-items-center sm:min-h-[164px]">
      <img
        className="h-40 w-40 object-contain drop-shadow-[0_0_12px_rgba(255,194,18,0.22)] sm:h-[132px] sm:w-[132px]"
        src={logoSrc}
        alt={brandName}
      />
    </div>
  );
}

type AuthFieldProps = {
  label: string;
  name: string;
  type?: "text" | "password";
  placeholder: string;
  autoComplete: string;
  value: string;
  onChange: (value: string) => void;
  icon: ReactNode;
  endAdornment?: ReactNode;
};

function AuthField({
  label,
  name,
  type = "text",
  placeholder,
  autoComplete,
  value,
  onChange,
  icon,
  endAdornment
}: AuthFieldProps) {
  return (
    <div className="grid gap-2 text-sm font-bold text-[var(--color-brand-text)] sm:text-base">
      <label htmlFor={name}>{label}</label>
      <span className="flex h-[50px] items-center rounded-[9px] border border-[var(--color-brand-border)] bg-[var(--color-brand-input)] text-[var(--color-brand-muted)] transition-[border-color,box-shadow] focus-within:border-[var(--color-brand-accent)] focus-within:shadow-[0_0_0_3px_color-mix(in_srgb,var(--color-brand-accent)_24%,transparent)] sm:h-[58px]">
        <span className="grid w-[42px] flex-none place-items-center [&>svg]:h-[18px] [&>svg]:w-[18px]" aria-hidden="true">
          {icon}
        </span>
        <input
          id={name}
          name={name}
          type={type}
          value={value}
          onChange={(event) => onChange(event.target.value)}
          placeholder={placeholder}
          autoComplete={autoComplete}
          required
          className="h-full min-w-0 flex-1 border-0 bg-transparent text-base text-[var(--color-brand-text)] outline-0 placeholder:text-[#777d89] sm:text-[17px]"
        />
        {endAdornment}
      </span>
    </div>
  );
}
