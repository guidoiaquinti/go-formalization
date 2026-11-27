package flate

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func dictBytes(b []byte) string {
	cells := make([]string, len(b))
	for i, v := range b {
		cells[i] = strconv.Itoa(int(v))
	}
	return strings.Join(cells, ",")
}

func dictStore(d *dictDecoder, value byte) (ok bool) {
	defer func() {
		if e := recover(); e != nil {
			failure, ok := e.(runtime.Error)
			if !ok || !strings.Contains(failure.Error(), "index out of range") {
				panic(e)
			}
		}
	}()
	d.writeByte(value)
	return true
}

func TestGoLeanDictState(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_FLATE_DICT_STATE"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	input := json.NewDecoder(f)
	for i := 0; ; i++ {
		var c struct {
			Backing                      []byte
			Length, Read, Written, Count int
			Full                         bool
			Value                        byte
			Want                         string
		}
		err := input.Decode(&c)
		if err == io.EOF && i > 0 {
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		backing := slices.Clone(c.Backing)
		d := dictDecoder{hist: backing[:c.Length:len(backing)], rdPos: c.Read, wrPos: c.Written, full: c.Full}
		state := func(state dictDecoder, b []byte) string {
			return fmt.Sprintf("%d:%d:%t:%s", state.rdPos, state.wrPos, state.full, dictBytes(b))
		}
		view := func(v []byte, origin int) string {
			if cap(v) > 0 && &v[:cap(v)][0] != &backing[origin] {
				t.Fatal("view lost backing identity")
			}
			return fmt.Sprintf("%t:%d:%d:%s", v != nil, cap(v), origin, dictBytes(v))
		}
		first := []string{fmt.Sprint(d.histSize()), fmt.Sprint(d.availRead()), fmt.Sprint(d.availWrite()), view(d.writeSlice(), c.Written)}
		marked := d
		marked.writeMark(c.Count)
		first = append(first, state(marked, backing))
		updated := slices.Clone(backing)
		stored := d
		stored.hist = updated[:c.Length]
		if dictStore(&stored, c.Value) {
			first = append(first, state(stored, updated))
		} else {
			if !slices.Equal(backing, updated) || stored.wrPos != d.wrPos {
				t.Fatal("failed store changed state")
			}
			first = append(first, "bounds")
		}
		first = append(first, view(d.readFlush(), c.Read), state(d, backing))
		if got := strings.Join(first, "|"); got != c.Want {
			t.Fatalf("row %d: Go %q want %q", i, got, c.Want)
		}
	}
}
