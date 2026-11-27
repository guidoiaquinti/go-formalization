package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"

	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Tag, Mode, Second int
	Payload, X, Y     int64
	Want              string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	file, err := os.CreateTemp("", "golean-sync-once")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()
	js := json.NewEncoder(file)
	compare := func(label string, cases []fixture) error {
		var rows bytes.Buffer
		for _, c := range cases {
			fmt.Fprintln(&rows, c.Tag, c.Mode, c.Second, c.Payload, c.X, c.Y)
		}
		for i, got := range check.LeanRows("godiff-sync-once", len(cases), &rows) {
			c := cases[i]
			c.Want = expected(c)
			if got != c.Want {
				return fmt.Errorf("row %+v Lean=%q want=%q", c, got, c.Want)
			}
			if err := js.Encode(c); err != nil {
				return err
			}
		}
		fmt.Printf("sync Once %s: %d Lean rows agree\n", label, len(cases))
		return nil
	}
	var edges []fixture
	for tag := range 3 {
		for mode := range 4 {
			for second := range 4 {
				for _, payload := range []int64{math.MinInt64, -1, 0, 1, math.MaxInt64} {
					for _, x := range []int64{math.MinInt64, 0, math.MaxInt64} {
						for _, y := range []int64{math.MinInt64, 0, math.MaxInt64} {
							edges = append(edges, fixture{Tag: tag, Mode: mode, Second: second, Payload: payload, X: x, Y: y})
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
		rows := make([]fixture, 20_000)
		for i := range rows {
			rows[i] = fixture{Tag: rng.Intn(3), Mode: rng.Intn(4), Second: rng.Intn(4),
				Payload: int64(rng.Uint64()), X: int64(rng.Uint64()), Y: int64(rng.Uint64())}
		}
		if err := compare(fmt.Sprintf("seed=%d", seed), rows); err != nil {
			return err
		}
	}
	if err := file.Close(); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_SYNC_ONCE_CASES", file.Name()); err != nil {
		return err
	}
	const overlay = "difftest/sync/once/testdata/"
	return check.TestInPackage("sync", "TestGoLeanSyncOnce", overlay+"bridge_test.go", overlay+"once_test.go")
}

func expected(c fixture) string {
	mode := []string{"nilCallback", "ok", fmt.Sprintf("panic:%d", c.Payload), "nilPanic"}[c.Mode]
	entered := c.Mode != 0
	once := "ok:false:true,ok:false:true"
	if c.Tag == 0 {
		once = "nilReceiver:false:nil,nilReceiver:false:nil"
	} else if c.Tag == 1 {
		once = fmt.Sprintf("%s:%t:true,ok:false:true", mode, entered)
	}
	cache := func(output string) string {
		if c.Mode == 1 {
			return output + ":true," + output + ":false"
		}
		return fmt.Sprintf("%s:%t,%s:false", mode, entered, mode)
	}
	return strings.Join([]string{once, once, cache("ok"), cache(fmt.Sprint(c.X)),
		cache(fmt.Sprintf("%d/%d", c.X, c.Y))}, "|")
}
