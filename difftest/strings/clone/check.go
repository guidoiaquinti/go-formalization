package main

import (
	"log"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"unsafe"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	backing    string
	start, end int
}

var sink string

func main() {
	bytes := []string{"", "ab", "\u00e9\U0001f680", "\x00x\x00", "\xff\xc0\x80", strings.Repeat("a", 128)}
	for b := range 256 {
		bytes = append(bytes, string([]byte{byte(b)}))
	}
	var edges []input
	for _, s := range bytes {
		edges = append(edges, input{s, 0, len(s)}, input{"pre" + s + "post", 3, 3 + len(s)})
	}
	check.CompareSamples("godiff-strings-clone", edges, sample, func(c input) string {
		return describe("0", c.backing[c.start:c.end])
	}, actual, func(c input) string {
		return describe(strconv.FormatBool(c.start != c.end), c.backing[c.start:c.end])
	})
	nonempty := runtime.GOOS + "/source"
	empty := nonempty[1:1]
	if n := testing.AllocsPerRun(1000, func() { sink = strings.Clone(empty) }); n != 0 {
		log.Fatalf("cloning an empty substring allocated %g objects per call", n)
	}
}

func sample(rng *rand.Rand, _ int) input {
	raw := make([]byte, rng.Intn(129))
	rng.Read(raw)
	start := rng.Intn(len(raw) + 1)
	end := start + rng.Intn(len(raw)-start+1)
	return input{string(raw), start, end}
}

func actual(c input) string {
	cloned := strings.Clone(c.backing[c.start:c.end])
	if len(cloned) == 0 {
		return describe("false", cloned)
	}
	// Compare the result with the whole source allocation represented by backing.
	// Keep both strings live while comparing their address ranges.
	src := uintptr(unsafe.Pointer(unsafe.StringData(c.backing)))
	dst := uintptr(unsafe.Pointer(unsafe.StringData(cloned)))
	disjoint := dst+uintptr(len(cloned)) <= src || src+uintptr(len(c.backing)) <= dst
	runtime.KeepAlive(c.backing)
	runtime.KeepAlive(cloned)
	return describe(strconv.FormatBool(disjoint), cloned)
}

func describe(prefix, s string) string {
	row := []byte(prefix)
	for i := range len(s) {
		row = strconv.AppendUint(append(row, ' '), uint64(s[i]), 10)
	}
	return string(row)
}
