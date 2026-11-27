package strings

func GoLeanStringFinderNext(pattern, text string) int {
	return makeStringFinder(pattern).next(text)
}
