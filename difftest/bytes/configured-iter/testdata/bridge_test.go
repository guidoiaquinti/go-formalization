package bytes

import "iter"

func GoLeanConfiguredSplit(s, sep []byte, save int) iter.Seq[[]byte] {
	return splitSeq(s, sep, save)
}
