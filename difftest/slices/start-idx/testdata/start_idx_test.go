package slices_test

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"testing"
)

const slots = 64

func goLeanStartIdx(a, al, b, bl int, separate bool) (row string) {
	// All element values are equal. Pointer identity must select the result.
	var first, second [slots]uint64
	source := &first
	if separate {
		source = &second
	}
	defer func() {
		switch p := recover().(type) {
		case nil:
		case runtime.Error:
			row = "bounds panic"
		case string:
			row = "panic " + p
		default:
			panic(p)
		}
	}()
	return fmt.Sprintf("index %d", slices.GoLeanStartIdx(first[a:a+al], source[b:b+bl]))
}

func TestGoLeanStartIdx(t *testing.T) {
	dir := os.Getenv("GOLEAN_START_IDX_DIR")
	in, err := os.Open(filepath.Join(dir, "inputs"))
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()
	out, err := os.Create(filepath.Join(dir, "outputs"))
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()
	reader, writer := bufio.NewReader(in), bufio.NewWriter(out)
	for {
		var a, al, b, bl int
		var separate bool
		n, err := fmt.Fscan(reader, &a, &al, &b, &bl, &separate)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if _, err := fmt.Fprintln(writer, goLeanStartIdx(a, al, b, bl, separate)); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Flush(); err != nil {
		t.Fatal(err)
	}
}
