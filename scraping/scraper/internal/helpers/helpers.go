package helpers

import "strings"

func ToSnakeCase(str string) string {
	lowerCasedStr := strings.ToLower(str)
	fragments := strings.Fields(lowerCasedStr)
	return strings.Join(fragments, "_")
}
