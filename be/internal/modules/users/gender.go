package users

import "strings"

// ValidGender reports whether g is a supported gender code.
func ValidGender(g string) bool {
	switch strings.TrimSpace(g) {
	case "male", "female", "other":
		return true
	default:
		return false
	}
}
