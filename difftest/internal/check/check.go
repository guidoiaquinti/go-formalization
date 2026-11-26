// Package check runs compiled Lean models from the repository root.
package check

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func LeanRows(binary string, count int, input io.Reader, args ...string) []string {
	pin, err := os.ReadFile(".go-version")
	if err != nil {
		log.Fatal(err)
	}
	if want := "go" + strings.TrimSpace(string(pin)); runtime.Version() != want {
		log.Fatalf("expected %s, built with %s", want, runtime.Version())
	}
	cmd := exec.Command(filepath.Join("lean", ".lake", "build", "bin", binary), args...)
	cmd.Stdin = input
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("%s %v: %v", binary, args, err)
	}
	var rows []string
	if len(output) > 0 {
		rows = strings.Split(strings.TrimSuffix(string(output), "\n"), "\n")
	}
	if len(rows) != count {
		log.Fatalf("%s %v produced %d rows, expected %d", binary, args, len(rows), count)
	}
	return rows
}

// Exhaustive checks one row per unsigned input, in ascending order.
func Exhaustive(binary string, count int, outcomes func(int) (actual, expected int)) {
	for x, model := range LeanRows(binary, count, nil) {
		actual, expected := outcomes(x)
		if actual != expected || model != strconv.Itoa(expected) {
			log.Fatalf("%s input %#x: Go=%d, Lean=%q, expected=%d", binary, x, actual, model, expected)
		}
	}
	fmt.Printf("%s: %s and Lean agree with the specification on all %d inputs.\n", binary, runtime.Version(), count)
}
