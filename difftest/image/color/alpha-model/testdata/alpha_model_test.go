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

type alphaModelProbe struct {
	value  uint32
	calls  int
	panics bool
}

func (p *alphaModelProbe) RGBA() (uint32, uint32, uint32, uint32) {
	p.calls++
	if p.panics {
		panic(p)
	}
	return 0x1234, 0x5678, 0x9abc, p.value
}

func alphaModelOutcome(tag int, value uint32) (out string) {
	probe := &alphaModelProbe{value: value, panics: tag == 2}
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
		c = Alpha{uint8(value)}
	case 1, 2:
		c = probe
	case 3:
	default:
		panic("unknown alpha input")
	}
	result := alphaModel(c).(Alpha)
	return fmt.Sprintf("%d %d", result.A, probe.calls)
}

func TestGoLeanAlphaModel(t *testing.T) {
	dir := os.Getenv("GOLEAN_ALPHAMODEL_DIR")
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
		var tag int
		var value uint32
		n, err := fmt.Fscan(r, &tag, &value)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprintln(w, alphaModelOutcome(tag, value))
	}
	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}
}
