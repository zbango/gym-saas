import type { ReactNode } from "react";

export type PageHeaderProps = {
  icon: ReactNode;
  title: string;
  description?: string;
  actions?: ReactNode;
  className?: string;
  iconClassName?: string;
};

export function PageHeader({ icon, title, description, actions, className, iconClassName }: PageHeaderProps) {
  return (
    <section className={`mb-[46px] flex items-start justify-between gap-7 max-[920px]:grid ${className ?? ""}`.trim()}>
      <div className="flex items-start gap-[18px]">
        <span className={`grid place-items-center pt-1 text-[#56616d] [&>svg]:h-[37px] [&>svg]:w-[37px] ${iconClassName ?? ""}`.trim()} aria-hidden="true">{icon}</span>
        <div>
          <h1 className="m-0 text-[38px] leading-[1.15] font-extrabold tracking-[-0.04em] max-[1260px]:text-[33px] max-[560px]:text-[28px]">{title}</h1>
          {description ? <p className="mt-2.5 text-[21px] text-[#5f6974] max-[560px]:text-[17px]">{description}</p> : null}
        </div>
      </div>
      {actions ? <div className="flex gap-[14px] max-[920px]:flex-wrap">{actions}</div> : null}
    </section>
  );
}
