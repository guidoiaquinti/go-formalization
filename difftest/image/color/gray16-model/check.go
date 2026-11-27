package main

import (
	"bytes"
	"encoding/binary"
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

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cases := []input{{tag: 2}, {tag: 3}}
	for y := range 65536 {
		cases = append(cases, input{tag: 0, r: uint32(y)})
	}
	bounds := []uint32{0, 1, 127, 128, 255, 256, 32767, 32768, 65535, ^uint32(0)}
	for _, r := range bounds {
		for _, g := range bounds {
			for _, b := range bounds {
				cases = append(cases, input{1, r, g, b, r ^ g ^ b})
			}
		}
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			r, g, b, a := rng.Uint32(), rng.Uint32(), rng.Uint32(), rng.Uint32()
			if i%2 == 0 {
				r, g, b, a = r&65535, g&65535, b&65535, a&65535
			}
			cases = append(cases, input{1, r, g, b, a})
		}
	}
	var in bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&in, c.tag, c.r, c.g, c.b, c.a)
	}
	dir, err := os.MkdirTemp("", "golean-gray16-model-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	if err := os.WriteFile(filepath.Join(dir, "inputs"), in.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_GRAY16MODEL_DIR", dir); err != nil {
		return err
	}
	lean := check.LeanRows("godiff-image-color-gray16-model", len(cases), &in)
	if err := check.TestInPackage("image/color", "TestGoLeanGray16Model", "difftest/image/color/gray16-model/testdata/gray16_model_test.go"); err != nil {
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
			want = fmt.Sprintf("%d 0", c.r)
		case 1:
			sum := 19595*uint64(c.r) + 38470*uint64(c.g) + 7471*uint64(c.b) + 32768
			var word [4]byte
			binary.BigEndian.PutUint32(word[:], uint32(sum))
			want = fmt.Sprintf("%d 1", binary.BigEndian.Uint16(word[:2]))
		case 2:
			want = "panic 1"
		case 3:
			want = "panic 0"
		}
		if rows[i] != want || lean[i] != want {
			return fmt.Errorf("case %+v: Go=%q Lean=%q want=%q", c, rows[i], lean[i], want)
		}
	}
	fmt.Printf("image/color.gray16Model: %d boundaries and 60000 samples agree, seeds %v\n", boundaries, seeds)
	return nil
}
