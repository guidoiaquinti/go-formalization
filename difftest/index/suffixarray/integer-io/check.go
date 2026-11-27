package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Length, Chunk, Terminal, Reply, Profile int
	Value                                   int
	Backing, Source                         []byte
	Integers                                []int64
	Want                                    string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	file, err := os.CreateTemp("", "golean-suffixarray-int-io")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()
	js := json.NewEncoder(file)
	compare := func(label string, cases []fixture) error {
		var rows bytes.Buffer
		for _, c := range cases {
			fmt.Fprint(&rows, strconv.IntSize, " ", c.Length, " ", len(c.Backing), " ", c.Value, " ", c.Chunk, " ", c.Terminal, " ", c.Reply, " ", len(c.Source), " ", c.Profile)
			for _, data := range [][]byte{c.Backing, c.Source} {
				for _, b := range data {
					fmt.Fprint(&rows, " ", b)
				}
			}
			for _, n := range c.Integers {
				fmt.Fprint(&rows, " ", n)
			}
			fmt.Fprintln(&rows)
		}
		for i, got := range check.LeanRows("godiff-suffixarray-integer-io", len(cases), &rows) {
			c := cases[i]
			c.Want = expected(c)
			if got != c.Want {
				return fmt.Errorf("row %+v Lean=%q want=%q", c, got, c.Want)
			}
			if err := js.Encode(c); err != nil {
				return err
			}
		}
		fmt.Printf("suffixarray integer IO %s: %d Lean rows agree\n", label, len(cases))
		return nil
	}
	var edges []fixture
	for _, size := range [][2]int{{0, 0}, {0, 10}, {1, 1}, {1, 12}, {9, 9}, {9, 10}, {10, 10}, {19, 20}, {20, 20}, {21, 40}, {29, 40}} {
		for _, value := range []int{math.MinInt, -65, -1, 0, 64, math.MaxInt} {
			for _, source := range [][]byte{nil, {14}, {0x80}, bytes.Repeat([]byte{0x80}, 10), bytes.Repeat([]byte{0xff}, 12)} {
				for terminal := range 3 {
					for reply := range 3 {
						for profile := range 2 {
							edges = append(edges, fixture{Length: size[0], Backing: bytes.Repeat([]byte{1}, size[1]),
								Value: value, Source: source, Chunk: 3, Terminal: terminal, Reply: reply, Profile: profile,
								Integers: []int64{1, -1, 256, math.MaxInt64}})
						}
					}
				}
			}
		}
	}
	edges = append(edges, fixture{Length: 1, Backing: bytes.Repeat([]byte{99}, 12), Value: -1, Chunk: 3, Reply: 1},
		fixture{Backing: bytes.Repeat([]byte{99}, 10), Source: []byte{14}, Chunk: 1},
		fixture{Length: 19, Backing: bytes.Repeat([]byte{99}, 20), Chunk: 1, Integers: []int64{1}},
		fixture{Length: 20, Backing: bytes.Repeat([]byte{99}, 20), Chunk: 1, Integers: []int64{1}})
	if err := compare("boundaries", edges); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			c := fixture{Backing: make([]byte, rng.Intn(60)), Source: make([]byte, rng.Intn(18)),
				Integers: make([]int64, rng.Intn(12)), Value: int(rng.Uint64()), Chunk: 1 + rng.Intn(15),
				Terminal: rng.Intn(3), Reply: rng.Intn(3), Profile: rng.Intn(2)}
			c.Length = rng.Intn(len(c.Backing) + 1)
			rng.Read(c.Backing)
			rng.Read(c.Source)
			for j := range c.Integers {
				c.Integers[j] = int64(rng.Uint64())
			}
			cases[i] = c
		}
		if err := compare(fmt.Sprintf("seed=%d", seed), cases); err != nil {
			return err
		}
	}
	if err := file.Close(); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_SUFFIXARRAY_INTEGER_IO_CASES", file.Name()); err != nil {
		return err
	}
	return check.TestInPackage("index/suffixarray", "TestGoLeanSuffixarrayIntegerIO", "difftest/index/suffixarray/integer-io/testdata/integer_io_test.go")
}

func expected(c fixture) string {
	write := func(pack bool) string {
		buf := slices.Clone(c.Backing)
		count, size := 0, 10
		var word int64
		if pack {
			for _, value := range c.Integers {
				if c.Length-size < 10 {
					break
				}
				if c.Profile == 0 {
					value = int64(int32(value))
				}
				size += copy(buf[size:], binary.AppendUvarint(nil, uint64(value)))
				count++
			}
			word = int64(size)
		} else {
			word = int64(c.Value)
		}
		header := binary.AppendVarint(nil, word)
		copy(buf[:c.Length], header)
		status, offer := "bounds", "none"
		if len(header) <= c.Length && size <= len(buf) {
			status = []string{"ok", "other", "nilIO"}[c.Reply]
			if c.Reply != 2 {
				offer = fmt.Sprintf("%d,%d:%s", size, len(buf), wire(buf[:size]))
			}
		}
		if offer == "none" {
			count = 0
		}
		return fmt.Sprintf("%s/%d/%s/%s", status, count, wire(buf), offer)
	}
	buf := slices.Clone(c.Backing)
	status, value := "bounds", int64(0)
	var needs []string
	if len(buf) >= 10 {
		status = "nilIO"
		if c.Terminal != 2 {
			for remaining, available := 10, len(c.Source); remaining > 0; {
				needs = append(needs, strconv.Itoa(remaining))
				used := min(remaining, available, c.Chunk)
				remaining -= used
				available -= used
				if available == 0 {
					break
				}
			}
			n := copy(buf[:10], c.Source)
			value, _ = binary.Varint(buf[:c.Length])
			switch {
			case n == 10:
				status = "ok"
			case c.Terminal == 1:
				status = "other"
			case n == 0:
				status = "eof"
			default:
				status = "unexpected"
			}
		}
	}
	return fmt.Sprintf("%s|%s/%d/%s/%s|%s", write(false), status, value, wire(buf), strings.Join(needs, ","), write(true))
}

func wire(data []byte) string {
	fields := make([]string, len(data))
	for i, b := range data {
		fields[i] = strconv.Itoa(int(b))
	}
	return strings.Join(fields, ",")
}
