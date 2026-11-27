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

type rgbaModelProbe struct {
	words  [4]uint32
	calls  int
	panics bool
}

func (p *rgbaModelProbe) RGBA() (uint32, uint32, uint32, uint32) {
	p.calls++
	if p.panics {
		panic(p)
	}
	return p.words[0], p.words[1], p.words[2], p.words[3]
}

func rgbaModelOutcome(width, tag int, words [4]uint32) (out string) {
	probe := &rgbaModelProbe{words: words, panics: tag == 2}
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
			c = RGBA{uint8(words[0]), uint8(words[1]), uint8(words[2]), uint8(words[3])}
		} else {
			c = RGBA64{uint16(words[0]), uint16(words[1]), uint16(words[2]), uint16(words[3])}
		}
	case 1, 2:
		c = probe
	case 3:
	default:
		panic("unknown rgba input")
	}
	if width == 8 {
		result := rgbaModel(c).(RGBA)
		return fmt.Sprintf("%d %d %d %d %d", result.R, result.G, result.B, result.A, probe.calls)
	}
	result := rgba64Model(c).(RGBA64)
	return fmt.Sprintf("%d %d %d %d %d", result.R, result.G, result.B, result.A, probe.calls)
}

func TestGoLeanRGBAModel(t *testing.T) {
	dir := os.Getenv("GOLEAN_RGBAMODEL_DIR")
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
		fmt.Fprintln(w, rgbaModelOutcome(width, tag, words))
	}
	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}
}
