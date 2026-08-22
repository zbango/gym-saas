import { useState, type FormEvent, type ReactNode } from "react";
import { EyeIcon, EyeOffIcon, LockIcon, UserIcon } from "@gym-saas/ui";
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
                {passwordVisible ? <EyeOffIcon /> : <EyeIcon />}
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
