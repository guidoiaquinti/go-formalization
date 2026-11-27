package time

func GoLeanClockFields(seconds, days, century uint64, year, day uint32) [9]uint64 {
	abs := absSeconds(seconds)
	h, m, s := abs.clock()
	month, monthDay := absYday(day).split()
	return [9]uint64{uint64(abs.days()), uint64(h), uint64(m), uint64(s),
		uint64(absDays(days).weekday()), uint64(absYday(day).janFeb()),
		uint64(absCentury(century).leap(absCyear(year))), uint64(month), uint64(monthDay)}
}
