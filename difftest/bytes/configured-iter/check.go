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

type input struct {
	NonNil              bool
	Raw, Spare, Sep     []byte
	Save, First, Second int
}

type fixture struct {
	Input input
	Want  string
}

func expected(c input, lines bool) string {
	text, tail, separator := string(c.Raw), string(c.Spare), string(c.Sep)
	nonNil := c.NonNil
	part := func(text string) string {
		tag := "nil"
		if nonNil {
			tag = "slice"
		}
		fields := make([]string, len(text))
		for i := range fields {
			fields[i] = strconv.Itoa(int(text[i]))
		}
		return fmt.Sprintf("%s:%d:%s", tag, len(text), strings.Join(fields, ","))
	}
	run := func(budget int) string {
		var parts []string
		for len(parts) < budget {
			if lines {
				if text == "" {
					break
				}
				line, rest, found := strings.Cut(text, "\n")
				if found {
					parts = append(parts, part(line+"\n"))
					text = rest
				} else {
					parts = append(parts, part(line))
					text, tail, nonNil = "", "", false
				}
				continue
			}
			if separator == "" {
				if text == "" {
					break
				}
				width := len(text)
				for index := range text {
					if index > 0 {
						width = index
						break
					}
				}
				parts = append(parts, part(text[:width]))
				if len(parts) < budget {
					text = text[width:]
				}
				continue
			}
			before, rest, found := strings.Cut(text, separator)
			if !found {
				parts = append(parts, part(text))
				break
			}
			// Check retention against the available extent before adding.
			if c.Save < -len(before) || c.Save > len(text)+len(tail)-len(before) {
				return "panic/" + strings.Join(parts, ";")
			}
			parts = append(parts, part((text + tail)[:len(before)+c.Save]))
			if len(parts) < budget {
				text = rest
			}
		}
		return "ok/" + strings.Join(parts, ";")
	}
	return strings.Join([]string{run(c.First), run(c.Second), run(4096), run(4096)}, "|")
}

func batch(label string, cases []input, all *[]fixture) {
	var protocol bytes.Buffer
	rows := make([]fixture, len(cases))
	for i, c := range cases {
		tag := 0
		if c.NonNil {
			tag = 1
		}
		fmt.Fprintf(&protocol, "%d %d %d %d %d %d %d", tag, len(c.Raw), len(c.Spare), len(c.Sep), c.Save, c.First, c.Second)
		for _, section := range [][]byte{c.Raw, c.Spare, c.Sep} {
			for _, b := range section {
				fmt.Fprintf(&protocol, " %d", b)
			}
		}
		fmt.Fprintln(&protocol)
		rows[i] = fixture{c, expected(c, true) + "#" + expected(c, false)}
	}
	for i, got := range check.LeanRows("godiff-bytes-configured-iter", len(cases), &protocol) {
		if got != rows[i].Want {
			log.Fatalf("%s row %d (%+v): got %q want %q", label, i, cases[i], got, rows[i].Want)
		}
	}
	*all = append(*all, rows...)
	fmt.Printf("Configured byte iterators %s: %d Lean inputs agree.\n", label, len(cases))
}

func main() {
	var edges []input
	for _, raw := range []string{"", "\n", "\n\n", "a\r\nb\nc", "a--b--", "ab--", "a\x00b\x00", "\xff\xef\xbf\xbd\xe2\x82", "a--"} {
		for _, sep := range []string{"", "\n", "-", "--", "\x00", "\xff"} {
			for _, save := range []int{math.MinInt, -2, -1, 0, 1, len(sep), len(raw) + 1, len(raw) + 3, math.MaxInt} {
				for first := range 4 {
					edges = append(edges, input{true, []byte(raw), []byte{0xee, 0xaf}, []byte(sep), save, first, 3 - first})
				}
			}
		}
	}
	for first := range 4 {
		edges = append(edges, input{First: first, Second: 1}, input{Sep: []byte{10}, First: first, Second: 1})
	}
	var all []fixture
	batch("boundaries", edges, &all)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]input, 20000)
		for i := range cases {
			raw, spare, sep := make([]byte, rng.Intn(33)), make([]byte, rng.Intn(9)), make([]byte, rng.Intn(5))
			rng.Read(raw)
			rng.Read(spare)
			rng.Read(sep)
			if i%2 == 0 {
				for j := range raw {
					raw[j] = []byte{'a', 'b', 10, 13, 0xff}[rng.Intn(5)]
				}
			}
			if i%3 == 0 && len(raw) > 0 {
				a := rng.Intn(len(raw))
				sep = raw[a : a+rng.Intn(len(raw)-a+1)]
			}
			cases[i] = input{true, raw, spare, sep, rng.Intn(2*len(raw)+16) - len(raw) - 8, rng.Intn(6), rng.Intn(6)}
		}
		batch(fmt.Sprintf("seed=%d", seed), cases, &all)
	}
	file, err := os.CreateTemp("", "golean-configured-iter-")
	must(err)
	defer os.Remove(file.Name())
	must(json.NewEncoder(file).Encode(all))
	must(file.Close())
	must(os.Setenv("GOLEAN_CONFIGURED_ITER_CASES", file.Name()))
	const data = "difftest/bytes/configured-iter/testdata/"
	must(check.TestInPackage("bytes", "TestGoLeanConfiguredIter", data+"bridge_test.go", data+"configured_iter_test.go"))
	fmt.Printf("Configured byte iterators pinned Go overlay: all %d inputs agree.\n", len(all))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
