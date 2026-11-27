package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Tag, Repeat   int
	Count, Offset int64
	Fill          byte
	Old, Tail     []byte
	Want          string
}

func expected(c fixture) string {
	state := func(status string, count int64, comment []byte) string {
		if c.Tag == 0 {
			return status + ":none"
		}
		countText := "none"
		if c.Tag == 2 {
			countText = strconv.FormatInt(count, 10)
		}
		out := make([]string, len(comment))
		for i, v := range comment {
			out[i] = strconv.Itoa(int(v))
		}
		return status + ":" + countText + ":" + strings.Join(out, ",")
	}
	var comment string
	switch {
	case c.Repeat+len(c.Tail) > math.MaxUint16:
		comment = state("tooLong", c.Count, c.Old)
	case c.Tag == 0:
		comment = state("nil", c.Count, c.Old)
	default:
		comment = state("ok", c.Count, append(bytes.Repeat([]byte{c.Fill}, c.Repeat), c.Tail...))
	}
	var offset string
	switch {
	case c.Tag != 2:
		offset = state("nil", c.Count, c.Old)
	case c.Count != 0:
		offset = state("afterWrite", c.Count, c.Old)
	default:
		offset = state("ok", c.Offset, c.Old)
	}
	return comment + "|" + offset
}

func main() {
	f, err := os.CreateTemp("", "golean-zip-writer-config")
	must(err)
	defer os.Remove(f.Name())
	defer f.Close()
	enc := json.NewEncoder(f)
	compare := func(label string, cases []fixture) {
		var input bytes.Buffer
		for _, c := range cases {
			fmt.Fprint(&input, c.Tag, " ", c.Count, " ", c.Offset, " ", c.Repeat, " ", c.Fill, " ", len(c.Old))
			for _, s := range [][]byte{c.Old, c.Tail} {
				for _, b := range s {
					fmt.Fprint(&input, " ", b)
				}
			}
			fmt.Fprintln(&input)
		}
		for i, row := range check.LeanRows("godiff-zip-writer-config", len(cases), &input) {
			c := cases[i]
			c.Want = expected(c)
			if row != c.Want {
				log.Fatalf("%s row %d: incorrect writer state", label, i)
			}
			must(enc.Encode(c))
		}
		fmt.Printf("%s: %d Lean rows agree\n", label, len(cases))
	}
	var fixed []fixture
	for tag := range 3 {
		for _, size := range []int{0, 65533, 65534, 65535, 65536, 65537} {
			for _, count := range []int64{0, -1, 1, math.MaxInt64} {
				for _, offset := range []int64{0, -1, 1, math.MinInt64, math.MaxInt64} {
					fixed = append(fixed, fixture{Tag: tag, Repeat: size, Count: count, Offset: offset, Fill: 255, Old: []byte{1, 0, 2}, Tail: []byte{0, 128}})
				}
			}
		}
	}
	compare("boundaries", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			c := fixture{Tag: r.Intn(3), Repeat: r.Intn(64), Count: int64(r.Uint64()), Offset: int64(r.Uint64()), Fill: byte(r.Uint32()),
				Old: make([]byte, r.Intn(30)), Tail: make([]byte, r.Intn(30))}
			if i%4 == 0 {
				c.Count = 0
			}
			r.Read(c.Old)
			r.Read(c.Tail)
			cases[i] = c
		}
		compare(fmt.Sprintf("seed %d", seed), cases)
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_ZIP_WRITER_CONFIG", f.Name()))
	must(check.TestInPackage("archive/zip", "TestGoLeanZipWriterConfig", "difftest/archive/zip/writer-config/testdata/config_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
