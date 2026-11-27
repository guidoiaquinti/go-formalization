package strings

func GoLeanASCIISetContains(raw *[256]bool, b byte) bool {
	return (*asciiSet)(raw).contains(b)
}
