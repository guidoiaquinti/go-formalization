package main

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

func main() {
	cases := []string{"", "longer", "é🙂", "\x00\xff"}
	var real, skippedByte int
	for _, input := range cases {
		reader := strings.NewReader(input)
		if previousRune(reader) != -1 {
			panic("new Reader did not initialize its previous-rune marker")
		}
		out, err := io.ReadAll(reader)
		if err != nil {
			panic(err)
		}
		if string(out) != input {
			real++
		}
		if input[min(1, len(input)):] != input {
			skippedByte++
		}
	}
	var zeroMarker strings.Reader
	markerFailed := previousRune(&zeroMarker) != -1
	fmt.Printf("NewReader: Go=%d/4 skip-first-byte=%d/4 zero-marker-rejected=%t\n", real, skippedByte, markerFailed)
	if real != 0 || skippedByte != 3 || !markerFailed {
		panic("unexpected NewReader control outcomes")
	}
}

func previousRune(reader *strings.Reader) int64 {
	return reflect.ValueOf(reader).Elem().FieldByName("prevRune").Int()
}
