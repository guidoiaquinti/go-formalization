package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	mode, argument, state, delta uint64
}

func (c input) expected() string {
	if c.mode < 2 {
		return fmt.Sprintf("nil-panic 0 0 0 %d", c.state)
	}
	state := c.state + c.delta + c.argument
	tag, result := "ok", c.argument
	switch c.mode {
	case 3:
		result = 0
	case 4:
		result = 4294967297
	case 5:
		tag, result = "callback-panic", c.delta
	}
	return fmt.Sprintf("%s %d 1 %d %d", tag, result, c.argument, state)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var cases []input
	for mode := range uint64(6) {
		for _, argument := range []uint64{0, 1, 1 << 32} {
			for _, state := range []uint64{0, 1<<32 - 1} {
				for _, delta := range []uint64{0, 1, 1<<32 - 1} {
					cases = append(cases, input{mode, argument, state, delta})
				}
			}
		}
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for i := range 20000 {
			argument := uint64(rng.Uint32()) + 1
			if i%3 == 0 {
				argument = 0
			}
			cases = append(cases, input{uint64(rng.Intn(6)), argument, uint64(rng.Uint32()), uint64(rng.Uint32())})
		}
	}
	var in bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&in, c.mode, c.argument, c.state, c.delta)
	}
	dir, err := os.MkdirTemp("", "golean-model-func-convert-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	if err := os.WriteFile(filepath.Join(dir, "inputs"), in.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MODEL_FUNC_CONVERT_DIR", dir); err != nil {
		return err
	}
	lean := check.LeanRows("godiff-image-color-model-func-convert", len(cases), &in)
	const fixture = "difftest/image/color/model-func-convert/testdata/model_func_convert_test.go"
	if err := check.TestInPackage("image/color", "TestGoLeanModelFuncConvert", fixture); err != nil {
		return err
	}
	out, err := os.ReadFile(filepath.Join(dir, "outputs"))
	if err != nil {
		return err
	}
	actual := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	if len(actual) != len(cases) {
		return fmt.Errorf("Go produced %d rows, want %d", len(actual), len(cases))
	}
	for i, c := range cases {
		want := c.expected()
		if actual[i] != want || lean[i] != want {
			return fmt.Errorf("case %+v: Go=%q Lean=%q want=%q", c, actual[i], lean[i], want)
		}
	}
	fmt.Printf("image/color.modelFunc.Convert: %d boundary and 60000 sampled cases agree; seeds %v\n", boundaries, seeds)
	return nil
}
