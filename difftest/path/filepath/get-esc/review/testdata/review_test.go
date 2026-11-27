package filepath_test

import (
	"errors"
	"path/filepath"
	"runtime"
	"testing"
	"unicode/utf8"
)

type answer struct {
	r    rune
	rest string
	err  error
}

func native(text string) answer {
	r, rest, err := filepath.GoLeanGetEsc(text)
	return answer{r, rest, err}
}

func TestGoLeanGetEscReview(t *testing.T) {
	badPattern := filepath.ErrBadPattern
	windows := runtime.GOOS == "windows"
	slashRune := rune(0)
	escaped := answer{'-', "]", nil}
	zeroOnError := 3
	if windows {
		slashRune = '\\'
		escaped = answer{'\\', "-]", nil}
		zeroOnError++
	}
	fixtures := []struct {
		text string
		want answer
	}{
		{"", answer{0, "", badPattern}},
		{"-x", answer{0, "", badPattern}},
		{"]x", answer{0, "", badPattern}},
		{"\\", answer{slashRune, "", badPattern}},
		{"\\-]", escaped},
		{"é", answer{'é', "", badPattern}},
		{"é\xff", answer{'é', "\xff", nil}},
		{"\xffz", answer{utf8.RuneError, "z", badPattern}},
		{"\xe2\x82", answer{utf8.RuneError, "\x82", badPattern}},
		{"�]", answer{utf8.RuneError, "]", nil}},
		{"🚀x", answer{'🚀', "x", nil}},
	}
	for _, c := range []struct {
		name  string
		parse func(string) answer
		want  int
	}{
		{"Go", native, 0},
		{"Discard tuple fields on error", func(s string) answer {
			got := native(s)
			if got.err != nil {
				return answer{err: got.err}
			}
			return got
		}, zeroOnError},
		{"Reject literal RuneError", func(s string) answer {
			got := native(s)
			if got.r == utf8.RuneError {
				got.err = badPattern
			}
			return got
		}, 1},
		{"Return equivalent error text with different identity", func(s string) answer {
			got := native(s)
			if got.err != nil {
				got.err = errors.New(got.err.Error())
			}
			return got
		}, 7},
	} {
		bad := 0
		for _, f := range fixtures {
			if got := c.parse(f.text); got != f.want {
				bad++
				if c.want == 0 {
					t.Errorf("getEsc(%x) gave (%U,%x,%v), want (%U,%x,%v)", f.text, got.r, got.rest, got.err, f.want.r, f.want.rest, f.want.err)
				}
			}
		}
		t.Logf("%s disagreed on %d of %d", c.name, bad, len(fixtures))
		if bad != c.want {
			t.Errorf("%s wanted %d differences", c.name, c.want)
		}
	}
}
