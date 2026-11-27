package time

func GoLeanArithmeticTime(wall uint64, ext int64, loc *Location) Time {
	return Time{wall: wall, ext: ext, loc: loc}
}

func GoLeanArithmeticFields(t Time) (uint64, int64, *Location) {
	return t.wall, t.ext, t.loc
}

func GoLeanArithmeticSubMono(t, u int64) Duration {
	return subMono(t, u)
}
