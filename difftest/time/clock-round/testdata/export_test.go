package time

func GoLeanClockDiv(value Time, duration Duration) (int, Duration) {
	return div(value, duration)
}

func GoLeanLessThanHalf(a, b Duration) bool {
	return lessThanHalf(a, b)
}
