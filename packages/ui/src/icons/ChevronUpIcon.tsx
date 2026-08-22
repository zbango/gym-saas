import React from 'react';

interface ChevronUpIconProps {
  className?: string;
}

export const ChevronUpIcon: React.FC<ChevronUpIconProps> = ({
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
    <polyline points="18 15 12 9 6 15" />
  </svg>
);
