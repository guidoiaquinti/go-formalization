package strings

func GoLeanAppendSliceWriterWrite(buffer *[]byte, input []byte) (int, error) {
	return (*appendSliceWriter)(buffer).Write(input)
}
