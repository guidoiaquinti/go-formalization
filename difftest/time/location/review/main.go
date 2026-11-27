package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Wall, Ext uint64
	Old, Next int
	Want      string
}

func main() {
	fixed := []fixture{
		{0, 0, 0, 0, "0,0,0|0,0,0|0,0,0|0,0,2|nil|1"},
		{1073741824, 13, 3, 4, "1073741824,13,3|1073741824,13,4|1073741824,13,0|1073741824,13,2|1073741824,13,4|3"},
		{9223372044370968581, 18446744073709551615, 4, 1, "5,59453308807,4|5,59453308807,0|5,59453308807,0|5,59453308807,2|5,59453308807,0|4"},
		{18446744073709551615, 0, 1, 0, "1073741823,68043243391,1|1073741823,68043243391,0|1073741823,68043243391,0|1073741823,68043243391,2|nil|1"},
		{9223372036854775808, 9223372036854775808, 2, 3, "0,59453308800,2|0,59453308800,3|0,59453308800,0|0,59453308800,2|0,59453308800,3|2"},
		{9223372036854775807, 18446744073709551615, 0, 2, "9223372036854775807,18446744073709551615,0|9223372036854775807,18446744073709551615,2|9223372036854775807,18446744073709551615,0|9223372036854775807,18446744073709551615,2|9223372036854775807,18446744073709551615,2|1"},
	}
	f, err := os.CreateTemp("", "golean-time-location-review")
	must(err)
	defer os.Remove(f.Name())
	enc := json.NewEncoder(f)
	var input bytes.Buffer
	for _, c := range fixed {
		fmt.Fprintln(&input, c.Wall, c.Ext, c.Old, c.Next)
		must(enc.Encode(c))
	}
	must(f.Close())
	for i, got := range check.LeanRows("godiff-time-location", len(fixed), &input) {
		if got != fixed[i].Want {
			log.Fatalf("fixture %d: Lean %q, expected %q", i, got, fixed[i].Want)
		}
	}
	must(os.Setenv("GOLEAN_TIME_LOCATION_ROWS", f.Name()))
	const bridge = "difftest/time/location/testdata/export_test.go"
	must(check.TestInPackage("time", "TestGoLeanLocation", bridge, "difftest/time/location/testdata/location_test.go"))
	must(check.TestInPackage("time", "TestGoLeanLocationIdentityControl", bridge, "difftest/time/location/review/testdata/identity_test.go"))
	fmt.Printf("%d literal rows and direct native identity controls passed\n", len(fixed))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
