package time

func GoLeanUnixTime(seconds int64, nanoseconds int32) Time {
	return unixTime(seconds, nanoseconds)
}
