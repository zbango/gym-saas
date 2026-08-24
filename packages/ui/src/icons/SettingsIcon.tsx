interface SettingsIconProps {
  className?: string;
}

export function SettingsIcon({ className }: SettingsIconProps) {
  return (
    <svg width="1em" height="1em" className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317a1.65 1.65 0 013.35 0 1.65 1.65 0 002.475 1.428 1.65 1.65 0 012.9 1.674 1.65 1.65 0 001.428 2.475 1.65 1.65 0 010 3.35 1.65 1.65 0 00-1.428 2.475 1.65 1.65 0 01-2.9 1.674 1.65 1.65 0 00-2.475 1.428 1.65 1.65 0 01-3.35 0 1.65 1.65 0 00-2.475-1.428 1.65 1.65 0 01-2.9-1.674 1.65 1.65 0 00-1.428-2.475 1.65 1.65 0 010-3.35A1.65 1.65 0 004.95 7.419a1.65 1.65 0 012.9-1.674 1.65 1.65 0 002.475-1.428z" />
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15.5a3.5 3.5 0 100-7 3.5 3.5 0 000 7z" />
    </svg>
  );
}
