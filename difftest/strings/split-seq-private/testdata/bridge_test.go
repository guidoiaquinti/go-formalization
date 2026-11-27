package strings

import "iter"

func GoLeanSplitSeqPrivate(input, sep string, save int) iter.Seq[string] {
	return splitSeq(input, sep, save)
}
