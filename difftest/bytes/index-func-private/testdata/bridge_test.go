package bytes

func GoLeanIndexFunc(s []byte, f func(rune) bool, truth bool) int {
	return indexFunc(s, f, truth)
}
