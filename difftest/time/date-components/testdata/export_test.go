package time

func GoLeanDateComponents(year, month, day, hi, lo, base int) (uint64, bool, [4]int) {
	m := Month(month)
	h, l := norm(hi, lo, base)
	return uint64(dateToAbsDays(int64(year), m, day)), isLeap(year),
		[4]int{daysBefore(m), daysIn(m, year), h, l}
}
