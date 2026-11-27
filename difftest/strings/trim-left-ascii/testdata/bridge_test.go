package strings

func GoLeanTrimLeftASCII(input string, members []byte, nilTable bool) string {
	if nilTable {
		return trimLeftASCII(input, nil)
	}
	var table asciiSet
	for _, member := range members {
		table[member] = true
	}
	return trimLeftASCII(input, &table)
}
