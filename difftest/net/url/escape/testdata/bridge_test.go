package url

func GoLeanEscapeClassification(c, mode byte) (bool, bool, byte) {
	return shouldEscape(c, encoding(mode)), ishex(c), unhex(c)
}

func GoLeanEscape(s string, mode byte) string {
	return escape(s, encoding(mode))
}
