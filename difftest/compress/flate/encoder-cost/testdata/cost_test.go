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

func observeCost(run func() int) (result string) {
	defer func() {
		if p := recover(); p != nil {
			message, ok := p.(runtime.Error)
			if !ok {
				panic(p)
			}
			switch {
			case strings.Contains(message.Error(), "nil pointer"):
				result = "nil"
			case strings.Contains(message.Error(), "index out of range"):
				result = "bounds"
			default:
				panic(p)
			}
		}
	}()
	return fmt.Sprint(uint(run()))
}

func TestGoLeanEncoderCost(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_FLATE_COST_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	for count := 0; ; count++ {
		var c struct {
			Present bool
			Words   []uint32
			Freq    []uint16
			Data    []byte
			Want    string
		}
		err := dec.Decode(&c)
		if err == io.EOF && count > 0 {
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		var encoder *huffmanEncoder
		if c.Present {
			encoder = &huffmanEncoder{codes: make([]hcode, len(c.Words))}
			for i, word := range c.Words {
				encoder.codes[i] = hcode(word)
			}
		}
		got := strings.Join([]string{
			observeCost(func() int { return encoder.bitLength(c.Freq) }),
			observeCost(func() int { return encoder.canEncodeLen(c.Freq) }),
			observeCost(func() int { return encoder.bitLengthRaw(c.Data) }),
		}, " ")
		if got != c.Want {
			t.Fatalf("row %d: Go=%q want=%q", count, got, c.Want)
		}
	}
}
