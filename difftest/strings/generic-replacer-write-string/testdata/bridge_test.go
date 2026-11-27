package strings

import "io"

func GoLeanGenericReplacerWriteString(oldnew []string, writer io.Writer, input string) (int, error) {
	return makeGenericReplacer(oldnew).WriteString(writer, input)
}
