package strings

func GoLeanAppendSliceWriterWriteString(buffer *[]byte, input string) (int, error) {
	return (*appendSliceWriter)(buffer).WriteString(input)
}
