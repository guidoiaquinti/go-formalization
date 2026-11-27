package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Wall, Ext        uint64
	Sec, Nsec, Units int64
	Want             string
}

func main() {
	cases := []fixture{
		{999999999, 62135596799, 0, -1, -1,
			"-1|-1|-1|18446744073709551615,62135596800|999999999,62135596799|999000000,62135596799|999999000,62135596799"},
		{987654321, 72135596800, math.MaxInt64, 1000000000, math.MinInt64,
			"10000000000987|10000000000987654|-8446744072721897295|1000000000,9223372098990372607|0,9223372098990372608|192000000,18437520763808293640|224192000,18446734912473111561"},
		{1<<63 + 1<<30 + 500123456, 7, math.MinInt64, -1000000001, math.MaxInt64,
			"-2682287998500|-2682287998499877|-2682287998499876544|18446744072709551615,9223372098990372608|999999999,9223372098990372606|807000000,9223434172451575|775807000,9285507633654"},
		{1073741823, 0, 0, 2147483648, -1001,
			"-62135596798927|-62135596798926259|-6795364577797603329|18446744071562067968,62135596800|147483648,62135596802|999000000,62135596798|998999000,62135596799"},
		{1, 9223372098990372607, 0, 0, 0,
			"-1000|-1000000|-999999999|0,62135596800|0,62135596800|0,62135596800|0,62135596800"},
	}
	var input bytes.Buffer
	truth, err := os.CreateTemp("", "review-time-epoch-")
	must(err)
	defer os.Remove(truth.Name())
	enc := json.NewEncoder(truth)
	for _, c := range cases {
		fmt.Fprintln(&input, c.Wall, c.Ext, c.Sec, c.Nsec, c.Units)
		must(enc.Encode(c))
	}
	must(truth.Close())
	must(os.Setenv("GOLEAN_TIME_EPOCH_ROWS", truth.Name()))
	must(check.TestInPackage("time", "TestGoLeanEpochUnits", "difftest/time/location/testdata/export_test.go", "difftest/time/epoch-units/testdata/export_test.go",
		"difftest/time/epoch-units/testdata/epoch_test.go"))
	rows := check.LeanRows("godiff-time-epoch-units", len(cases), &input)
	for i, row := range rows {
		if row != cases[i].Want {
			log.Fatalf("row %d got %q, want %q", i, row, cases[i].Want)
		}
	}
	for _, m := range []struct {
		row, col int
		value    string
	}{
		{0, 0, "0"},                                // round pre-epoch value toward zero at the end
		{0, 1, "0"},                                // round pre-epoch value toward zero at the end
		{1, 2, "9223372036854775807"},              // saturate nanoseconds at the maximum
		{4, 0, "9223372036854775807"},              // saturate milliseconds on overflow
		{4, 1, "9223372036854775807"},              // saturate microseconds on overflow
		{2, 2, "-6795364571371221696"},             // read monotonic storage as seconds
		{3, 3, "2147483648,62135596800"},           // zero-extend the private signed nanoseconds
		{1, 4, "1000000000,9223372098990372607"},   // leave positive nanoseconds unnormalized
		{0, 4, "18446744073709551615,62135596800"}, // leave negative nanoseconds unnormalized
		{0, 5, "18446744073708551616,62135596800"}, // split units by truncation without normalization
		{0, 6, "18446744073709550616,62135596800"},
	} {
		if strings.Split(rows[m.row], "|")[m.col] == m.value {
			log.Fatalf("row %d column %d accepted mutant %q", m.row, m.col, m.value)
		}
	}
	log.Print("five independent pinned Go and Lean rows passed and rejected eleven wrong observations")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
