package strings

func GoLeanMakeStringFinder(pattern string) (string, [256]int, []int) {
	finder := makeStringFinder(pattern)
	return finder.pattern, finder.badCharSkip, finder.goodSuffixSkip
}
