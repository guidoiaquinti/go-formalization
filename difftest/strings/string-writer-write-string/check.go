package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	text           string
	count, failure int
	bounded        bool
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dir, err := os.MkdirTemp("", "golean-strings-string-writer-write-string-")
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
	if err := os.Setenv("GOLEAN_STRINGS_STRING_WRITER_WRITE_STRING_CASES", path); err != nil {
		return err
	}
	fixture := bufio.NewWriter(file)

	inputs := []string{"", "one", "\xff\x00\x80", "\xe2\x82", "�", "é🚀", strings.Repeat("\x00\xfc", 128)}
	for b := range 256 {
		inputs = append(inputs, string([]byte{byte(b)}))
	}
	var edges []input
	for _, payload := range inputs {
		for failure := range 3 {
			for _, count := range []int{0, 1, len(payload), len(payload) + 1, -1, math.MinInt, math.MaxInt} {
				edges = append(edges, input{payload, count, failure, false})
			}
			for _, limit := range []int{0, 1, 256} {
				edges = append(edges, input{payload, limit, failure, true})
			}
		}
	}
	if err := compare(fixture, "boundaries", edges); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]input, 20_000)
		for i := range cases {
			raw := make([]byte, rng.Intn(40))
			rng.Read(raw)
			count := int(rng.Uint64())
			if i%3 != 0 {
				count = rng.Intn(50)
			}
			cases[i] = input{string(raw), count, rng.Intn(3), i%3 != 0}
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
	const fixtures = "difftest/strings/string-writer-write-string/testdata/"
	return check.TestInPackage("strings", "TestGoLeanStringWriterWriteString",
		fixtures+"bridge_test.go", fixtures+"string_writer_write_string_test.go")
}

func compare(fixture *bufio.Writer, label string, cases []input) error {
	var modelInput bytes.Buffer
	wants := make([]string, len(cases))
	for i, c := range cases {
		mode := 0
		returned := c.count
		if c.bounded {
			mode = 1
			returned = min(returned, len(c.text))
		}
		row := fmt.Appendf(nil, "%d %d %d", mode, c.count, c.failure)
		want := fmt.Appendf(nil, "%d:%d:B", returned, c.failure)
		for j := range len(c.text) {
			value := uint64(c.text[j])
			row = strconv.AppendUint(append(row, ' '), value, 10)
			want = strconv.AppendUint(append(want, ','), value, 10)
		}
		fmt.Fprintln(&modelInput, string(row))
		wants[i] = string(want)
		if _, err := fmt.Fprintf(fixture, "%d\t%d\t%d\t%x\t%s\n", mode, c.count, c.failure, c.text, want); err != nil {
			return err
		}
	}
	for i, got := range check.LeanRows("godiff-strings-string-writer-write-string", len(cases), &modelInput) {
		if got != wants[i] {
			c := cases[i]
			return fmt.Errorf("%s WriteString(%x) %v: Lean=%s want=%s", label, c.text, c, got, wants[i])
		}
	}
	fmt.Printf("strings.stringWriter.WriteString %s: %d Lean traces agree.\n", label, len(cases))
	return nil
}
