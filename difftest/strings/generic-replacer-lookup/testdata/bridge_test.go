package strings

func GoLeanGenericReplacerLookup(oldnew []string, text string, ignoreRoot bool) (string, int, bool) {
	return makeGenericReplacer(oldnew).lookup(text, ignoreRoot)
}
