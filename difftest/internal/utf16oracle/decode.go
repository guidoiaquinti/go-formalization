// Package utf16oracle defines the surrogate-block recurrence independently of Go's decoder.
package utf16oracle

func Decode(words []uint16) []rune {
	var runes []rune
	for len(words) > 0 {
		first := rune(words[0])
		words = words[1:]
		decoded := first
		switch first >> 10 {
		case 54:
			decoded = 65533
			if len(words) > 0 && words[0]>>10 == 55 {
				decoded = 65536 + 1024*(first%1024) + rune(words[0]%1024)
				words = words[1:]
			}
		case 55:
			decoded = 65533
		}
		runes = append(runes, decoded)
	}
	return runes
}
