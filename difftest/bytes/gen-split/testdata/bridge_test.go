package bytes

func GoLeanGenSplit(s, sep []byte, save, limit int) [][]byte {
	return genSplit(s, sep, save, limit)
}
