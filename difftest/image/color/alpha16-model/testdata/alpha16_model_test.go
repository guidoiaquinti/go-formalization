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

type alpha16ModelProbe struct {
	value  uint32
	calls  int
	panics bool
}

func (p *alpha16ModelProbe) RGBA() (uint32, uint32, uint32, uint32) {
	p.calls++
	if p.panics {
		panic(p)
	}
	return 0x1234, 0x5678, 0x9abc, p.value
}

func alpha16ModelOutcome(tag int, value uint32) (out string) {
	probe := &alpha16ModelProbe{value: value, panics: tag == 2}
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
		c = Alpha16{uint16(value)}
	case 1, 2:
		c = probe
	case 3:
	default:
		panic("unknown alpha input")
	}
	result := alpha16Model(c).(Alpha16)
	return fmt.Sprintf("%d %d", result.A, probe.calls)
}

func TestGoLeanAlpha16Model(t *testing.T) {
	dir := os.Getenv("GOLEAN_ALPHA16MODEL_DIR")
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
		fmt.Fprintln(w, alpha16ModelOutcome(tag, value))
	}
	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}
}
