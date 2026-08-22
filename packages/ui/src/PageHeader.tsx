import type { ReactNode } from "react";

export type PageHeaderProps = {
  icon: ReactNode;
  title: string;
  description?: string;
  actions?: ReactNode;
  className?: string;
};

export function PageHeader({ icon, title, description, actions, className }: PageHeaderProps) {
  return (
    <section className={`gs-page-header ${className ?? ""}`.trim()}>
      <div className="gs-page-header-copy">
        <span className="gs-page-header-icon" aria-hidden="true">{icon}</span>
        <div>
          <h1>{title}</h1>
          {description ? <p>{description}</p> : null}
        </div>
      </div>
      {actions ? <div className="gs-page-header-actions">{actions}</div> : null}
    </section>
  );
}
