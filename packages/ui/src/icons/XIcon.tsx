interface XIconProps {
  className?: string;
}

export function XIcon({ className }: XIconProps) {
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
        d="M6 18L18 6M6 6l12 12"
      />
    </svg>
  );
}
