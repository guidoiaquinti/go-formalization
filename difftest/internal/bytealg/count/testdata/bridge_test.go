package bytealg

func GoLeanBytealgCounts(slice []byte, text string, target byte) [4]int {
	return [4]int{Count(slice, target), CountString(text, target),
		countGeneric(slice, target), countGenericString(text, target)}
}
