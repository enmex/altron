package utils

func ContainsFunc[T any](slice []T, pred func(T) bool) bool {
	for idx := range slice {
		if pred(slice[idx]) {
			return true
		}
	}
	return false
}

func ToggleElement[T any](slice []T, el T, comparisonFunc func(i, j T) bool) []T {
	for idx, v := range slice {
		if comparisonFunc(el, v) {
			return append(slice[:idx], slice[idx+1:]...)
		}
	}
	return append(slice, el)
}

func DeleteElement[T any](slice []T, el T, comparisonFunc func(i, j T) bool) []T {
	for idx, v := range slice {
		if comparisonFunc(el, v) {
			slice[idx] = slice[len(slice)-1]
			return slice[:len(slice)-1]
		}
	}
	return slice
}

func Difference[T struct{}](a, b []T) []T {
	mb := make(map[T]bool, 0)
	for _, x := range b {
		mb[x] = true
	}
	diff := make([]T, 0)
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
