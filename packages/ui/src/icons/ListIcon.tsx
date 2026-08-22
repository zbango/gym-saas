interface ListIconProps {
  className?: string;
}

export function ListIcon({ className }: ListIconProps) {
  return (
    <svg
      width="1em"
      height="1em"
      className={className}
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M4 6h16M4 12h16M4 18h16"
      />
    </svg>
  );
}

