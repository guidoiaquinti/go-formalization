package strings

func GoLeanIndexFunc(s string, f func(rune) bool, truth bool) int {
	return indexFunc(s, f, truth)
}
