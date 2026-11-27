package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct{ a, b string }

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dir, err := os.MkdirTemp("", "golean-strings-longest-common-suffix-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := os.Setenv("GOLEAN_STRINGS_LONGEST_COMMON_SUFFIX_CASES", path); err != nil {
		return err
	}
	fixture := bufio.NewWriter(file)

	patterns := []string{"", "aa", "\x00a", "\x80", "\xff\x00", "é", "ê", "🚀", "🤀",
		"\xff🚀", "\xe2\x82", "\xfd\xe2\x82", strings.Repeat("\xff\x00", 128)}
	var edges []input
	for _, a := range patterns {
		for _, b := range patterns {
			edges = append(edges, input{a, b})
		}
	}
	for i := range 1 << 16 {
		edges = append(edges, input{string([]byte{byte(i >> 8)}), string([]byte{byte(i)})})
	}
	if err := compare(fixture, "boundaries", edges); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]input, 20_000)
		for i := range cases {
			a, b := raw(rng, 80), raw(rng, 80)
			tail := raw(rng, 120)
			if i%4 == 0 {
				a = ""
			} else if i%4 == 1 {
				b = a
			}
			cases[i] = input{a + tail, b + tail}
		}
		if err := compare(fixture, fmt.Sprintf("seed=%d", seed), cases); err != nil {
			return err
		}
	}
	if err := fixture.Flush(); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	const fixtures = "difftest/strings/longest-common-suffix/testdata/"
	return check.TestInPackage("strings", "TestGoLeanLongestCommonSuffix",
		fixtures+"bridge_test.go", fixtures+"longest_common_suffix_test.go")
}

func raw(rng *rand.Rand, limit int) string {
	buf := make([]byte, rng.Intn(limit))
	rng.Read(buf)
	return string(buf)
}

func expected(in input) int {
	for n := min(len(in.a), len(in.b)); n > 0; n-- {
		if strings.HasSuffix(in.a, in.b[len(in.b)-n:]) {
			return n
		}
	}
	return 0
}

func compare(fixture *bufio.Writer, label string, cases []input) error {
	var modelInput bytes.Buffer
	wants := make([]string, len(cases))
	for i, in := range cases {
		row := strconv.AppendInt(nil, int64(len(in.a)), 10)
		for _, raw := range []string{in.a, in.b} {
			for j := range len(raw) {
				row = strconv.AppendUint(append(row, ' '), uint64(raw[j]), 10)
			}
		}
		fmt.Fprintln(&modelInput, string(row))
		wants[i] = strconv.Itoa(expected(in))
		if _, err := fmt.Fprintf(fixture, "%x\t%x\t%s\n", in.a, in.b, wants[i]); err != nil {
			return err
		}
	}
	for i, got := range check.LeanRows("godiff-strings-longest-common-suffix", len(cases), &modelInput) {
		if got != wants[i] {
			return fmt.Errorf("%s longestCommonSuffix(%x,%x): Lean=%s want=%s", label,
				cases[i].a, cases[i].b, got, wants[i])
		}
	}
	fmt.Printf("strings.longestCommonSuffix %s: %d Lean inputs agree.\n", label, len(cases))
	return nil
}
