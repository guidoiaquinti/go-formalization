package suffixarray

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"testing"
)

type signedWord interface{ ~int32 | ~int64 }

type expandFixture struct {
	Mode, Bins int
	Text       []byte
	Prefix     []int
	Want       string
}

func TestGoLeanExpandLms(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_SUFFIX_EXPAND_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	dec := gob.NewDecoder(f)
	for row := 0; ; row++ {
		var c expandFixture
		if err := dec.Decode(&c); err != nil {
			if err == io.EOF && row > 0 {
				t.Logf("%d pinned native cases per declaration", row)
				return
			}
			t.Fatal(err)
		}
		n32 := nativeExpand(c, true, func(freq, bucket, sa []int32) {
			expand_8_32(c.Text, freq, bucket, sa, len(c.Prefix))
		})
		n64 := nativeExpand(c, true, func(freq, bucket, sa []int64) {
			expand_8_64(c.Text, freq, bucket, sa, len(c.Prefix))
		})
		w32 := nativeExpand(c, false, func(freq, bucket, sa []int32) {
			input := make([]int32, len(c.Text))
			for i, b := range c.Text {
				input[i] = 2 * int32(b)
			}
			expand_32(input, freq, bucket, sa, len(c.Prefix))
		})
		w64 := nativeExpand(c, false, func(freq, bucket, sa []int64) {
			input := make([]int64, len(c.Text))
			for i, b := range c.Text {
				input[i] = 2 * int64(b)
			}
			expand_64(input, freq, bucket, sa, len(c.Prefix))
		})
		got := strings.Join([]string{n32, n64, w32, w64}, "|")
		byteRef, wideRef := referenceExpand(c, true), referenceExpand(c, false)
		want := strings.Join([]string{byteRef, byteRef, wideRef, wideRef}, "|")
		if got != c.Want || got != want {
			t.Fatalf("row %d %+v: Go %q Lean %q independent %q", row, c, got, c.Want, want)
		}
	}
}

func geometry(c expandFixture, byteMode bool) (bins, scale int) {
	if byteMode {
		return 256, 1
	}
	return c.Bins, 2
}

func nativeExpand[T signedWord](c expandFixture, byteMode bool, run func([]T, []T, []T)) string {
	bins, scale := geometry(c, byteMode)
	var freq []T
	if c.Mode != 0 {
		freq = make([]T, bins+2)
		freq[bins], freq[bins+1] = 93, 94
		if c.Mode == 1 {
			for i := range bins {
				freq[i] = 17
			}
			freq[0] = -1
		} else {
			for _, b := range c.Text {
				freq[int(b)*scale]++
			}
		}
	}
	bucket := make([]T, bins+2)
	for i := range bins {
		bucket[i] = 19
	}
	bucket[bins], bucket[bins+1] = 91, 92
	sa := make([]T, len(c.Text))
	for i := range sa {
		sa[i] = -7
	}
	for i, pos := range c.Prefix {
		sa[i] = T(pos)
	}
	run(freq[:min(len(freq), bins)], bucket[:bins], sa)
	return encodeExpansion(sa, freq, bucket)
}

func referenceExpand(c expandFixture, byteMode bool) string {
	bins, scale := geometry(c, byteMode)
	counts := make([]int64, bins)
	group := make([][]int, bins)
	for _, b := range c.Text {
		counts[int(b)*scale]++
	}
	for _, pos := range c.Prefix {
		sym := int(c.Text[pos]) * scale
		group[sym] = append(group[sym], pos)
	}
	var freq []int64
	if c.Mode != 0 {
		freq = append(slices.Clone(counts), 93, 94)
	}
	sa := make([]int64, len(c.Text))
	bucket := make([]int64, bins)
	offset := 0
	for sym, count := range counts {
		offset += int(count)
		start := offset - len(group[sym])
		bucket[sym] = int64(start)
		for i, pos := range group[sym] {
			sa[start+i] = int64(pos)
		}
	}
	return encodeExpansion(sa, freq, append(bucket, 91, 92))
}

func encodeExpansion[T signedWord](sa, freq, bucket []T) string {
	return encodeRuns(sa) + "/" + encodeRuns(freq) + "/" + encodeRuns(bucket)
}

func encodeRuns[T signedWord](words []T) string {
	if words == nil {
		return "nil"
	}
	var parts []string
	for i := 0; i < len(words); {
		j := i + 1
		for j < len(words) && words[j] == words[i] {
			j++
		}
		parts = append(parts, fmt.Sprintf("%d.%d", words[i], j-i))
		i = j
	}
	return strings.Join(parts, ",")
}
