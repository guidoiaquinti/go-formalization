package bytes

func GoLeanLastIndexFunc(s []byte, f func(rune) bool, truth bool) int {
	return lastIndexFunc(s, f, truth)
}
