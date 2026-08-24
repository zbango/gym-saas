export type TabItem<Id extends string> = {
  id: Id;
  label: string;
  disabled?: boolean;
};

export type TabsProps<Id extends string> = {
  tabs: readonly TabItem<Id>[];
  activeTab: Id;
  onChange: (tab: Id) => void;
  ariaLabel: string;
  className?: string;
};

export function Tabs<Id extends string>({ tabs, activeTab, onChange, ariaLabel, className }: TabsProps<Id>) {
  return (
    <div className={`flex min-w-max items-end gap-7 border-b border-[#e6e9ed] ${className ?? ""}`.trim()} role="tablist" aria-label={ariaLabel}>
      {tabs.map((tab) => {
        const selected = tab.id === activeTab;
        return (
          <button
            key={tab.id}
            id={`tab-${tab.id}`}
            role="tab"
            type="button"
            aria-selected={selected}
            aria-controls={`tabpanel-${tab.id}`}
            disabled={tab.disabled}
            onClick={() => onChange(tab.id)}
            className={[
              "relative min-h-[72px] shrink-0 border-b-4 px-1 text-[20px] font-bold transition-colors focus-visible:outline-3 focus-visible:outline-offset-[-3px] focus-visible:outline-[var(--color-brand-accent)] disabled:cursor-not-allowed disabled:opacity-50",
              selected
                ? "border-[var(--color-brand-accent)] text-[#895c2a]"
                : "border-transparent text-[#a2a9b2] hover:text-[#626d79]"
            ].join(" ")}
          >
            {tab.label}
          </button>
        );
      })}
    </div>
  );
}
