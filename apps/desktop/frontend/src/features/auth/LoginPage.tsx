import { useState, type FormEvent, type ReactNode } from "react";
import logoSrc from "../../assets/logo.png";
import "./login-page.css";

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
    <main className="login-page">
      <section className="login-card" aria-labelledby="login-title">
        <BrandMark brandName={brandName} />

        <header className="login-intro">
          <h1 id="login-title">Control de Gimnasio</h1>
          <p>Inicia sesión para continuar</p>
        </header>

        <form className="login-form" onSubmit={(event) => void submit(event)}>
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
                className="password-toggle"
                type="button"
                onClick={() => setPasswordVisible((visible) => !visible)}
                aria-label={passwordVisible ? "Ocultar contraseña" : "Mostrar contraseña"}
                aria-pressed={passwordVisible}
              >
                <EyeIcon crossed={passwordVisible} />
              </button>
            }
          />

          <button className="forgot-password" type="button" onClick={onForgotPassword}>
            Olvidé mi contraseña
          </button>
          <button className="login-submit" type="submit" disabled={submitting}>
            {submitting ? "INICIANDO SESIÓN..." : "INICIAR SESIÓN"}
          </button>
        </form>
      </section>
    </main>
  );
}

function BrandMark({ brandName }: { brandName: string }) {
  return (
    <div className="brand-mark">
      <img className="brand-logo" src={logoSrc} alt={brandName} />
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
    <div className="auth-field">
      <label htmlFor={name}>{label}</label>
      <span className="auth-input-shell">
        <span className="auth-input-icon" aria-hidden="true">{icon}</span>
        <input
          id={name}
          name={name}
          type={type}
          value={value}
          onChange={(event) => onChange(event.target.value)}
          placeholder={placeholder}
          autoComplete={autoComplete}
          required
        />
        {endAdornment}
      </span>
    </div>
  );
}

function UserIcon() {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.9" aria-hidden="true">
      <circle cx="12" cy="8" r="3.5" />
      <path d="M4.5 20c.8-4.1 3.3-6.2 7.5-6.2s6.7 2.1 7.5 6.2" />
    </svg>
  );
}

function LockIcon() {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.9" aria-hidden="true">
      <rect x="5.5" y="10" width="13" height="10" rx="1.6" />
      <path d="M8.5 10V7.4a3.5 3.5 0 0 1 7 0V10M12 14v2.4" />
    </svg>
  );
}

function EyeIcon({ crossed }: { crossed: boolean }) {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.9" aria-hidden="true">
      <path d="M2.7 12s3.1-5.2 9.3-5.2 9.3 5.2 9.3 5.2-3.1 5.2-9.3 5.2S2.7 12 2.7 12Z" />
      <circle cx="12" cy="12" r="2.5" />
      {crossed ? <path d="m4 4 16 16" /> : null}
    </svg>
  );
}
