package flate

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"testing"
)

func tryLE(action func() string) (answer string) {
	defer func() {
		if e := recover(); e != nil {
			fault, ok := e.(runtime.Error)
			if !ok || !strings.Contains(fault.Error(), "out of range") {
				panic(e)
			}
			answer = "bounds"
		}
	}()
	return action()
}

func loadObservations[I indexer](view []byte, start I) []string {
	return []string{
		tryLE(func() string { return fmt.Sprint(loadLE8(view, start)) }),
		tryLE(func() string { return fmt.Sprint(loadLE32(view, start)) }),
		tryLE(func() string { return fmt.Sprint(loadLE64(view, start)) }),
	}
}

func TestGoLeanFlateEndian(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_FLATE_ENDIAN_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	for row := 0; ; row++ {
		var c struct {
			Tag, Length  int
			Index, Value uint64
			Cells        []byte
			Want         string
		}
		err := dec.Decode(&c)
		if err == io.EOF && row > 0 {
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		view := c.Cells[:c.Length:len(c.Cells)]
		var loads []string
		switch c.Tag {
		case 0:
			loads = loadObservations(view, int8(c.Index))
		case 1:
			loads = loadObservations(view, uint8(c.Index))
		case 2:
			loads = loadObservations(view, int16(c.Index))
		case 3:
			loads = loadObservations(view, uint16(c.Index))
		case 4:
			loads = loadObservations(view, int32(c.Index))
		case 5:
			loads = loadObservations(view, uint32(c.Index))
		case 6:
			loads = loadObservations(view, int64(c.Index))
		case 7:
			loads = loadObservations(view, c.Index)
		case 8:
			loads = loadObservations(view, int(c.Index))
		case 9:
			loads = loadObservations(view, uint(c.Index))
		default:
			t.Fatalf("unknown numeric type %d", c.Tag)
		}
		stored := tryLE(func() string { storeLE64(view, c.Value); return "true" })
		if stored == "bounds" {
			stored = "false"
		}
		words := make([]string, len(c.Cells))
		for i, b := range c.Cells {
			words[i] = fmt.Sprint(b)
		}
		got := strings.Join(append(loads, stored, strings.Join(words, ",")), "|")
		if got != c.Want {
			t.Fatalf("row %d tag %d index %#x: Go %q want %q", row, c.Tag, c.Index, got, c.Want)
		}
	}
}
