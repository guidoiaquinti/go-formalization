package bytealg

func GoLeanRabinKarpHashes(slice []byte, text string) [8]uint32 {
	a, b := hashStr(slice)
	c, d := hashStr(text)
	e, f := hashStrRev(slice)
	g, h := hashStrRev(text)
	return [8]uint32{a, b, c, d, e, f, g, h}
}

func GoLeanRabinKarpFirstSlice(input, sep []byte) int  { return IndexRabinKarp(input, sep) }
func GoLeanRabinKarpFirstString(input, sep string) int { return IndexRabinKarp(input, sep) }
func GoLeanRabinKarpLastSlice(input, sep []byte) int   { return LastIndexRabinKarp(input, sep) }
func GoLeanRabinKarpLastString(input, sep string) int  { return LastIndexRabinKarp(input, sep) }
