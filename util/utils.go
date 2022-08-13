package util

func Contains[E comparable](v []E, find E) bool {
	for _, e := range v {
		if e == find {
			return true
		}
	}
	return false
}
