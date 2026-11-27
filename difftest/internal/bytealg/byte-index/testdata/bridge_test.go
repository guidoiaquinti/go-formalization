package bytealg

func GoLeanBytealgByteIndices(slice []byte, text string, target byte) [4]int {
	return [4]int{IndexByte(slice, target), IndexByteString(text, target),
		LastIndexByte(slice, target), LastIndexByteString(text, target)}
}
