package strings

func GoLeanReviewSingleMaker(pattern, value string) (string, string, [256]int, []int) {
	replacer := makeSingleStringReplacer(pattern, value)
	return replacer.finder.pattern, replacer.value, replacer.finder.badCharSkip, replacer.finder.goodSuffixSkip
}
