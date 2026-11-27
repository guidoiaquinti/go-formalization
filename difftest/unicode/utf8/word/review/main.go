package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	width := strconv.IntSize / 8
	fixtures := []struct {
		in     string
		want32 uint64
		want64 uint64
	}{
		{"\x01\x23\x45\x67\x89\xab\xcd\xef", 0x67452301, 0xefcdab8967452301},
		{"\x80\x00\x00\x00\x00\x00\x00\x00", 0x80, 0x80},
		{"\x00\x00\x00\xff\x00\x00\x00\x80\xff", 0xff000000, 0x80000000ff000000},
		{"\x01\x00\x00\x00\x00\x00\x00\x00\xff", 1, 1},
		{"\xff\xff\xff", 0, 0},
		{"", 0, 0},
	}
	var rows bytes.Buffer
	var endianFailures, paddingFailures int
	for _, f := range fixtures {
		if len(f.in) < width {
			fmt.Fprintf(&rows, "x%x bounds\n", f.in)
			paddingFailures++
			continue
		}
		want := f.want64
		reversed := uint64(0)
		if width == 4 {
			want = f.want32
			reversed = uint64(binary.BigEndian.Uint32([]byte(f.in)))
		} else {
			reversed = binary.BigEndian.Uint64([]byte(f.in))
		}
		fmt.Fprintf(&rows, "x%x %d\n", f.in, want)
		if reversed != want {
			endianFailures++
		}
	}
	if endianFailures != 4 || paddingFailures != 2 {
		log.Fatalf("big-endian mutant failed %d, zero-pad-short mutant failed %d", endianFailures, paddingFailures)
	}
	dir, err := os.MkdirTemp("", "golean-word-review-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, rows.Bytes(), 0600); err != nil {
		log.Fatal(err)
	}
	if err := os.Setenv("GOLEAN_UTF8_WORD_CASES", path); err != nil {
		log.Fatal(err)
	}
	const overlay = "difftest/unicode/utf8/word/testdata/"
	if err := check.TestInPackage("unicode/utf8", "TestGoLeanWord", overlay+"bridge_test.go", overlay+"word_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("both native %d-byte word instantiations passed %d fixtures; mutants failed %d and %d\n", width, len(fixtures), endianFailures, paddingFailures)
}
