package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Date, Clock   uint16
	Seconds       int64
	Nanos, Offset int
	Want          string
}

func civil(t time.Time) [6]int {
	y, m, d := t.Date()
	h, minute, s := t.Clock()
	return [6]int{y, int(m), d, h, minute, s}
}

func encode(c [6]int) (uint16, uint16) {
	return uint16(c[2] + 32*c[1] + 512*(c[0]-1980)), uint16(c[5]/2 + 32*c[4] + 2048*c[3])
}

func main() {
	f, err := os.CreateTemp("", "golean-zip-dos")
	must(err)
	defer os.Remove(f.Name())
	defer f.Close()
	enc := json.NewEncoder(f)
	compare := func(label string, cases []fixture) {
		var input bytes.Buffer
		for i, c := range cases {
			local := civil(time.Unix(c.Seconds+int64(c.Offset), int64(c.Nanos)).UTC())
			utc := civil(time.Unix(c.Seconds, int64(c.Nanos)).UTC())
			fmt.Fprint(&input, c.Date, " ", c.Clock, " ", c.Seconds, " ", c.Nanos)
			for _, fields := range [][6]int{local, utc} {
				for _, x := range fields {
					fmt.Fprint(&input, " ", x)
				}
			}
			fmt.Fprintln(&input)
			date, clock := int(c.Date), int(c.Clock)
			base := time.Date(1980+date/512, time.Month(date/32%16), date%32, 0, 0, 0, 0, time.UTC)
			decoded := base.Unix() + int64(clock/2048*3600+clock/32%64*60+clock%32*2)
			ld, lt := encode(local)
			ud, ut := encode(utc)
			cases[i].Want = fmt.Sprintf("%d|%d:%d|%d|%d:%d:%d:%d", decoded, ld, lt, decoded, c.Seconds, c.Nanos, ud, ut)
		}
		for i, row := range check.LeanRows("godiff-zip-dos-time", len(cases), &input) {
			c := cases[i]
			if row != c.Want {
				log.Fatalf("%s row %d: Lean=%q want=%q", label, i, row, c.Want)
			}
			must(enc.Encode(c))
		}
		fmt.Printf("%s: %d Lean rows agree\n", label, len(cases))
	}
	var fixed []fixture
	for _, d := range []uint16{0, 1, 32, 33, 63, 64, 0x1ff, 0x200, 0xffff} {
		for _, clock := range []uint16{0, 1, 31, 32, 0x7ff, 0x800, 0xffff} {
			for _, year := range []int{1899, 1900, 1969, 1979, 1980, 2000, 2099, 2100, 2107, 2108} {
				for _, offset := range []int{-43200, -901, 0, 901, 50400} {
					fixed = append(fixed, fixture{Date: d, Clock: clock, Seconds: time.Date(year, 3, 1, 0, 0, 1, 0, time.UTC).Unix(),
						Nanos: 999999999, Offset: offset})
				}
			}
		}
	}
	compare("boundaries", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			cases[i] = fixture{Date: uint16(r.Uint32()), Clock: uint16(r.Uint32()), Seconds: r.Int63n(31556952000) - 15778476000,
				Nanos: r.Intn(1_000_000_000), Offset: r.Intn(100801) - 50400}
		}
		compare(fmt.Sprintf("seed %d", seed), cases)
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_ZIP_DOS_ROWS", f.Name()))
	must(check.TestInPackage("archive/zip", "TestGoLeanDOSTime", "difftest/archive/zip/dos-time/testdata/dos_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
