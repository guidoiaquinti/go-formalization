package slices

import "cmp"

func GoLeanXorshiftNext(word uint64, present bool) (uint64, uint64) {
	state := xorshift(word)
	var receiver *xorshift
	if present {
		receiver = &state
	}
	value := receiver.Next()
	return value, uint64(state)
}

func GoLeanNextPowerOfTwo(length int) uint { return nextPowerOfTwo(length) }

func GoLeanSortIsNaN[T cmp.Ordered](x T) bool { return isNaN(x) }
