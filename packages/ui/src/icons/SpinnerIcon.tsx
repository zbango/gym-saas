interface SpinnerIconProps {
  className?: string;
}

export function SpinnerIcon({ className }: SpinnerIconProps) {
  return (
    <svg
      width="1em"
      height="1em"
      className={className}
      viewBox="0 0 24 24"
      aria-label="Cargando"
      role="img"
    >
      <circle cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" opacity="0.25" />
      <path fill="currentColor" opacity="0.75" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z">
        <animateTransform attributeName="transform" attributeType="XML" type="rotate" from="0 12 12" to="360 12 12" dur="1s" repeatCount="indefinite" />
      </path>
    </svg>
  );
}
