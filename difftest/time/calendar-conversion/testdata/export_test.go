package time

func GoLeanCalendarConversion(days, century uint64, y, jf, m, yd, l int) (uint64, [10]int) {
	d := absDays(days)
	c, cy, ay := d.split()
	dateYear, dateMonth, monthDay := d.date()
	dayYear, yearDay := d.yearYday()
	janFeb := absJanFeb(jf)
	return uint64(c), [10]int{int(cy), int(ay), absCentury(century).year(absCyear(y), janFeb),
		int(absMonth(m).month(janFeb)), absYday(yd).yday(janFeb, absLeap(l)),
		dateYear, int(dateMonth), monthDay, dayYear, yearDay}
}
