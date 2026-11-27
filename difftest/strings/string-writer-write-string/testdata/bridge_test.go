package strings

import "io"

func GoLeanStringWriterWriteString(writer io.Writer, input string) (int, error) {
	return (stringWriter{writer}).WriteString(input)
}
