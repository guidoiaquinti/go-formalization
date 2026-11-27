package main

import (
	"bufio"
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
	Kind, A, B, C, Initial int
	Script                 []bool
	Want                   string
}

func main() {
	if strconv.IntSize != 64 {
		fmt.Println("Pivot helper comparison covers signed 64-bit Go ints.")
		return
	}
	file, err := os.CreateTemp("", "golean-sort-pivot-")
	must(err)
	defer os.Remove(file.Name())
	buffer := bufio.NewWriter(file)
	rows := json.NewEncoder(buffer)
	for kind := range 4 {
		compare(fmt.Sprintf("helper=%d boundaries", kind), boundaries(kind), rows)
		for _, seed := range []int64{17, 2718, 65537} {
			rng := rand.New(rand.NewSource(seed + int64(kind)))
			for block := range 4 {
				cases := make([]fixture, 5000)
				for i := range cases {
					cases[i] = sample(rng, kind, i)
				}
				compare(fmt.Sprintf("helper=%d seed=%d block=%d", kind, seed+int64(kind), block), cases, rows)
			}
		}
	}
	must(buffer.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_SORT_PIVOT_CASES", file.Name()))
	const data = "difftest/sort/pivot-helpers/testdata/"
	must(check.TestInPackage("sort", "TestGoLeanPivotHelpers", data+"bridge_test.go", data+"pivot_helpers_test.go"))
}

func replies(bits, width int) []bool {
	script := make([]bool, width)
	for i := range script {
		script[i] = bits&(1<<i) != 0
	}
	return script
}

func boundaries(kind int) []fixture {
	var cases []fixture
	if kind < 3 {
		width := 3
		if kind == 0 {
			width = 1
		}
		for _, indices := range [][3]int{{1, 2, 3}, {3, 2, 1}, {2, 2, 2}, {20, 11, 20}, {math.MaxInt - 1, 0, math.MaxInt}} {
			for _, initial := range []int{0, 1, -1, math.MinInt, math.MaxInt - 1, math.MaxInt} {
				for bits := range 1 << width {
					cases = append(cases, fixture{kind, indices[0], indices[1], indices[2], initial, replies(bits, width), ""})
				}
			}
		}
		return cases
	}
	for _, width := range []int{0, 1, 3, 4, 7, 8, 9, 11, 12, 48, 49, 50, 51, 63, 64, 65} {
		for _, a := range []int{0, 13, math.MaxInt - width} {
			for bits := range 8 {
				cases = append(cases, fixture{kind, a, a + width, 0, 0, replies(bits, 3), ""})
			}
		}
	}
	for bits := range 1 << 12 {
		cases = append(cases, fixture{kind, 9, 59, 0, 0, replies(bits, 12), ""})
	}
	return cases
}

func sample(rng *rand.Rand, kind, index int) fixture {
	a, b, c := int(rng.Uint64()>>1), int(rng.Uint64()>>1), int(rng.Uint64()>>1)
	initial := int(rng.Uint64())
	if index%9 == 0 {
		initial = math.MaxInt
	}
	if kind == 2 {
		a = max(1, min(a, math.MaxInt-1))
	}
	if kind == 3 {
		width := rng.Intn(1025)
		if index%4 == 0 {
			width = []int{7, 8, 49, 50}[rng.Intn(4)]
		}
		a = min(a, math.MaxInt-width)
		b = a + width
	}
	return fixture{kind, a, b, c, initial, replies(rng.Intn(1<<12), 12), ""}
}

// Apply an index sorting network to the three positions. The registered
// swap count is separate from the data and no Swap callback is allowed.
func expected(c fixture) string {
	var trace []string
	cost := 0
	exchange := func(items []int, a, b int) {
		i, j := items[a], items[b]
		trace = append(trace, fmt.Sprintf("%d,%d", j, i))
		step := len(trace) - 1
		if step < len(c.Script) && c.Script[step] {
			items[a], items[b] = j, i
			cost++
		}
	}
	middle := func(a, b, c int) int {
		items := []int{a, b, c}
		for _, edge := range [][2]int{{0, 1}, {1, 2}, {0, 1}} {
			exchange(items, edge[0], edge[1])
		}
		return items[1]
	}
	first, second := 0, -1
	switch c.Kind {
	case 0:
		items := []int{c.A, c.B}
		exchange(items, 0, 1)
		first, second = items[0], items[1]
	case 1:
		first = middle(c.A, c.B, c.C)
	case 2:
		first = middle(c.A-1, c.A, c.A+1)
	default:
		width := c.B - c.A
		centers := [3]int{}
		for i := range centers {
			centers[i] = c.A + (i+1)*(width>>2)
			if width >= 50 {
				center := centers[i]
				centers[i] = middle(center-1, center, center+1)
			}
		}
		first = centers[1]
		if width >= 8 {
			first = middle(centers[0], centers[1], centers[2])
		}
		second = 0
		switch cost {
		case 0:
			second = 1
		case 12:
			second = 2
		}
		return fmt.Sprintf("%d:%d:0:%s", first, second, strings.Join(trace, ";"))
	}
	return fmt.Sprintf("%d:%d:%d:%s", first, second, int(uint(c.Initial)+uint(cost)), strings.Join(trace, ";"))
}

func compare(label string, cases []fixture, rows *json.Encoder) {
	var wire bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&wire, "%d %d %d %d %d", c.Kind, c.A, c.B, c.C, c.Initial)
		for _, reply := range c.Script {
			if reply {
				wire.WriteString(" 1")
			} else {
				wire.WriteString(" 0")
			}
		}
		fmt.Fprintln(&wire)
	}
	for i, got := range check.LeanRows("godiff-sort-pivot-helpers", len(cases), &wire) {
		c := cases[i]
		c.Want = expected(c)
		if got != c.Want {
			log.Fatalf("%s row %d %+v: Lean=%s want=%s", label, i, c, got, c.Want)
		}
		must(rows.Encode(c))
	}
	fmt.Printf("sort pivot helpers %s: %d Lean rows agree.\n", label, len(cases))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
