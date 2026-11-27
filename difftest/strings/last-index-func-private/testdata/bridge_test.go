package strings

func GoLeanLastIndexFunc(s string, f func(rune) bool, truth bool) int {
	return lastIndexFunc(s, f, truth)
}
