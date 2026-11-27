package filepath

func GoLeanMatchChunk(pattern, input string) (string, bool, error) {
	return matchChunk(pattern, input)
}
