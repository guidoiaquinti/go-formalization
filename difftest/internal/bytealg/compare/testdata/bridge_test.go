package bytealg

func GoLeanBytealgCompare(left, right []byte, a, b string) ([3]int, bool) {
	return [3]int{Compare(left, right), CompareString(a, b), abigen_runtime_cmpstring(a, b)},
		Equal(left, right)
}
