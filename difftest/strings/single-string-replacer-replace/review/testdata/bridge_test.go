package strings

func GoLeanReviewSingleReplacer(pattern, value, input string) string {
	return makeSingleStringReplacer(pattern, value).Replace(input)
}
