package domain

import (
	"fmt"
	"time"
)

// FormatTimestamp returns the canonical persisted representation for a time.
func FormatTimestamp(value time.Time) string {
	return value.UTC().Format(time.RFC3339Nano)
}

// ParseTimestamp reads the canonical UTC timestamp representation.
func ParseTimestamp(value string) (time.Time, error) {
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse timestamp: %w", err)
	}
	if parsed.Location() != time.UTC || parsed.Format(time.RFC3339Nano) != value {
		return time.Time{}, fmt.Errorf("parse timestamp: value must be canonical UTC RFC 3339")
	}
	return parsed, nil
}
