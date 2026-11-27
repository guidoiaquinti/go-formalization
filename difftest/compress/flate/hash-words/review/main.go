package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	const mixed = 0xdeadbeeffedcba98
	cases := []struct {
		word           uint64
		bits, count    uint8
		four, variable uint32
	}{
		{0xffffffff, 32, 4, 1640531535, 1640531535},
		{mixed, 32, 3, 1903614744, 652883968},
		{mixed, 33, 3, 0, 0},
		{mixed, 0, 255, 0, 0},
		{mixed, 40, 8, 0, 1704895668},
		{mixed, 64, 7, 0, 136366080},
		{mixed, 65, 5, 0, 0},
		{mixed, 32, 6, 1903614744, 1597372796},
		{mixed, 24, 0, 7435995, 7435995},
		{mixed, 24, 4, 7435995, 7435995},
	}
	var wire []string
	var native bytes.Buffer
	var input bytes.Buffer
	var counts [3]int
	for _, c := range cases {
		fmt.Fprintf(&input, "%d %d %d\n", c.word, c.bits, c.count)
		fmt.Fprintf(&native, "%d %d %d %d %d\n", c.word, c.bits, c.count, c.four, c.variable)
		wire = append(wire, fmt.Sprintf("%d %d", c.four, c.variable))
		// Saturating a negative shift reports the product instead of the Go overshift.
		var satShift uint8
		if c.bits <= 32 {
			satShift = 32 - c.bits
		}
		if uint32(c.word)*2654435761>>satShift != c.four {
			counts[0]++
		}
		if c.count == 0 || c.count == 255 {
			if c.variable != 0 {
				counts[1]++
			}
		}
		if c.count >= 5 && c.count <= 8 {
			shift := []uint{24, 16, 8, 0}[c.count-5]
			prime := []uint64{889523592379, 227718039650203, 58295818150454627, 0xcf1bbcdcb7a56463}[c.count-5]
			wrong := uint32((((c.word << shift) * prime) >> (64 - c.bits)) >> 32)
			if wrong != c.variable {
				counts[2]++
			}
		}
	}
	if want := [3]int{4, 1, 3}; counts != want {
		log.Fatalf("mutants %v, want %v", counts, want)
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-flate-hash-words")
	cmd.Stdin = &input
	got, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	if actual := strings.TrimSuffix(string(got), "\n"); actual != strings.Join(wire, "\n") {
		log.Fatalf("Lean observations differ: %s", actual)
	}
	tmp, err := os.CreateTemp("", "golean-flate-hash-review-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(native.Bytes()); err != nil {
		log.Fatal(err)
	}
	if err := tmp.Close(); err != nil {
		log.Fatal(err)
	}
	if err := os.Setenv("GOLEAN_FLATE_HASH_ROWS", tmp.Name()); err != nil {
		log.Fatal(err)
	}
	if err := check.TestInPackage("compress/flate", "TestGoLeanHashWords", "difftest/compress/flate/hash-words/testdata/hash_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("flate hash review: %s literal rows, wrong rules rejected %v\n", strconv.Itoa(len(cases)), counts)
}
