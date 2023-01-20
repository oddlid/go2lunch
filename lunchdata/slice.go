package lunchdata

func sliceIndex[T any](slice []T, match func(T) bool) int {
	for i := range slice {
		if match(slice[i]) {
			return i
		}
	}
	return -1
}

// deleteByIndex deletes the element at the given index, if valid. Order is NOT kept.
// note: this should not be used with slices containing pointers, as they will not get GC'ed
func deleteByIndex[T any](slice []T, index int) []T {
	lastIndex := len(slice) - 1

	if index < 0 || index > lastIndex {
		return slice
	}

	if index != lastIndex {
		slice[index] = slice[lastIndex]
	}

	return slice[:lastIndex]
}
