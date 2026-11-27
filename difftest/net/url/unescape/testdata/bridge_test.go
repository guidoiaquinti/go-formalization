package url

func GoLeanUnescape(s string, mode byte) (string, error) {
	return unescape(s, encoding(mode))
}
