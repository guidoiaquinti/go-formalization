package strings

func GoLeanGenericReplacerReplace(oldnew []string, input string) string {
	return makeGenericReplacer(oldnew).Replace(input)
}
