package flate

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func geometry(size int) (description string) {
	defer func() {
		if failure := recover(); failure != nil {
			if _, ok := failure.(runtime.Error); !ok {
				panic(failure)
			}
			description = "bounds"
		}
	}()
	encoder := newHuffmanEncoder(size)
	clear := encoder.bitCount == [17]int32{} && encoder.freqcache == [literalCount + 1]literalNode{}
	for _, code := range encoder.codes[:cap(encoder.codes)] {
		clear = clear && code == 0
	}
	return fmt.Sprintf("%d %d %t", len(encoder.codes), cap(encoder.codes), clear)
}

func TestGoLeanEncoderInit(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_FLATE_INIT_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	s := bufio.NewScanner(file)
	node := maxNode()
	count := 0
	for s.Scan() {
		sizeField, want, ok := strings.Cut(s.Text(), " ")
		if !ok {
			t.Fatalf("bad fixture %q", s.Text())
		}
		size, err := strconv.Atoi(sizeField)
		if err != nil {
			t.Fatal(err)
		}
		got := fmt.Sprintf("%s|%d %d", geometry(size), node.literal, node.freq)
		if got != want {
			t.Fatalf("size %d: got %q want %q", size, got, want)
		}
		count++
	}
	if err := s.Err(); err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("empty fixture file")
	}
}
