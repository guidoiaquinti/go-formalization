package flate

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
)

type noEOFPointer struct{}

func (*noEOFPointer) Error() string { return "EOF" }

type noEOFList []int

func (noEOFList) Error() string { return "EOF" }

func TestGoLeanNoEOF(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_FLATE_NO_EOF"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	input := bufio.NewReader(f)
	cases := []error{nil, io.EOF, io.ErrUnexpectedEOF, errors.New("EOF"),
		fmt.Errorf("wrapped: %w", io.EOF), (*noEOFPointer)(nil), noEOFList{42}, new(noEOFPointer)}
	for i, err := range cases {
		var want int
		if _, readErr := fmt.Fscan(input, &want); readErr != nil {
			t.Fatal(readErr)
		}
		got := noEOF(err)
		if i == 6 {
			before := err.(noEOFList)
			after, ok := got.(noEOFList)
			if want != i || !ok || !reflect.DeepEqual(before, after) || &before[0] != &after[0] {
				t.Fatalf("incomparable error changed: %T %v", got, got)
			}
			continue
		}
		if want < 0 || want >= len(cases) || got != cases[want] {
			t.Fatalf("case %d: got %T %v, expected identity %d", i, got, got, want)
		}
	}
	var extra int
	if n, err := fmt.Fscan(input, &extra); n != 0 || err != io.EOF {
		t.Fatal("unexpected trailing result")
	}
}
