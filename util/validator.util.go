package util

import (
	"github.com/google/uuid"
)

// HasAnyInvalidUUID checks if any of the given UUIDs is invalid
func HasAnyInvalidUUID(ids []string) bool {
	for _, id := range ids {
		if _, err := uuid.Parse(id); err != nil {
			return true
		}
	}
	return false
}
