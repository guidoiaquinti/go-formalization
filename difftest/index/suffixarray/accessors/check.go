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
	Tag, Len32, Len64, ByteLen, I, J int
	Value                            int64
	Words32, Words64                 []int64
	Bytes                            []byte
	Want                             string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	file, err := os.CreateTemp("", "golean-suffixarray-accessors")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()
	writer := json.NewEncoder(file)
	compare := func(label string, cases []fixture) error {
		var rows bytes.Buffer
		for _, c := range cases {
			fmt.Fprint(&rows, c.Tag, " ", c.Len32, " ", len(c.Words32), " ", c.Len64, " ", len(c.Words64), " ", c.ByteLen, " ", c.I, " ", c.J, " ", c.Value)
			for _, words := range [][]int64{c.Words32, c.Words64} {
				for _, word := range words {
					fmt.Fprint(&rows, " ", word)
				}
			}
			for _, b := range c.Bytes {
				fmt.Fprint(&rows, " ", b)
			}
			fmt.Fprintln(&rows)
		}
		for i, got := range check.LeanRows("godiff-suffixarray-accessors", len(cases), &rows) {
			c := cases[i]
			c.Want = expected(c)
			if got != c.Want {
				return fmt.Errorf("row %+v: Lean=%q", c, got)
			}
			if err := writer.Encode(c); err != nil {
				return err
			}
		}
		fmt.Printf("suffixarray accessors %s: %d Lean rows agree\n", label, len(cases))
		return nil
	}
	var boundaries []fixture
	values := []int64{math.MinInt64, -1, 0, 1, 2, 3, 1<<32 + 1, math.MaxInt64}
	for tag := range 16 {
		for _, length := range []int{0, 1, 3} {
			for _, capDelta := range []int{0, 1, 3} {
				for _, i := range []int{math.MinInt, -1, 0, 1, 2, 3, math.MaxInt} {
					for _, j := range []int{-1, 0, 1, 2, 3, 4, math.MaxInt} {
						for _, value := range values {
							boundaries = append(boundaries, makeFixture(tag, length, capDelta, i, j, value))
						}
					}
				}
			}
		}
	}
	if err := compare("boundaries", boundaries); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			length, delta := rng.Intn(8), rng.Intn(5)
			c := makeFixture(rng.Intn(16), length, delta, rng.Intn(length+delta+3)-1,
				rng.Intn(length+delta+3)-1, int64(rng.Uint64()))
			for j := range c.Words32 {
				c.Words32[j] = rng.Int63()
				if j%2 == 0 {
					c.Words32[j] = int64(rng.Intn(length+3) - 1)
				}
			}
			for j := range c.Words64 {
				c.Words64[j] = int64(rng.Uint64())
				if j%2 == 0 {
					c.Words64[j] = int64(rng.Intn(length+3) - 1)
				}
			}
			rng.Read(c.Bytes)
			cases[i] = c
		}
		if err := compare(fmt.Sprintf("seed=%d", seed), cases); err != nil {
			return err
		}
	}
	if err := file.Close(); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_SUFFIXARRAY_ACCESSORS_CASES", file.Name()); err != nil {
		return err
	}
	return check.TestInPackage("index/suffixarray", "TestGoLeanSuffixarrayAccessors",
		"difftest/index/suffixarray/accessors/testdata/accessors_test.go")
}

func makeFixture(tag, length, delta, i, j int, value int64) fixture {
	c := fixture{Tag: tag, I: i, J: j, Value: value}
	if tag&1 != 0 {
		c.Len32 = length
		c.Words32 = make([]int64, length+delta)
		for k := range c.Words32 {
			c.Words32[k] = int64(k) - 1
		}
	}
	if tag&2 != 0 {
		c.Len64 = length
		c.Words64 = make([]int64, length+delta)
		for k := range c.Words64 {
			c.Words64[k] = 0x100000000 + int64(k)
			if k%2 == 0 {
				c.Words64[k] = int64(k)
			}
		}
	}
	if tag&4 != 0 {
		c.ByteLen = length
		c.Bytes = make([]byte, length+delta)
		for k := range c.Bytes {
			c.Bytes[k] = byte(22 + k)
		}
	}
	return c
}

func expected(c fixture) string {
	if c.Tag&8 != 0 {
		return strings.TrimSuffix(strings.Repeat("nilReceiver|", 6), "|")
	}
	x := make([]int64, len(c.Words32))
	for i, v := range c.Words32 {
		x[i] = int64(int32(v))
	}
	if c.Words32 == nil {
		x = nil
	}
	y := c.Words64
	selected, length := y, c.Len64
	prefer32 := x != nil
	if prefer32 {
		selected, length = x, c.Len32
	}
	get, store, sub, at := "bounds", "bounds", "bounds", "bounds"
	if c.I >= 0 && c.I < length {
		start := selected[c.I]
		get = fmt.Sprintf("ok:%d", start)
		if start >= 0 && start <= int64(c.ByteLen) {
			at = "ok:" + wire(c.ByteLen-int(start), c.Bytes[start:])
		}
		if prefer32 {
			x := append([]int64(nil), x...)
			x[c.I] = int64(int32(c.Value))
			store = "ok:" + wire(c.Len32, x) + "/" + wire(c.Len64, y)
		} else {
			y := append([]int64(nil), y...)
			y[c.I] = c.Value
			store = "ok:" + wire(c.Len32, x) + "/" + wire(c.Len64, y)
		}
	}
	if c.I >= 0 && c.I <= c.J && c.J <= len(selected) {
		image := wire(c.J-c.I, selected[c.I:])
		sub = "ok:nil/" + image
		if prefer32 {
			sub = "ok:" + image + "/nil"
		}
	}
	return fmt.Sprintf("ok:%d|%s|%s|%s|ok:%s|%s", c.Len32+c.Len64, get, store, sub, wire(c.ByteLen, c.Bytes), at)
}

func wire[T ~byte | ~int64](length int, backing []T) string {
	if backing == nil {
		return "nil"
	}
	fields := make([]string, len(backing))
	for i, v := range backing {
		fields[i] = strconv.FormatInt(int64(v), 10)
	}
	return fmt.Sprintf("%d:%s", length, strings.Join(fields, ","))
}
