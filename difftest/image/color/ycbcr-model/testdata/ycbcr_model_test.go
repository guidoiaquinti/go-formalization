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

type goLeanYCbCrProbe struct {
	r, g, b, a uint32
	calls      int
	panics     bool
}

func (p *goLeanYCbCrProbe) RGBA() (uint32, uint32, uint32, uint32) {
	p.calls++
	if p.panics {
		panic(p)
	}
	return p.r, p.g, p.b, p.a
}

func goLeanYCbCrOutcome(tag int, r, g, b, a uint32) (row string) {
	probe := &goLeanYCbCrProbe{r: r, g: g, b: b, a: a, panics: tag == 2}
	defer func() {
		if p := recover(); p != nil {
			_, runtimePanic := p.(runtime.Error)
			if p != probe && (tag != 3 || !runtimePanic) {
				panic(p)
			}
			row = fmt.Sprintf("panic %d", probe.calls)
		}
	}()
	var c Color
	switch tag {
	case 0:
		c = YCbCr{uint8(r), uint8(g), uint8(b)}
	case 1, 2:
		c = probe
	case 3:
	default:
		panic("unknown YCbCr input")
	}
	out := yCbCrModel(c).(YCbCr)
	return fmt.Sprintf("%d %d %d %d", out.Y, out.Cb, out.Cr, probe.calls)
}

func TestGoLeanYCbCrModel(t *testing.T) {
	dir := os.Getenv("GOLEAN_YCBCR_MODEL_DIR")
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
	reader, writer := bufio.NewReader(in), bufio.NewWriter(out)
	for {
		var tag int
		var r, g, b, a uint32
		n, err := fmt.Fscan(reader, &tag, &r, &g, &b, &a)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if _, err := fmt.Fprintln(writer, goLeanYCbCrOutcome(tag, r, g, b, a)); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Flush(); err != nil {
		t.Fatal(err)
	}
}
