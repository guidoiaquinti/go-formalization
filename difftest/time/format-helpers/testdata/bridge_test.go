package time

func GoLeanFmtInt(buf []byte, value uint64) int {
	return fmtInt(buf, value)
}

func GoLeanFmtFrac(buf []byte, value uint64, precision int) (int, uint64) {
	return fmtFrac(buf, value, precision)
}
