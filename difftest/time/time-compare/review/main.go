package main

import (
	"bytes"
	"cmp"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	// Each outcome lists nsec, sec, unixSec, mono, Nanosecond, Unix,
	// IsZero, Before, After, Equal, and Compare in that order.
	fixed := []struct{ input, want string }{
		{"0 0 0 0", "0 0 -62135596800 0 0 -62135596800 true false false true 0"},
		{"1073741824 0 0 0", "0 0 -62135596800 0 0 -62135596800 false false false true 0"},
		{"1073741823 -9223372036854775808 0 0", "1073741823 -9223372036854775808 9223371974719179008 0 1073741823 9223371974719179008 false true false false -1"},
		{"18446744073709551615 -1 9223372036854775808 0", "1073741823 68043243391 5907646591 -1 1073741823 5907646591 false true false false -1"},
		{"9223372036854775808 -1 18446744073709551615 -1", "0 59453308800 -2682288000 -1 0 -2682288000 false false false true 0"},
		{"9223372036854775810 -100 1 59453308800", "2 59453308800 -2682288000 -100 2 -2682288000 false false true false 1"},
		{"2 59453308800 9223372036854775810 -100", "2 59453308800 -2682288000 0 2 -2682288000 false false false true 0"},
		{"0 62135596800 1 62135596800", "0 62135596800 0 0 0 0 false true false false -1"},
		{"0 9223372036854775807 1073741823 9223372036854775806", "0 9223372036854775807 9223371974719179007 0 0 9223371974719179007 false false true false 1"},
	}
	f, err := os.CreateTemp("", "golean-time-compare-review")
	must(err)
	defer os.Remove(f.Name())
	var input bytes.Buffer
	for _, c := range fixed {
		fmt.Fprintln(&input, c.input)
		_, err := fmt.Fprintln(f, c.input, c.want)
		must(err)
	}
	must(f.Close())
	for i, got := range check.LeanRows("godiff-time-compare", len(fixed), &input) {
		if got != fixed[i].want {
			log.Fatalf("review row %d: Lean returned %q, expected %q", i, got, fixed[i].want)
		}
	}
	must(os.Setenv("GOLEAN_CLOCK_COMPARE", f.Name()))
	must(check.TestInPackage("time", "TestGoLeanClockCompare",
		"difftest/time/time-compare/testdata/bridge_test.go",
		"difftest/time/time-compare/testdata/compare_test.go"))
	var mutations [4]int
	for _, c := range fixed {
		var tw, uw uint64
		var te, ue int64
		_, err := fmt.Sscan(c.input, &tw, &te, &uw, &ue)
		must(err)
		fields := strings.Fields(c.want)
		order, err := strconv.Atoi(fields[10])
		must(err)
		const flag = uint64(1 << 63)
		const nanosField = uint64(1 << 30)
		tMono, uMono := tw >= flag, uw >= flag
		if decodedZero := !tMono && te == 0 && tw%nanosField == 0; decodedZero != (fields[6] == "true") {
			mutations[0]++
		}
		if tMono && uMono {
			if cmp.Compare(uint64(te), uint64(ue)) != order {
				mutations[1]++
			}
			// Packed wall timestamps order by the entire word when both flags are set.
			if cmp.Compare(tw, uw) != order {
				mutations[2]++
			}
		}
		if tMono != uMono && cmp.Compare(te, ue) != order {
			mutations[3]++
		}
	}
	if want := [4]int{1, 1, 2, 2}; mutations != want {
		log.Fatalf("mutation disagreements %v, want %v", mutations, want)
	}
	fmt.Printf("%d literal native/Lean controls; mutation disagreements %v\n", len(fixed), mutations)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
