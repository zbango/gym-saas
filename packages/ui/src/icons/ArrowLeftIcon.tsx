interface ArrowLeftIconProps {
  className?: string;
}

export function ArrowLeftIcon({ className }: ArrowLeftIconProps) {
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
        d="M10 19l-7-7m0 0l7-7m-7 7h18"
      />
    </svg>
  );
}
