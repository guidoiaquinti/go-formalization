package time

func GoLeanTimeFromWords(wall uint64, ext int64, loc *Location) Time {
	return Time{wall: wall, ext: ext, loc: loc}
}

func GoLeanTimeWords(t Time) (int32, int64, int64, int64) {
	return t.nsec(), t.sec(), t.unixSec(), t.mono()
}
