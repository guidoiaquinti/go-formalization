package time

func GoLeanClockMutations(wall, ext uint64, loc *Location, argument uint64) (Time, Time) {
	first := Time{wall: wall, ext: int64(ext), loc: loc}
	second := first
	first.setMono(int64(argument))
	second.addSec(int64(argument))
	return first, second
}

func GoLeanClockMutationParts(t Time) (uint64, uint64, *Location) {
	return t.wall, uint64(t.ext), t.loc
}
