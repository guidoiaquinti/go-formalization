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
	tag        int
	r, g, b, a uint32
}

func (c input) expected() string {
	switch c.tag {
	case 0:
		return fmt.Sprintf("%d %d %d %d 0", c.r, c.g, c.b, c.a)
	case 2:
		return "panic 1"
	case 3:
		return "panic 0"
	}
	r, g, b := c.r/256%256, c.g/256%256, c.b/256%256
	white := max(r, g, b)
	if white == 0 {
		return "0 0 0 255 1"
	}
	ink := func(v uint32) uint32 {
		return 255 - (255*v+white-1)/white
	}
	return fmt.Sprintf("%d %d %d %d 1", ink(r), ink(g), ink(b), 255-white)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cases := []input{{tag: 2}, {tag: 3}}
	for v := range uint32(256) {
		cases = append(cases, input{0, v, 255 - v, v ^ 128, v ^ 64})
	}
	edges := []uint32{0, 1, 255, 256, 0x8000, 0xffff, 0x10000, 0xffffffff}
	for _, r := range edges {
		for _, g := range edges {
			for _, b := range edges {
				for _, a := range []uint32{0, 0xffff} {
					cases = append(cases, input{1, r, g, b, a})
				}
			}
		}
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for i := range 20000 {
			if i%5 == 0 {
				cases = append(cases, input{0, rng.Uint32() % 256, rng.Uint32() % 256, rng.Uint32() % 256, rng.Uint32() % 256})
			} else {
				cases = append(cases, input{1, rng.Uint32(), rng.Uint32(), rng.Uint32(), rng.Uint32()})
			}
		}
	}
	var in bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&in, c.tag, c.r, c.g, c.b, c.a)
	}
	dir, err := os.MkdirTemp("", "golean-cmyk-model-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	if err := os.WriteFile(filepath.Join(dir, "inputs"), in.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_CMYK_MODEL_DIR", dir); err != nil {
		return err
	}
	lean := check.LeanRows("godiff-image-color-cmyk-model", len(cases), &in)
	const fixture = "difftest/image/color/cmyk-model/testdata/cmyk_model_test.go"
	if err := check.TestInPackage("image/color", "TestGoLeanCMYKModel", fixture); err != nil {
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
	fmt.Printf("image/color.cmykModel: %d boundary and 60000 sampled cases agree; seeds %v\n", boundaries, seeds)
	return nil
}
