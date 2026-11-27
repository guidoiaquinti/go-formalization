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

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Words                          []uint64
	A, B, Kind, Mode, Reply, Count int
	HasCounter                     bool
	Want                           string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	file, err := os.CreateTemp("", "golean-slices-order2-")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()
	encoder := json.NewEncoder(file)
	compare := func(label string, cases []fixture) error {
		var rows bytes.Buffer
		for _, c := range cases {
			present := 0
			if c.HasCounter {
				present = 1
			}
			fmt.Fprint(&rows, strconv.IntSize, " ", c.Kind, " ", c.A, " ", c.B, " ", present, " ", c.Count, " ", c.Mode, " ", c.Reply)
			for _, word := range c.Words {
				fmt.Fprint(&rows, " ", word)
			}
			fmt.Fprintln(&rows)
		}
		for i, got := range check.LeanRows("godiff-slices-order2", len(cases), &rows) {
			c := cases[i]
			c.Want = expected(c)
			if got != c.Want {
				return fmt.Errorf("%+v: Lean=%q", c, got)
			}
			if err := encoder.Encode(c); err != nil {
				return err
			}
		}
		fmt.Printf("order2 ordered and comparator %s: %d Lean rows agree.\n", label, len(cases))
		return nil
	}
	words := binary64.Boundaries()
	var edges []fixture
	for _, x := range words {
		for _, y := range words {
			for kind := range 2 {
				for mode := range 3 {
					for _, reply := range []int{math.MinInt, 0, math.MaxInt} {
						for _, hasCounter := range []bool{false, true} {
							edges = append(edges, fixture{Words: []uint64{x, y}, A: 0, B: 1, Kind: kind,
								Mode: mode, Reply: reply, Count: math.MaxInt, HasCounter: hasCounter})
						}
					}
				}
			}
		}
	}
	if err := compare("boundaries", edges); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			cases[i] = sample(rng, i)
		}
		if err := compare(fmt.Sprintf("seed=%d", seed), cases); err != nil {
			return err
		}
	}
	if err := file.Close(); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_SLICES_ORDER2_CASES", file.Name()); err != nil {
		return err
	}
	const data = "difftest/slices/order2/testdata/"
	return check.TestInPackage("slices", "TestGoLeanOrder2", data+"bridge_test.go", data+"order2_test.go")
}

func sample(rng *rand.Rand, i int) fixture {
	words := make([]uint64, 1+rng.Intn(10))
	for j := range words {
		words[j] = rng.Uint64()
	}
	a, b := rng.Intn(len(words)), rng.Intn(len(words))
	if i%7 == 0 {
		words[b] = words[a]
	}
	counter := int(rng.Uint64())
	if i%11 == 0 {
		counter = math.MaxInt
	}
	return fixture{Words: words, A: a, B: b, Kind: i % 2, Mode: rng.Intn(3), Reply: int(rng.Uint64()),
		Count: counter, HasCounter: i%4 != 0}
}

func expected(c fixture) string {
	left, right := c.Words[c.A], c.Words[c.B]
	inverted := int64(right) < int64(left)
	if c.Kind == 1 {
		x, y := math.Float64frombits(right), math.Float64frombits(left)
		inverted = x < y || math.IsNaN(x) && !math.IsNaN(y)
	}
	count := "nil"
	if c.HasCounter {
		count = strconv.Itoa(c.Count)
	}
	ordered := order(c, inverted, count)
	switch c.Mode {
	case 0:
		return ordered + "|nil/" + count + "/"
	case 2:
		return ordered + fmt.Sprintf("|panic:%d/%s/%d,%d", c.Reply, count, right, left)
	default:
		return ordered + "|" + order(c, c.Reply < 0, count) + fmt.Sprintf("%d,%d", right, left)
	}
}

func order(c fixture, inverted bool, count string) string {
	if !inverted {
		return fmt.Sprintf("ok:%d,%d/%s/", c.A, c.B, count)
	}
	if !c.HasCounter {
		return "nil/nil/"
	}
	return fmt.Sprintf("ok:%d,%d/%d/", c.B, c.A, c.Count+1)
}
