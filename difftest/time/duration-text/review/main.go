package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	fixed := []struct {
		value   int64
		present bool
		text    string
	}{
		{math.MinInt64, true, "-2562047h47m16.854775808s"},
		{math.MaxInt64, true, "2562047h47m16.854775807s"},
		{0, true, "0s"},
		{-999, true, "-999ns"},
		{1000, true, "1\xc2\xb5s"},
		{-1010, true, "-1.01\xc2\xb5s"},
		{999999, true, "999.999\xc2\xb5s"},
		{1002000, true, "1.002ms"},
		{999999999, true, "999.999999ms"},
		{1020000000, true, "1.02s"},
		{60000000001, true, "1m0.000000001s"},
		{3600000000000, true, "1h0m0s"},
		{-1001, false, "-1.001\xc2\xb5s"},
		{0, false, "0s"},
	}
	f, err := os.CreateTemp("", "golean-duration-review")
	must(err)
	defer os.Remove(f.Name())
	enc := json.NewEncoder(f)
	var input bytes.Buffer
	wants := make([]string, 0, len(fixed))
	var mutations [3]int
	for i, c := range fixed {
		var cells [32]byte
		for j := range cells {
			cells[j] = byte(0x80 + i + j)
		}
		present := 0
		if c.present {
			present = 1
		}
		fmt.Fprint(&input, c.value, " ", present)
		for _, b := range cells {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
		if got := time.Duration(c.value).String(); got != c.text {
			log.Fatalf("public row %d: native %q want %q", i, got, c.text)
		}
		if strings.ReplaceAll(c.text, "\xc2\xb5", "\xce\xbc") != c.text {
			mutations[0]++
		}
		if _, fraction, ok := strings.Cut(c.text, "."); ok {
			precision, unit := 9, "s"
			if strings.HasSuffix(fraction, "µs") {
				precision, unit = 3, "µs"
			} else if strings.HasSuffix(fraction, "ms") {
				precision, unit = 6, "ms"
			}
			digits := strings.TrimSuffix(fraction, unit)
			if len(digits) < precision {
				mutations[1]++
			}
		}
		private := "nil"
		if c.present {
			start := len(cells) - len(c.text)
			next := cells
			copy(next[start:], c.text)
			private = fmt.Sprintf("%d:%s", start, wire(next[:]))
			cleared := next
			clear(cleared[:start])
			if cleared != next {
				mutations[2]++
			}
		}
		wants = append(wants, wire([]byte(c.text))+"|"+private)
		must(enc.Encode(struct {
			Value   int64
			Present bool
			Cells   [32]byte
			Want    string
		}{c.value, c.present, cells, private}))
	}
	must(f.Close())
	for i, got := range check.LeanRows("godiff-time-duration-text", len(fixed), &input) {
		if got != wants[i] {
			log.Fatalf("fixture %d: Lean %q, want %q", i, got, wants[i])
		}
	}
	must(os.Setenv("GOLEAN_DURATION_TEXT_ROWS", f.Name()))
	must(check.TestInPackage("time", "TestGoLeanDurationText",
		"difftest/time/duration-text/testdata/export_test.go",
		"difftest/time/duration-text/testdata/duration_text_test.go"))
	if want := [3]int{4, 3, 12}; mutations != want {
		log.Fatalf("mutation disagreements %v, want %v", mutations, want)
	}
	fmt.Printf("%d literal public and full-buffer controls; mutation disagreements %v\n", len(fixed), mutations)
}

func wire(raw []byte) string {
	cells := make([]string, len(raw))
	for i, b := range raw {
		cells[i] = strconv.Itoa(int(b))
	}
	return strings.Join(cells, ",")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
