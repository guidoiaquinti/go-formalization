package strings

func GoLeanTrimRightASCII(input string, members []byte, nilTable bool) string {
	if nilTable {
		return trimRightASCII(input, nil)
	}
	var table asciiSet
	for _, member := range members {
		table[member] = true
	}
	return trimRightASCII(input, &table)
}
