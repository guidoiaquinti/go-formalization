package color

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

type nrgbaModelProbe struct {
	words  [4]uint32
	calls  int
	panics bool
}

func (p *nrgbaModelProbe) RGBA() (uint32, uint32, uint32, uint32) {
	p.calls++
	if p.panics {
		panic(p)
	}
	return p.words[0], p.words[1], p.words[2], p.words[3]
}

func nrgbaModelOutcome(width, tag int, words [4]uint32) (out string) {
	probe := &nrgbaModelProbe{words: words, panics: tag == 2}
	defer func() {
		if p := recover(); p != nil {
			_, runtimePanic := p.(runtime.Error)
			if p == probe || tag == 3 && runtimePanic {
				out = fmt.Sprintf("panic %d", probe.calls)
			} else {
				panic(p)
			}
		}
	}()
	var c Color
	switch tag {
	case 0:
		if width == 8 {
			c = NRGBA{uint8(words[0]), uint8(words[1]), uint8(words[2]), uint8(words[3])}
		} else {
			c = NRGBA64{uint16(words[0]), uint16(words[1]), uint16(words[2]), uint16(words[3])}
		}
	case 1, 2:
		c = probe
	case 3:
	default:
		panic("unknown rgba input")
	}
	if width == 8 {
		result := nrgbaModel(c).(NRGBA)
		return fmt.Sprintf("%d %d %d %d %d", result.R, result.G, result.B, result.A, probe.calls)
	}
	result := nrgba64Model(c).(NRGBA64)
	return fmt.Sprintf("%d %d %d %d %d", result.R, result.G, result.B, result.A, probe.calls)
}

func TestGoLeanNRGBAModel(t *testing.T) {
	dir := os.Getenv("GOLEAN_NRGBAMODEL_DIR")
	in, err := os.Open(filepath.Join(dir, "inputs"))
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()
	out, err := os.Create(filepath.Join(dir, "outputs"))
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()
	r, w := bufio.NewReader(in), bufio.NewWriter(out)
	for {
		var width, tag int
		var words [4]uint32
		n, err := fmt.Fscan(r, &width, &tag, &words[0], &words[1], &words[2], &words[3])
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprintln(w, nrgbaModelOutcome(width, tag, words))
	}
	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}
}
