package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

type write func(*strings.Builder, []byte) (int, error)

func failure(fn write, b *strings.Builder, p []byte) (out string) {
	defer func() {
		switch value := recover().(type) {
		case nil:
		case runtime.Error:
			out = "nil"
		case string:
			if value != "strings: illegal use of non-zero Builder copied by value" {
				panic(value)
			}
			out = "copy"
		default:
			panic(value)
		}
	}()
	_, err := fn(b, p)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func check(fn write) [5]bool {
	var b strings.Builder
	b.WriteString("a")
	original := b.String()
	input := []byte{0xff, 0}
	n, err := fn(&b, input)
	input[0] = 'z'
	good := n == 2 && err == nil && original == "a" && b.String() == "a\xff\x00"
	copyOfB := b
	rejected := failure(fn, &copyOfB, []byte("x")) == "copy" && copyOfB.String() == b.String()
	var zero strings.Builder
	n, err = fn(&zero, nil)
	copyOfZero := zero
	nilSetsIdentity := n == 0 && err == nil && failure(fn, &copyOfZero, []byte("x")) == "copy"
	var empty strings.Builder
	n, err = fn(&empty, []byte{})
	copyOfEmpty := empty
	emptySetsIdentity := n == 0 && err == nil && failure(fn, &copyOfEmpty, []byte("x")) == "copy"
	return [5]bool{good, rejected, nilSetsIdentity, emptySetsIdentity, failure(fn, nil, []byte("x")) == "nil"}
}

func main() {
	for _, c := range []struct {
		name string
		fn   write
		want int
	}{
		{"Go", (*strings.Builder).Write, 0},
		{"Skip empty writes", func(b *strings.Builder, p []byte) (int, error) {
			if len(p) == 0 {
				return 0, nil
			}
			return b.Write(p)
		}, 2},
		{"Report combined length", func(b *strings.Builder, p []byte) (int, error) { _, err := b.Write(p); return b.Len(), err }, 1},
	} {
		bad := 0
		for _, pass := range check(c.fn) {
			if !pass {
				bad++
			}
		}
		fmt.Printf("%s disagreed on %d of 5 checks\n", c.name, bad)
		if bad != c.want {
			log.Fatalf("%s: wanted %d failures", c.name, c.want)
		}
	}
}
