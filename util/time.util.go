package util

import (
	"github.com/samber/lo"
	"time"
)

const (
	LayoutISO         = "2006-01-02"
	LayoutISOWithTime = "2006-01-02T15:04:05Z"
)

// Format is a function that formats a time.Time to a string with a given layout.
func Format(t *time.Time, layout string) string {
	return lo.Ternary(t == nil, "", t.Format(layout))
}
