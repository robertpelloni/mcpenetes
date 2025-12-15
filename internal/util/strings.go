package util

import "strings"

// CaseInsensitiveContains checks if s contains substr in a case-insensitive manner
func CaseInsensitiveContains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
