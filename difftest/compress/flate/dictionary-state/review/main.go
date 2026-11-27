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
	Backing                      []byte
	Length, Read, Written, Count int
	Full                         bool
	Value                        byte
	Want                         string
}

func main() {
	fixed := []fixture{
		{nil, 0, 0, 0, 0, false, 1, "0|0|0|false:0:0:|0:0:false:|bounds|false:0:0:|0:0:true:"},
		{[]byte{}, 0, 0, 0, 0, false, 1, "0|0|0|true:0:0:|0:0:false:|bounds|true:0:0:|0:0:true:"},
		{[]byte{17, 26}, 0, 0, 0, 0, false, 1, "0|0|0|true:2:0:|0:0:false:17,26|bounds|true:2:0:|0:0:true:17,26"},
		{[]byte{11, 22, 33, 44, 55}, 3, 1, 2, 1, false, 99, "2|1|1|true:3:2:33|1:3:false:11,22,33,44,55|1:3:false:11,22,99,44,55|true:4:1:22|2:2:false:11,22,33,44,55"},
		{[]byte{11, 22, 33, 44, 55}, 3, 1, 3, 0, false, 99, "3|2|0|true:2:3:|1:3:false:11,22,33,44,55|bounds|true:4:1:22,33|0:0:true:11,22,33,44,55"},
		{[]byte{4, 5, 6}, 3, 1, 1, 2, true, 255, "3|0|2|true:2:1:5,6|1:3:true:4,5,6|1:2:true:4,255,6|true:2:1:|1:1:true:4,5,6"},
		{[]byte{31, 211}, 2, 2, 2, 0, true, 17, "2|0|0|true:0:2:|2:2:true:31,211|bounds|true:0:2:|0:0:true:31,211"},
	}
	f, err := os.CreateTemp("", "golean-dict-review")
	must(err)
	defer os.Remove(f.Name())
	enc := json.NewEncoder(f)
	var input bytes.Buffer
	var mutations [3]int
	for _, c := range fixed {
		present, full := 0, 0
		if c.Backing != nil {
			present = 1
		}
		if c.Full {
			full = 1
		}
		fmt.Fprint(&input, present, " ", c.Length, " ", c.Read, " ", c.Written, " ", full, " ", c.Count, " ", c.Value)
		for _, v := range c.Backing {
			fmt.Fprint(&input, " ", v)
		}
		fmt.Fprintln(&input)
		must(enc.Encode(c))
		fields := strings.Split(c.Want, "|")
		if clippedCapacity := strconv.Itoa(c.Length - c.Written); clippedCapacity != strings.Split(fields[3], ":")[1] {
			mutations[0]++
		}
		if wrapAtCapacity := c.Written == len(c.Backing); wrapAtCapacity != strings.HasPrefix(fields[7], "0:0:true:") {
			mutations[1]++
		}
		if storeAtCapacity := c.Written < len(c.Backing); storeAtCapacity != (fields[5] != "bounds") {
			mutations[2]++
		}
	}
	must(f.Close())
	for i, got := range check.LeanRows("godiff-flate-dictionary-state", len(fixed), &input) {
		if got != fixed[i].Want {
			log.Fatalf("control %d: Lean %q, expected %q", i, got, fixed[i].Want)
		}
	}
	must(os.Setenv("GOLEAN_FLATE_DICT_STATE", f.Name()))
	must(check.TestInPackage("compress/flate", "TestGoLeanDictState", "difftest/compress/flate/dictionary-state/testdata/dict_test.go"))
	must(check.TestInPackage("compress/flate", "TestGoLeanDictReview", "difftest/compress/flate/dictionary-state/review/testdata/dict_review_test.go"))
	if expected := [3]int{3, 2, 2}; mutations != expected {
		log.Fatalf("mutation disagreements %v, want %v", mutations, expected)
	}
	fmt.Printf("%d literal Go/Lean controls; mutation disagreements %v; native alias, atomicity and wrap controls pass\n", len(fixed), mutations)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
