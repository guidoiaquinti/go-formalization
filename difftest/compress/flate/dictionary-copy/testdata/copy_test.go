package flate

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanDictionaryCopy(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_FLATE_DICT_COPY"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	in := json.NewDecoder(f)
	for row := 0; ; row++ {
		var c struct {
			Backing                                []byte
			Length, Read, Written, Distance, Count int
			Full                                   bool
			Want                                   string
		}
		err := in.Decode(&c)
		if err == io.EOF && row > 0 {
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		var outcomes []string
		for _, method := range []func(*dictDecoder, int, int) int{(*dictDecoder).writeCopy, (*dictDecoder).tryWriteCopy} {
			b := slices.Clone(c.Backing)
			state := dictDecoder{hist: b[:c.Length:len(b)], rdPos: c.Read, wrPos: c.Written, full: c.Full}
			n := method(&state, c.Distance, c.Count)
			words := make([]string, len(b))
			for i, cell := range b {
				words[i] = strconv.Itoa(int(cell))
			}
			outcomes = append(outcomes, fmt.Sprintf("%d:%d:%d:%t:%s", n, state.rdPos, state.wrPos, state.full, strings.Join(words, ",")))
		}
		if got := strings.Join(outcomes, "|"); got != c.Want {
			t.Fatalf("row %d: Go=%q want=%q", row, got, c.Want)
		}
	}
}
