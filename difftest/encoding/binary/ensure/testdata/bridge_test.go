package binary

func GoLeanEnsure(buf []byte, n int) (buf2, pos []byte) {
	return ensure(buf, n)
}
