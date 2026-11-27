package url

func GoLeanValidEncoded(s string, mode byte) bool {
	return validEncoded(s, encoding(mode))
}

func GoLeanSetPath(u *URL, s string) error     { return u.setPath(s) }
func GoLeanSetFragment(u *URL, s string) error { return u.setFragment(s) }
