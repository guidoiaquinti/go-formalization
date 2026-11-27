package strings

import "io"

func GoLeanSingleStringReplacerWriteString(pattern, value, input string, writer io.Writer) (int, error) {
	return makeSingleStringReplacer(pattern, value).WriteString(writer, input)
}
