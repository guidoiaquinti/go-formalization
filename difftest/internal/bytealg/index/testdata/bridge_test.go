package bytealg

func GoLeanBytealgIndex(slice, needle []byte, text, pattern string) (int, int) {
	return Index(slice, needle), IndexString(text, pattern)
}
