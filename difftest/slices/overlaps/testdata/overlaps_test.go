package slices_test

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

const slots = 64

func goLeanOverlaps[E any](aStart, aLength, bStart, bLength int, separate bool) bool {
	var array [slots]E
	var second [slots]E
	b := &array
	if separate {
		b = &second
	}
	return slices.GoLeanOverlaps(array[aStart:aStart+aLength], b[bStart:bStart+bLength])
}

func TestGoLeanOverlaps(t *testing.T) {
	dir := os.Getenv("GOLEAN_OVERLAPS_DIR")
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
		var a, al, b, bl, size int
		var separate bool
		n, err := fmt.Fscan(reader, &a, &al, &b, &bl, &size, &separate)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		var got bool
		switch size {
		case 0:
			got = goLeanOverlaps[struct{}](a, al, b, bl, separate)
		case 1:
			got = goLeanOverlaps[byte](a, al, b, bl, separate)
		case 3:
			got = goLeanOverlaps[[3]byte](a, al, b, bl, separate)
		case 8:
			got = goLeanOverlaps[uint64](a, al, b, bl, separate)
		default:
			t.Fatalf("unexpected element size %d", size)
		}
		if _, err := fmt.Fprintln(writer, got); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Flush(); err != nil {
		t.Fatal(err)
	}
}
