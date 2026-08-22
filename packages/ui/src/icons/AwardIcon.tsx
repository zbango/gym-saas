import React from 'react';

interface AwardIconProps {
  className?: string;
}

export const AwardIcon: React.FC<AwardIconProps> = ({
  className,
}) => (
  <svg
      width="1em"
      height="1em"
    className={className}
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <circle cx="12" cy="8" r="6" />
    <path d="M15.477 12.89 17 22l-5-3-5 3 1.523-9.11" />
  </svg>
);
