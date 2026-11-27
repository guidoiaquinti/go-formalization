package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Count             int64
	Written, Existing int
	Failed            bool
	Data              []byte
	Want              string
}

func oracle(c fixture) string {
	counted := fmt.Sprintf("%d:%d:%t:1", int64(uint64(c.Count)+uint64(c.Written)), c.Written, c.Failed)
	remaining, called := 0, 0
	sticky := false
	if len(c.Data) > 0 {
		remaining, called = len(c.Data)-c.Written, 1
		sticky = c.Failed
	}
	capacity, preloaded, reused := 4096, 0, false
	if c.Existing >= 4096 {
		capacity, preloaded, reused = c.Existing, 1, true
	}
	return fmt.Sprintf("%s|%d:%t:%d:%d|%d:%t:0:%d|%d:%d:%t:0:0|false:0",
		counted, remaining, sticky, called, c.Count, remaining, sticky, c.Count, capacity, preloaded, reused)
}

func main() {
	f, err := os.CreateTemp("", "golean-zip-writer-io")
	must(err)
	defer os.Remove(f.Name())
	defer f.Close()
	w := bufio.NewWriter(f)
	enc := gob.NewEncoder(w)
	compare := func(label string, cases []fixture) {
		var input bytes.Buffer
		for _, c := range cases {
			failed := 0
			if c.Failed {
				failed = 1
			}
			fmt.Fprint(&input, c.Count, " ", c.Written, " ", failed, " ", c.Existing)
			for _, cell := range c.Data {
				fmt.Fprint(&input, " ", cell)
			}
			fmt.Fprintln(&input)
		}
		for i, got := range check.LeanRows("godiff-zip-writer-io", len(cases), &input) {
			c := cases[i]
			c.Want = oracle(c)
			if got != c.Want {
				log.Fatalf("%s row %d: Lean %q, oracle %q", label, i, got, c.Want)
			}
			must(enc.Encode(c))
		}
		fmt.Printf("%s: %d Lean writer observations agree\n", label, len(cases))
	}
	geometries := []int{-1, 0, 1, 4095, 4096, 4097, 8192}
	var fixed []fixture
	for _, count := range []int64{math.MinInt64, -1, 0, 1, math.MaxInt64 - 1, math.MaxInt64} {
		for _, capacity := range geometries {
			for _, data := range [][]byte{nil, {0}, {255, 0}, {65, 255, 0, 128, 65}} {
				fixed = append(fixed, fixture{Count: count, Written: len(data), Existing: capacity, Data: data})
				for n := 0; n <= len(data); n++ {
					fixed = append(fixed, fixture{Count: count, Written: n, Existing: capacity, Failed: true, Data: data})
				}
			}
		}
	}
	compare("fixed", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([]fixture, 5000)
			for i := range cases {
				data := make([]byte, r.Intn(33))
				_, err := r.Read(data)
				must(err)
				n, failed := len(data), r.Intn(2) == 0
				if failed {
					n = r.Intn(n + 1)
				}
				cases[i] = fixture{Count: int64(r.Uint64()), Written: n, Existing: geometries[r.Intn(len(geometries))], Failed: failed, Data: data}
			}
			compare(fmt.Sprintf("seed %d block %d", seed, block), cases)
		}
	}
	must(w.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_ZIP_WRITER_IO_ROWS", f.Name()))
	must(check.TestInPackage("archive/zip", "TestGoLeanZipWriterIO", "difftest/archive/zip/writer-io/testdata/writer_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
