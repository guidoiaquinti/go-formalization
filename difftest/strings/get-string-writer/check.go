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

type reply struct{ limit, failure int }

type input struct {
	text                   string
	mode                   int
	byteReply, stringReply reply
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dir, err := os.MkdirTemp("", "golean-strings-get-string-writer-")
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
	if err := os.Setenv("GOLEAN_STRINGS_GET_STRING_WRITER_CASES", path); err != nil {
		return err
	}
	fixture := bufio.NewWriter(file)
	payloads := []string{"", "word", "\xff\x00\x80", "\xe2\x82", "�", "é🚀", strings.Repeat("\x00\xfc", 128)}
	for b := range 256 {
		payloads = append(payloads, string([]byte{byte(b)}))
	}
	replies := []reply{{0, 0}, {1, 0}, {1, 1}, {256, 2}}
	var edges []input
	for _, text := range payloads {
		for _, byteReply := range replies {
			for _, stringReply := range replies {
				for mode := range 2 {
					edges = append(edges, input{text, mode, byteReply, stringReply})
				}
			}
		}
	}
	for mode := 2; mode <= 4; mode++ {
		edges = append(edges, input{mode: mode})
	}
	if err := compare(fixture, "boundaries", edges); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]input, 20_000)
		for i := range cases {
			payload := make([]byte, rng.Intn(40))
			rng.Read(payload)
			cases[i] = input{string(payload), i % 2,
				reply{rng.Intn(50), rng.Intn(3)}, reply{rng.Intn(50), rng.Intn(3)}}
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
	const fixtures = "difftest/strings/get-string-writer/testdata/"
	return check.TestInPackage("strings", "TestGoLeanGetStringWriter",
		fixtures+"bridge_test.go", fixtures+"get_string_writer_test.go")
}

func expected(c input) string {
	if c.mode > 1 {
		if c.mode == 2 {
			return "1"
		}
		return "0"
	}
	selected, method := c.byteReply, 'B'
	if c.mode == 1 {
		selected, method = c.stringReply, 'S'
	}
	result := fmt.Appendf(nil, "%d:%d:%d:%c", c.mode, min(selected.limit, len(c.text)), selected.failure, method)
	for i := range len(c.text) {
		result = strconv.AppendUint(append(result, ','), uint64(c.text[i]), 10)
	}
	return string(result)
}

func compare(fixture *bufio.Writer, label string, cases []input) error {
	var modelInput bytes.Buffer
	wants := make([]string, len(cases))
	for i, c := range cases {
		row := fmt.Appendf(nil, "%d %d %d %d %d", c.mode,
			c.byteReply.limit, c.byteReply.failure, c.stringReply.limit, c.stringReply.failure)
		for j := range len(c.text) {
			row = strconv.AppendUint(append(row, ' '), uint64(c.text[j]), 10)
		}
		fmt.Fprintln(&modelInput, string(row))
		wants[i] = expected(c)
		if _, err := fmt.Fprintf(fixture, "%d\t%d\t%d\t%d\t%d\t%x\t%s\n", c.mode,
			c.byteReply.limit, c.byteReply.failure, c.stringReply.limit, c.stringReply.failure,
			c.text, wants[i]); err != nil {
			return err
		}
	}
	for i, got := range check.LeanRows("godiff-strings-get-string-writer", len(cases), &modelInput) {
		if got != wants[i] {
			return fmt.Errorf("%s getStringWriter input %+v: Lean=%s want=%s", label, cases[i], got, wants[i])
		}
	}
	fmt.Printf("strings.getStringWriter %s: %d Lean selections agree.\n", label, len(cases))
	return nil
}
