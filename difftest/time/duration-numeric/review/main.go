package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	cases := []struct {
		duration, step int64
		want           string
	}{
		{-999, 0, "-999|0|0|13740696507596399576|13713989419439935718|13687303298590655802|-999|-999|999"},
		{-999999, -1, "-999999|-999|0|13785626541160459129|13758912388823352245|13732218732848053099|-999999|-999999|999999"},
		{-3, 2, "-3|0|0|13702700021751366111|13676161994314988987|13649656204900468252|-2|-4|3"},
		{3, 2, "3|0|0|4479327984896590303|4452789957460213179|4426284168045692444|2|4|3"},
		{math.MaxInt64, 2, "9223372036854775807|9223372036854775|9223372036854|4756133310154331797|4729434092608672722|4702756367077584574|9223372036854775806|9223372036854775807|9223372036854775807"},
		{math.MinInt64 + 1, 4, "-9223372036854775807|-9223372036854775|-9223372036854|13979505347009107605|13952806129463448530|13926128403932360382|-9223372036854775804|-9223372036854775808|9223372036854775807"},
		{math.MinInt64, 1, "-9223372036854775808|-9223372036854775|-9223372036854|13979505347009107605|13952806129463448530|13926128403932360382|-9223372036854775808|-9223372036854775808|9223372036854775807"},
		{4, 3, "4|0|0|4481413732884731541|4454714515339072466|4428036789807984318|3|3|4"},
		{-814638300595526980, 12345, "-814638300595526980|-814638300595526|-814638300595|13963488926586649146|13936922907580124973|13910387260490780720|-814638300595519980|-814638300595532325|814638300595526980"},
	}
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&input, c.duration, c.step)
		d, m := time.Duration(c.duration), time.Duration(c.step)
		got := fmt.Sprintf("%d|%d|%d|%d|%d|%d|%d|%d|%d", d.Nanoseconds(), d.Microseconds(), d.Milliseconds(),
			math.Float64bits(d.Seconds()), math.Float64bits(d.Minutes()), math.Float64bits(d.Hours()),
			d.Truncate(m), d.Round(m), d.Abs())
		if got != c.want {
			log.Fatalf("Go at %d %d = %q, want %q", c.duration, c.step, got, c.want)
		}
	}
	rows := check.LeanRows("godiff-time-duration-numeric", len(cases), &input)
	for i, got := range rows {
		if got != cases[i].want {
			log.Fatalf("Lean row %d = %q, want %q", i, got, cases[i].want)
		}
	}

	reject := func(row, column int, mutant any) {
		if fmt.Sprint(mutant) == strings.Split(rows[row], "|")[column] {
			log.Fatalf("row %d column %d did not reject %v", row, column, mutant)
		}
	}
	reject(6, 0, uint64(cases[6].duration)) // incorrectly expose an unsigned count
	reject(0, 1, -1)                        // floor rounding instead of truncation
	reject(1, 2, -1)
	for column, scale := range []float64{1e9, 6e10, 3.6e12} {
		reject(8, column+3, math.Float64bits(float64(cases[8].duration)/scale))
	}
	reject(2, 6, -4)                   // truncate toward negative infinity
	reject(2, 7, -2)                   // resolve a negative tie toward zero
	reject(3, 7, 2)                    // resolve a positive tie toward zero
	reject(4, 7, int64(math.MinInt64)) // wrap positive overflow
	reject(7, 7, 6)                    // round a below-half remainder up for odd step
	reject(6, 8, -cases[6].duration)   // negate minimum without saturating
	log.Print("nine independent Go and Lean rows passed and rejected twelve wrong observations")
}
