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
	width, tag int
	channels   [4]uint32
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var cases []input
	for _, width := range []int{8, 16} {
		cases = append(cases, input{width: width, tag: 2}, input{width: width, tag: 3})
		for value := range 1 << width {
			c := uint32(value)
			mask := uint32(1<<width - 1)
			cases = append(cases, input{width, 0, [4]uint32{c, c ^ mask, (c + 1) & mask, (c + 3) & mask}})
		}
		for _, r := range []uint32{0, 255, 256, 65535, 65536, 0x8000ff00, 0xffffffff} {
			for _, g := range []uint32{0, 255, 256, 65535, 65536, 0x8000ff00, 0xffffffff} {
				cases = append(cases, input{width, 1, [4]uint32{r, g, ^r, ^g}})
			}
		}
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for range 20_000 {
			words := [4]uint32{rng.Uint32(), rng.Uint32(), rng.Uint32(), rng.Uint32()}
			cases = append(cases, input{8, 1, words}, input{16, 1, words})
		}
	}
	var in bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&in, c.width, c.tag, c.channels[0], c.channels[1], c.channels[2], c.channels[3])
	}
	dir, err := os.MkdirTemp("", "golean-rgba-model-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	if err := os.WriteFile(filepath.Join(dir, "inputs"), in.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_RGBAMODEL_DIR", dir); err != nil {
		return err
	}
	lean := check.LeanRows("godiff-image-color-rgba-model", len(cases), &in)
	if err := check.TestInPackage("image/color", "TestGoLeanRGBAModel", "difftest/image/color/rgba-model/testdata/rgba_model_test.go"); err != nil {
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
		case 0, 1:
			words := c.channels
			if c.tag == 1 {
				for j, word := range words {
					if c.width == 8 {
						words[j] = word / 256 % 256
					} else {
						words[j] = word % 65536
					}
				}
			}
			want = fmt.Sprintf("%d %d %d %d %d", words[0], words[1], words[2], words[3], c.tag)
		case 2:
			want = "panic 1"
		case 3:
			want = "panic 0"
		}
		if rows[i] != want || lean[i] != want {
			return fmt.Errorf("case %+v: Go=%q Lean=%q want=%q", c, rows[i], lean[i], want)
		}
	}
	fmt.Printf("image/color.rgbaModel and rgba64Model: %d boundaries and 120000 samples agree, seeds %v\n", boundaries, seeds)
	return nil
}
