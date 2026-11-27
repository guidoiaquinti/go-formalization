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
		return fmt.Sprintf("%d %d %d 0", c.r, c.g, c.b)
	case 2:
		return "panic 1"
	case 3:
		return "panic 0"
	}
	r := int64(c.r / 256 % 256)
	g := int64(c.g / 256 % 256)
	b := int64(c.b / 256 % 256)
	round := func(delta int64) int64 {
		return (delta+257*(1<<15))/(1<<16) - 128
	}
	y := g + round(19595*(r-g)+7471*(b-g))
	cb := min(255, max(0, 128+round(32768*(b-g)-11056*(r-g))))
	cr := min(255, max(0, 128+round(32768*(r-g)-5328*(b-g))))
	return fmt.Sprintf("%d %d %d 1", y, cb, cr)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cases := []input{{tag: 2}, {tag: 3}}
	for v := range uint32(256) {
		cases = append(cases, input{0, v, 255 - v, v ^ 128, 0})
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
				cases = append(cases, input{0, rng.Uint32() % 256, rng.Uint32() % 256, rng.Uint32() % 256, 0})
			} else {
				cases = append(cases, input{1, rng.Uint32(), rng.Uint32(), rng.Uint32(), rng.Uint32()})
			}
		}
	}
	var in bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&in, c.tag, c.r, c.g, c.b, c.a)
	}
	dir, err := os.MkdirTemp("", "golean-ycbcr-model-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	if err := os.WriteFile(filepath.Join(dir, "inputs"), in.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_YCBCR_MODEL_DIR", dir); err != nil {
		return err
	}
	lean := check.LeanRows("godiff-image-color-ycbcr-model", len(cases), &in)
	const fixture = "difftest/image/color/ycbcr-model/testdata/ycbcr_model_test.go"
	if err := check.TestInPackage("image/color", "TestGoLeanYCbCrModel", fixture); err != nil {
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
	fmt.Printf("image/color.yCbCrModel: %d boundary and 60000 sampled cases agree; seeds %v\n", boundaries, seeds)
	return nil
}
