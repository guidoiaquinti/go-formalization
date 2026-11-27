package strings

func GoLeanByteReplacerReplace(mapping *[256]byte, input string) string {
	return (*byteReplacer)(mapping).Replace(input)
}
