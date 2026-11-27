package strings

func GoLeanMakeSingleStringReplacer(pattern, value string) (string, string, [256]int, []int) {
	replacer := makeSingleStringReplacer(pattern, value)
	finder := replacer.finder
	return finder.pattern, replacer.value, finder.badCharSkip, finder.goodSuffixSkip
}
