package util


func SliceContains(slice []string, searchString string) bool {
	for _, str := range slice {
		if str == searchString {
			return true
		}
	}
	return false
}
