interface ActivityIconProps {
  className?: string;
}

export function ActivityIcon({ className }: ActivityIconProps) {
  return (
    <svg
      width="1em"
      height="1em"
      className={className}
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M13 10V3L4 14h7v7l9-11h-7z"
      />
    </svg>
  );
}
