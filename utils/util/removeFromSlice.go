package util

func RemoveFromSlice(slice []string, value string) []string {
	for i, v := range slice {
		if v == value {
			slice = append(slice[:i], slice[i+1:]...)
			return slice
		}
	}
	return slice
}
