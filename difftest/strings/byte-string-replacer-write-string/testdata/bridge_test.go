package strings

import "io"

func GoLeanByteStringReplacerWriteString(mapping *[256][]byte, writer io.Writer, input string) (int, error) {
	var replacer *byteStringReplacer
	if mapping != nil {
		replacer = &byteStringReplacer{replacements: *mapping}
	}
	return replacer.WriteString(writer, input)
}
