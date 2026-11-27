package strings

func GoLeanByteStringReplacerReplace(mapping *[256][]byte, input string) string {
	var replacer *byteStringReplacer
	if mapping != nil {
		replacer = &byteStringReplacer{replacements: *mapping}
		for b, replacement := range mapping {
			if replacement != nil {
				replacer.toReplace = append(replacer.toReplace, string([]byte{byte(b)}))
			}
		}
	}
	return replacer.Replace(input)
}
