package utils

// StringSliceContains returns true if a string is contained in a slice
func StringSliceContains(haystack []string, needle string) bool {
	for _, element := range haystack {
		if element == needle {
			return true
		}
	}

	return false
}
