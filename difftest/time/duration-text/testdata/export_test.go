package time

func GoLeanDurationFormat(d Duration, backing *[32]byte) int {
	return d.format(backing)
}
