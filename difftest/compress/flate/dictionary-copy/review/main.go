package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Backing                                []byte
	Length, Read, Written, Distance, Count int
	Full                                   bool
	Want                                   string
}

func main() {
	fixed := []fixture{
		{[]byte{10, 20, 30, 40, 50, 60, 70, 80, 91, 92}, 8, 1, 2, 4, 6, true,
			"6:1:8:true:10,20,70,80,10,20,70,80,91,92|0:1:2:true:10,20,30,40,50,60,70,80,91,92"},
		{[]byte{10, 20, 30, 40, 50, 60, 70, 80}, 8, 0, 2, 4, 1, true,
			"1:0:3:true:10,20,70,40,50,60,70,80|0:0:2:true:10,20,30,40,50,60,70,80"},
		{[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9}, 8, 2, 3, 2, 5, false,
			"5:2:8:false:1,2,3,2,3,2,3,2,9|5:2:8:false:1,2,3,2,3,2,3,2,9"},
		{[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9}, 8, 2, 5, 2, 8, false,
			"3:2:8:false:1,2,3,4,5,4,5,4,9|0:2:5:false:1,2,3,4,5,6,7,8,9"},
		{[]byte{19, 29, 39, 49}, 4, 0, 0, 4, 4, true,
			"4:0:4:true:19,29,39,49|0:0:0:true:19,29,39,49"},
		{[]byte{11, 14, 17, 20}, 4, 1, 2, 2, 0, false,
			"0:1:2:false:11,14,17,20|0:1:2:false:11,14,17,20"},
		{[]byte{11, 14, 17, 20, 23}, 4, 1, 2, 2, 2, false,
			"2:1:4:false:11,14,11,14,23|2:1:4:false:11,14,11,14,23"},
	}
	f, err := os.CreateTemp("", "golean-dict-copy-review")
	must(err)
	defer os.Remove(f.Name())
	enc := json.NewEncoder(f)
	var input bytes.Buffer
	var mutations [3]int
	for _, c := range fixed {
		full := 0
		if c.Full {
			full = 1
		}
		fmt.Fprint(&input, c.Length, " ", c.Read, " ", c.Written, " ", full, " ", c.Distance, " ", c.Count)
		for _, b := range c.Backing {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
		must(enc.Encode(c))
		ordinary, fast, _ := strings.Cut(c.Want, "|")
		if immutableCopy(c) != ordinary {
			mutations[0]++
		}
		if ordinary != fast {
			if c.Distance > c.Written {
				mutations[1]++
			}
			if c.Count > c.Length-c.Written {
				mutations[2]++
			}
		}
	}
	must(f.Close())
	for i, got := range check.LeanRows("godiff-flate-dictionary-copy", len(fixed), &input) {
		if got != fixed[i].Want {
			log.Fatalf("control %d: Lean %q, expected %q", i, got, fixed[i].Want)
		}
	}
	must(os.Setenv("GOLEAN_FLATE_DICT_COPY", f.Name()))
	must(check.TestInPackage("compress/flate", "TestGoLeanDictionaryCopy", "difftest/compress/flate/dictionary-copy/testdata/copy_test.go"))
	if expected := [3]int{3, 3, 1}; mutations != expected {
		log.Fatalf("mutation disagreements %v, want %v", mutations, expected)
	}
	fmt.Printf("%d literal states for both methods; mutation disagreements %v\n", len(fixed), mutations)
}

// immutableCopy never reads newly written bytes when the source overlaps output.
func immutableCopy(c fixture) string {
	backing := slices.Clone(c.Backing)
	count := min(c.Count, c.Length-c.Written)
	for i := range count {
		source := (c.Written + i + c.Length - c.Distance) % c.Length
		backing[c.Written+i] = c.Backing[source]
	}
	fields := make([]string, len(backing))
	for i, b := range backing {
		fields[i] = strconv.Itoa(int(b))
	}
	return fmt.Sprintf("%d:%d:%d:%t:%s", count, c.Read, c.Written+count, c.Full, strings.Join(fields, ","))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
