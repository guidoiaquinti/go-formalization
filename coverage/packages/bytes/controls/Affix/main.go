package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

type sliceView struct {
	text string
	cap  int
	nil  bool
}

func view(s []byte) sliceView {
	return sliceView{string(s), cap(s), s == nil}
}

type testCase struct {
	src, affix        []byte
	pre, post         bool
	leading, trailing sliceView
}

func main() {
	pin, err := os.ReadFile(".go-version")
	if err != nil {
		log.Fatal(err)
	}
	if runtime.Version() != "go"+strings.TrimSpace(string(pin)) {
		log.Fatalf("unexpected Go version %s", runtime.Version())
	}
	storage := [8]byte{0xff, 0, 'x'}
	cases := []testCase{
		{nil, []byte{}, true, true, sliceView{"", 0, true}, sliceView{"", 0, true}},
		{[]byte{}, nil, true, true, sliceView{"", 0, false}, sliceView{"", 0, false}},
		{storage[:3], []byte{0xff}, true, false, sliceView{"\x00x", 7, false}, sliceView{"\xff\x00x", 8, false}},
		{storage[:3], []byte{'x'}, false, true, sliceView{"\xff\x00x", 8, false}, sliceView{"\xff\x00", 8, false}},
		{storage[:3], []byte{0xff, 0, 'x'}, true, true, sliceView{"", 5, false}, sliceView{"", 8, false}},
		{storage[:3], []byte{0xff, 0, 'x', 0}, false, false, sliceView{"\xff\x00x", 8, false}, sliceView{"\xff\x00x", 8, false}},
	}
	var real, mutants [6]int
	for _, c := range cases {
		pre, post := bytes.HasPrefix(c.src, c.affix), bytes.HasSuffix(c.src, c.affix)
		leading, preFound := bytes.CutPrefix(c.src, c.affix)
		trailing, postFound := bytes.CutSuffix(c.src, c.affix)
		trimPre, trimPost := bytes.TrimPrefix(c.src, c.affix), bytes.TrimSuffix(c.src, c.affix)
		passed := [6]bool{
			pre == c.pre, post == c.post,
			view(leading) == c.leading && preFound == c.pre,
			view(trailing) == c.trailing && postFound == c.post,
			view(trimPre) == c.leading, view(trimPost) == c.trailing,
		}
		// Substitute an empty-affix rejection or clip each matching result to length.
		mutated := [6]bool{
			(pre && len(c.affix) != 0) == c.pre,
			(post && len(c.affix) != 0) == c.post,
			view(clipMatch(leading, preFound)) == c.leading && preFound == c.pre,
			view(clipMatch(trailing, postFound)) == c.trailing && postFound == c.post,
			view(clipMatch(trimPre, pre)) == c.leading,
			view(clipMatch(trimPost, post)) == c.trailing,
		}
		for i := range passed {
			if !passed[i] {
				real[i]++
			}
			if !mutated[i] {
				mutants[i]++
			}
		}
	}
	if real != [6]int{} || mutants != [6]int{2, 2, 2, 2, 2, 2} {
		log.Fatalf("affix failures: Go=%v mutants=%v", real, mutants)
	}
	fmt.Printf("%s %s/%s: bytes affix Go=%v mutants=%v across %d cases\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, real, mutants, len(cases))
}

func clipMatch(s []byte, found bool) []byte {
	if found {
		return s[:len(s):len(s)]
	}
	return s
}
