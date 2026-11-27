package flate

import (
	"errors"
	"io"
	"testing"
)

type reviewNoEOFHooks struct {
	calls int
}

func (e *reviewNoEOFHooks) Error() string {
	e.calls++
	return "EOF"
}

func (e *reviewNoEOFHooks) Is(target error) bool {
	e.calls++
	return target == io.EOF
}

func (e *reviewNoEOFHooks) Unwrap() error {
	e.calls++
	return io.EOF
}

type reviewNoEOFMap map[string]int

func (reviewNoEOFMap) Error() string { return "EOF" }

func TestReviewGoLeanNoEOF(t *testing.T) {
	if got := noEOF(io.EOF); got != io.ErrUnexpectedEOF {
		t.Fatal("exact EOF was not replaced by the unexpected EOF sentinel")
	}
	joined := errors.Join(io.EOF, io.ErrUnexpectedEOF)
	if got := noEOF(joined); got != joined {
		t.Fatal("wrapped sentinel identity changed")
	}
	hooked := new(reviewNoEOFHooks)
	if got := noEOF(hooked); got != hooked || hooked.calls != 0 {
		t.Fatal("normalization changed the error or called a user hook")
	}
	mapping := reviewNoEOFMap{"before": 1}
	got, ok := noEOF(mapping).(reviewNoEOFMap)
	if !ok || got["before"] != 1 {
		t.Fatal("incomparable map-backed error changed")
	}
	got["after"] = 2
	if mapping["after"] != 2 {
		t.Fatal("incomparable map was copied")
	}
	if got := noEOF((*reviewNoEOFHooks)(nil)); got != (*reviewNoEOFHooks)(nil) {
		t.Fatal("typed nil changed")
	}
}
