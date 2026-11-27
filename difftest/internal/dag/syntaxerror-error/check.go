package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cases := [][]byte{nil, []byte("missing <"), []byte("%d %s"), []byte("a\x00\nb"), []byte("λ🐹"), {0xc0, 0xff}}
	for b := range 256 {
		cases = append(cases, []byte{byte(b)})
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for range 20_000 {
			input := make([]byte, rng.Intn(65))
			rng.Read(input)
			cases = append(cases, input)
		}
	}

	var stdin bytes.Buffer
	for _, input := range cases {
		fmt.Fprintln(&stdin, format("0", input))
	}
	rows := check.LeanRows("godiff-dag-syntaxerror-error", len(cases), &stdin)
	for i, input := range cases {
		if want := format("error", input); rows[i] != want {
			return fmt.Errorf("case %d input=%q: Lean=%q want=%q", i, input, rows[i], want)
		}
	}

	dir, err := os.MkdirTemp("", "golean-dag-error-method-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	data, err := json.Marshal(cases)
	if err != nil {
		return err
	}
	path := filepath.Join(dir, "inputs.json")
	if err := os.WriteFile(path, data, 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_DAG_ERROR_METHOD_INPUT", path); err != nil {
		return err
	}
	const fixture = "difftest/internal/dag/syntaxerror-error/testdata/error_method_test.go"
	if err := check.TestInPackage("internal/dag", "TestGoLeanSyntaxErrorMethod", fixture); err != nil {
		return err
	}
	fmt.Printf("internal/dag.syntaxError.Error: %d boundary and %d sampled cases agree, seeds %v.\n",
		boundaries, len(cases)-boundaries, seeds)
	return nil
}

func format(marker string, input []byte) string {
	row := []byte(marker)
	for _, b := range input {
		row = strconv.AppendUint(append(row, ' '), uint64(b), 10)
	}
	return string(row)
}
