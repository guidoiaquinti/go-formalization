package bytes

func GoLeanTrimRightASCII(input []byte, members []byte, nilTable bool) []byte {
	if nilTable {
		return trimRightASCII(input, nil)
	}
	var table asciiSet
	for _, member := range members {
		table[member] = true
	}
	return trimRightASCII(input, &table)
}
