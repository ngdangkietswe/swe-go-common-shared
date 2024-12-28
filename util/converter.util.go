package util

import (
	"github.com/google/uuid"
	"github.com/samber/lo"
)

// Convert2UUID converts list of string to list of UUID
func Convert2UUID(ids []string) []uuid.UUID {
	distinct := lo.Uniq(ids)
	return lo.Map(distinct, func(id string, _ int) uuid.UUID {
		return uuid.MustParse(id)
	})
}

// Convert2String converts list of UUID to list of string
func Convert2String(ids []uuid.UUID) []string {
	distinct := lo.Uniq(ids)
	return lo.Map(distinct, func(id uuid.UUID, _ int) string {
		return id.String()
	})
}
