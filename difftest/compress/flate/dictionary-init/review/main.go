package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Old, Preset                 []byte
	Length, Read, Written, Size int
	Full                        bool
	Want                        string
}

func main() {
	fixed := []fixture{
		{nil, []byte{99}, 0, 0, 0, 0, false, "false|0|0|0|0|true|true||"},
		{[]byte{}, []byte{99}, 0, 0, 0, 0, false, "true|0|0|0|0|true|true||"},
		{[]byte{12, 13}, []byte{99}, 1, 1, 1, 0, false, "true|0|2|0|0|true|true|12,13|12,13"},
		{[]byte{11, 22, 33, 44, 55}, []byte{61, 62, 63, 64, 65}, 1, 1, 1, 3, false, "true|3|5|0|0|true|true|63,64,65,44,55|63,64,65,44,55"},
		{[]byte{71, 72, 73, 74}, []byte{0, 255}, 2, 1, 2, 4, true, "true|4|4|2|2|false|true|0,255,73,74|0,255,73,74"},
		{[]byte{11, 22, 33}, []byte{231, 232}, 2, 1, 2, 5, true, "true|5|5|2|2|false|false|231,232,0,0,0|11,22,33"},
		{nil, []byte{81, 82, 83, 84}, 0, 0, 0, 2, false, "true|2|2|0|0|true|false|83,84|"},
		{[]byte{91, 92, 93}, []byte{11}, 3, 1, 3, 1, false, "true|1|3|0|0|true|true|11,92,93|11,92,93"},
		{[]byte{101, 102, 103}, nil, 3, 1, 3, 2, true, "true|2|3|0|0|false|true|101,102,103|101,102,103"},
	}
	f, err := os.CreateTemp("", "golean-dict-init-review")
	must(err)
	defer os.Remove(f.Name())
	enc := json.NewEncoder(f)
	var input bytes.Buffer
	var mutations [3]int
	for _, c := range fixed {
		present, full := 0, 0
		if c.Old != nil {
			present = 1
		}
		if c.Full {
			full = 1
		}
		fmt.Fprint(&input, present, " ", c.Length, " ", len(c.Old), " ", c.Read, " ", c.Written, " ", full, " ", c.Size)
		for _, raw := range [][]byte{c.Old, c.Preset} {
			for _, b := range raw {
				fmt.Fprint(&input, " ", b)
			}
		}
		fmt.Fprintln(&input)
		must(enc.Encode(c))

		fields := strings.Split(c.Want, "|")
		if c.Size > 0 && len(c.Preset) > c.Size &&
			strconv.Itoa(int(c.Preset[0])) != strings.Split(fields[7], ",")[0] {
			mutations[0]++
		}
		if reuseByLength := c.Length >= c.Size; reuseByLength != (fields[6] == "true") {
			mutations[1]++
		}
		if nonzeroFullOnly := c.Size != 0 && len(c.Preset) >= c.Size; nonzeroFullOnly != (fields[5] == "true") {
			mutations[2]++
		}
	}
	must(f.Close())
	for i, got := range check.LeanRows("godiff-flate-dictionary-init", len(fixed), &input) {
		if got != fixed[i].Want {
			log.Fatalf("control %d: Lean %q, expected %q", i, got, fixed[i].Want)
		}
	}
	must(os.Setenv("GOLEAN_FLATE_DICT_INIT", f.Name()))
	must(check.TestInPackage("compress/flate", "TestGoLeanDictionaryInit", "difftest/compress/flate/dictionary-init/testdata/init_test.go"))
	if expected := [3]int{2, 2, 3}; mutations != expected {
		log.Fatalf("mutation disagreements %v, want %v", mutations, expected)
	}
	fmt.Printf("%d literal dictionary initialization controls; mutation disagreements %v\n", len(fixed), mutations)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
