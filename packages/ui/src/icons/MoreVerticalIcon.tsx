interface MoreVerticalIconProps {
  className?: string;
}

export function MoreVerticalIcon({ className }: MoreVerticalIconProps) {
  return (
    <svg width="1em" height="1em" className={className} viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
      <circle cx="12" cy="5" r="1.8" />
      <circle cx="12" cy="12" r="1.8" />
      <circle cx="12" cy="19" r="1.8" />
    </svg>
  );
}
