package strings

func GoLeanSingleStringReplacerReplace(pattern, value, input string) string {
	return makeSingleStringReplacer(pattern, value).Replace(input)
}
