package main

import (
	"bytes"
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
	Length, Count int
	Backing       []int64
	Want          string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	file, err := os.CreateTemp("", "golean-suffixarray-map-lms")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()
	js := json.NewEncoder(file)
	compare := func(label string, cases []fixture) error {
		var rows bytes.Buffer
		for _, c := range cases {
			fmt.Fprint(&rows, c.Length, " ", c.Count)
			for _, value := range c.Backing {
				fmt.Fprint(&rows, " ", value)
			}
			fmt.Fprintln(&rows)
		}
		for i, got := range check.LeanRows("godiff-suffixarray-map-lms", len(cases), &rows) {
			c := cases[i]
			c.Want = expected(c)
			if got != c.Want {
				return fmt.Errorf("row %+v Lean=%q want=%q", c, got, c.Want)
			}
			if err := js.Encode(c); err != nil {
				return err
			}
		}
		fmt.Printf("suffixarray LMS map %s: %d Lean rows agree\n", label, len(cases))
		return nil
	}
	var edges []fixture
	for n := range 7 {
		for pattern := range 1 << (2 * n) {
			backing := make([]int64, n+2)
			for i := range n {
				backing[i] = []int64{0, -1, 1, 1<<32 + 1}[(pattern>>(2*i))&3]
			}
			backing[n], backing[n+1] = 99, 100
			for _, count := range []int{math.MinInt, -1, 0, n, math.MaxInt} {
				edges = append(edges, fixture{Length: n, Count: count, Backing: backing})
			}
		}
	}
	edges = append(edges, fixture{Length: 6, Count: -8, Backing: []int64{4, 0, 5, 8, 19, 21, 99}},
		fixture{Length: 6, Count: 3, Backing: []int64{4, 0, 5, 8, 19, 21, 99}}, fixture{})
	if err := compare("boundaries", edges); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			n := rng.Intn(70)
			c := fixture{Length: n, Count: int(rng.Uint64()), Backing: make([]int64, n+rng.Intn(6))}
			for j := range c.Backing {
				switch i % 3 {
				case 0:
					c.Backing[j] = int64(rng.Intn(7) - 3)
				case 1:
					c.Backing[j] = int64(rng.Uint64())
				case 2:
					c.Backing[j] = int64(rng.Intn(99) + 1)
				}
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
	if err := os.Setenv("GOLEAN_SUFFIXARRAY_MAP_LMS_CASES", file.Name()); err != nil {
		return err
	}
	return check.TestInPackage("index/suffixarray", "TestGoLeanSuffixarrayMapLms", "difftest/index/suffixarray/map-lms/testdata/map_test.go")
}

func expected(c fixture) string {
	image := func(narrow bool) string {
		backing := slices.Clone(c.Backing)
		if narrow {
			for i, v := range backing {
				backing[i] = int64(int32(v))
			}
		}
		status := "bounds"
		if c.Length > 0 {
			status = "ok"
			var packed []int64
			for _, v := range backing[:c.Length/2+1] {
				if v > 0 {
					packed = append(packed, v-1)
				}
			}
			copy(backing[c.Length-len(packed):c.Length], packed)
		}
		words := make([]string, len(backing))
		for i, v := range backing {
			words[i] = strconv.FormatInt(v, 10)
		}
		return status + ":" + strings.Join(words, ",")
	}
	return image(true) + "|" + image(false)
}
