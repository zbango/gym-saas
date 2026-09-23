package main

import (
	"time"

	"github.com/zbango/gym-saas/go/core/domain"
)

func formatOptionalTimestamp(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return domain.FormatTimestamp(value)
}
