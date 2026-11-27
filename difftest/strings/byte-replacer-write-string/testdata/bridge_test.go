package strings

import "io"

func GoLeanByteReplacerWriteString(mapping *[256]byte, writer io.Writer, input string) (int, error) {
	return (*byteReplacer)(mapping).WriteString(writer, input)
}
