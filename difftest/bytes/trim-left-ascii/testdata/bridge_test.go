package bytes

func GoLeanTrimLeftASCII(input []byte, members []byte, nilTable bool) []byte {
	if nilTable {
		return trimLeftASCII(input, nil)
	}
	var table asciiSet
	for _, member := range members {
		table[member] = true
	}
	return trimLeftASCII(input, &table)
}
