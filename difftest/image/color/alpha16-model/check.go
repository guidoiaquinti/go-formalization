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
	tag   int
	value uint32
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cases := []input{{2, 0}, {3, 0}}
	for a := range 65536 {
		cases = append(cases, input{0, uint32(a)}, input{1, uint32(a)})
	}
	for _, high := range []uint32{1, 0x8000, 0xffff} {
		for _, low := range []uint32{0, 1, 0x7fff, 0x8000, 0xffff} {
			cases = append(cases, input{1, high<<16 | low})
		}
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for range 20_000 {
			cases = append(cases, input{1, rng.Uint32()})
		}
	}
	var in bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&in, c.tag, c.value)
	}
	dir, err := os.MkdirTemp("", "golean-alpha16-model-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	if err := os.WriteFile(filepath.Join(dir, "inputs"), in.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_ALPHA16MODEL_DIR", dir); err != nil {
		return err
	}
	lean := check.LeanRows("godiff-image-color-alpha16-model", len(cases), &in)
	if err := check.TestInPackage("image/color", "TestGoLeanAlpha16Model", "difftest/image/color/alpha16-model/testdata/alpha16_model_test.go"); err != nil {
		return err
	}
	data, err := os.ReadFile(filepath.Join(dir, "outputs"))
	if err != nil {
		return err
	}
	rows := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
	if len(rows) != len(cases) {
		return fmt.Errorf("Go produced %d rows, want %d", len(rows), len(cases))
	}
	for i, c := range cases {
		var want string
		switch c.tag {
		case 0:
			want = fmt.Sprintf("%d 0", c.value)
		case 1:
			want = fmt.Sprintf("%d 1", c.value%65536)
		case 2:
			want = "panic 1"
		case 3:
			want = "panic 0"
		}
		if rows[i] != want || lean[i] != want {
			return fmt.Errorf("case %+v: Go=%q Lean=%q want=%q", c, rows[i], lean[i], want)
		}
	}
	fmt.Printf("image/color.alpha16Model: %d boundaries and 60000 samples agree, seeds %v\n", boundaries, seeds)
	return nil
}
