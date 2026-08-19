package domain

import (
	"testing"
	"time"
)

func TestFormatTimestampUsesCanonicalUTC(t *testing.T) {
	value := time.Date(2026, time.August, 19, 10, 11, 12, 123456789, time.FixedZone("ECT", -5*60*60))
	if got, want := FormatTimestamp(value), "2026-08-19T15:11:12.123456789Z"; got != want {
		t.Fatalf("FormatTimestamp() = %q, want %q", got, want)
	}
}

func TestParseTimestamp(t *testing.T) {
	value, err := ParseTimestamp("2026-08-19T15:11:12.123456789Z")
	if err != nil {
		t.Fatalf("ParseTimestamp returned error: %v", err)
	}
	if value.Location() != time.UTC {
		t.Fatalf("location = %v, want UTC", value.Location())
	}
}

func TestParseTimestampRejectsNonCanonicalValues(t *testing.T) {
	for _, value := range []string{
		"2026-08-19T10:11:12-05:00",
		"2026-08-19T15:11:12+00:00",
		"2026-08-19 15:11:12Z",
	} {
		if _, err := ParseTimestamp(value); err == nil {
			t.Fatalf("ParseTimestamp(%q) succeeded", value)
		}
	}
}
