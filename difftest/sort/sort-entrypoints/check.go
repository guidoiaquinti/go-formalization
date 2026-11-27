package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"math/rand"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	kind, tag, reported int
	values              []int
}
type tracedInts struct {
	values             []int
	reported, lenCalls int
	trace              []string
}

func (x *tracedInts) Len() int { x.lenCalls++; return x.reported }
func (x *tracedInts) Less(a, b int) bool {
	x.trace = append(x.trace, fmt.Sprintf("L%d,%d", a, b))
	return x.values[a] < x.values[b]
}
func (x *tracedInts) Swap(a, b int) {
	x.trace = append(x.trace, fmt.Sprintf("S%d,%d", a, b))
	x.values[a], x.values[b] = x.values[b], x.values[a]
}

type namedInts []int

func native(c fixture) (status string, x *tracedInts) {
	x = &tracedInts{values: slices.Clone(c.values), reported: c.reported}
	if c.values == nil {
		x.values = nil
	}
	status = "ok"
	defer func() {
		if p := recover(); p != nil {
			text := fmt.Sprint(p)
			switch {
			case strings.HasPrefix(text, "reflect: call of Swapper on zero Value"):
				status = "invalid"
			case strings.HasPrefix(text, "reflect: call of Swapper on "):
				status = "non-slice"
			case strings.Contains(text, "invalid memory address or nil pointer dereference"):
				if c.kind == 0 {
					status = "nil-interface"
				} else {
					status = "nil-less"
				}
			default:
				panic(p)
			}
		}
	}()
	switch c.kind {
	case 0:
		if c.tag == 0 {
			sort.Sort(x)
		} else {
			sort.Sort(nil)
		}
	case 1:
		padded := append([]int{311}, x.values...)
		padded = append(padded, -251)
		view := sort.IntSlice(padded[1 : len(padded)-1 : len(padded)])
		if c.tag == 1 {
			view = nil
		}
		view.Sort()
		if padded[0] != 311 || padded[len(padded)-1] != -251 {
			log.Fatal("IntSlice.Sort changed out-of-view storage")
		}
		x.values = []int(view)
	case 2:
		switch c.tag {
		case 0:
			sort.Slice(namedInts(x.values), x.Less)
		case 1:
			sort.Slice(x.values, nil)
		case 2:
			sort.Slice(nil, x.Less)
		case 3:
			sort.Slice([2]int{2, 1}, x.Less)
		case 4:
			sort.Slice("input", x.Less)
		}
	}
	return
}

func valuesWire(values []int) string {
	parts := make([]string, len(values))
	for i, n := range values {
		parts[i] = strconv.Itoa(n)
	}
	return strings.Join(parts, ",")
}

func observe(c fixture, model string) {
	status, data := native(c)
	wantStatus := "ok"
	wantLen := 0
	values := slices.Clone(c.values)
	switch c.kind {
	case 0:
		if c.tag == 1 {
			wantStatus = "nil-interface"
		} else {
			wantLen = 1
			if c.reported > 1 {
				slices.Sort(values[:c.reported])
			}
		}
	case 1:
		slices.Sort(values)
	case 2:
		switch c.tag {
		case 0:
			slices.Sort(values)
		case 1:
			if len(values) > 1 {
				wantStatus = "nil-less"
			}
		case 2:
			wantStatus = "invalid"
		default:
			wantStatus = "non-slice"
		}
	}
	if status != wantStatus || data.lenCalls != wantLen || !slices.Equal(data.values, values) {
		log.Fatalf("%+v: native %s, Len=%d, values=%v want %s, Len=%d, values=%v", c, status, data.lenCalls, data.values, wantStatus, wantLen, values)
	}
	if c.kind == 1 && (data.values == nil) != (c.values == nil) {
		log.Fatal("IntSlice.Sort changed modeled nilness")
	}
	if c.kind == 2 {
		fields := strings.SplitN(model, ":", 4)
		if len(fields) != 4 {
			log.Fatalf("invalid Lean row %q", model)
		}
		calls := strings.Split(fields[3], ";")
		calls = slices.DeleteFunc(calls, func(call string) bool { return strings.HasPrefix(call, "S") })
		fields[3] = strings.Join(calls, ";")
		model = strings.Join(fields, ":")
	}
	got := fmt.Sprintf("%s:%d:%s:%s", status, data.lenCalls, valuesWire(data.values), strings.Join(data.trace, ";"))
	if got != model {
		log.Fatalf("%+v: native %q, Lean %q", c, got, model)
	}
}

func compare(label string, cases []fixture) {
	var wire bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&wire, "%d %d %d %d", c.kind, strconv.IntSize, c.tag, c.reported)
		for _, v := range c.values {
			fmt.Fprintf(&wire, " %d", v)
		}
		fmt.Fprintln(&wire)
	}
	for i, got := range check.LeanRows("godiff-sort-entrypoints", len(cases), &wire) {
		observe(cases[i], got)
	}
	fmt.Printf("sort entry points %s: %d cases matched pin and Lean.\n", label, len(cases))
}

func boundaries(kind int) []fixture {
	var cases []fixture
	for _, width := range []int{0, 1, 2, 11, 12, 13, 49, 50, 51, 100, 101} {
		values := make([]int, width)
		for i := range values {
			values[i] = width - i
		}
		switch kind {
		case 0:
			for _, report := range []int{math.MinInt, -1, 0, min(1, width), width / 2, width} {
				cases = append(cases, fixture{kind, 0, report, values})
			}
			cases = append(cases, fixture{kind, 1, 0, values})
		case 1:
			cases = append(cases, fixture{kind, 0, 0, values})
		case 2:
			for tag := range 5 {
				cases = append(cases, fixture{kind, tag, 0, values})
			}
		}
	}
	if kind == 1 {
		cases = append(cases, fixture{kind, 1, 0, nil})
	}
	if kind == 2 {
		cases = append(cases, fixture{kind, 0, 0, nil}, fixture{kind, 1, 0, nil})
	}
	return cases
}

func main() {
	for kind := range 3 {
		compare(fmt.Sprintf("kind=%d boundaries", kind), boundaries(kind))
		for _, seed := range []int64{17, 2718, 65537} {
			rng := rand.New(rand.NewSource(seed))
			for block := range 4 {
				cases := make([]fixture, 5000)
				for i := range cases {
					width := rng.Intn(65)
					values := make([]int, width)
					for j := range values {
						switch i % 5 {
						case 0:
							values[j] = j
						case 1:
							values[j] = width - j
						case 2:
							values[j] = rng.Intn(4)
						default:
							values[j] = int(rng.Uint64())
						}
					}
					c := fixture{kind: kind, values: values, reported: width}
					if kind == 0 && i%7 == 0 {
						c.reported = rng.Intn(width+3) - 2
					}
					if kind == 0 && i%11 == 0 {
						c.tag = 1
					}
					if kind == 2 && i%5 == 0 {
						c.tag = 1 + rng.Intn(4)
					}
					cases[i] = c
				}
				compare(fmt.Sprintf("kind=%d seed=%d block=%d", kind, seed, block), cases)
			}
		}
	}
}
