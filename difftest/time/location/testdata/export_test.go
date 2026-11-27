package time

func GoLeanLocatedTime(wall, ext uint64, loc *Location) Time {
	return Time{wall: wall, ext: int64(ext), loc: loc}
}

func GoLeanLocatedParts(t Time) (uint64, uint64, *Location) {
	return t.wall, uint64(t.ext), t.loc
}

func GoLeanStripMono(t Time) Time {
	t.stripMono()
	return t
}

func GoLeanSetLoc(t Time, loc *Location) Time {
	t.setLoc(loc)
	return t
}
