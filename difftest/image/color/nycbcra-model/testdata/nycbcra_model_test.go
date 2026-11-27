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

type goLeanNYCbCrAProbe struct {
	r, g, b, a uint32
	calls      int
	panics     bool
}

func (p *goLeanNYCbCrAProbe) RGBA() (uint32, uint32, uint32, uint32) {
	p.calls++
	if p.panics {
		panic(p)
	}
	return p.r, p.g, p.b, p.a
}

func goLeanNYCbCrAOutcome(tag int, r, g, b, a uint32) (row string) {
	probe := &goLeanNYCbCrAProbe{r: r, g: g, b: b, a: a, panics: tag == 3}
	defer func() {
		if p := recover(); p != nil {
			_, runtimePanic := p.(runtime.Error)
			if p != probe && (tag != 4 || !runtimePanic) {
				panic(p)
			}
			row = fmt.Sprintf("panic %d", probe.calls)
		}
	}()
	var c Color
	switch tag {
	case 0:
		c = NYCbCrA{YCbCr{uint8(r), uint8(g), uint8(b)}, uint8(a)}
	case 1:
		c = YCbCr{uint8(r), uint8(g), uint8(b)}
	case 2, 3:
		c = probe
	case 4:
	default:
		panic("unknown NYCbCrA input")
	}
	out := nYCbCrAModel(c).(NYCbCrA)
	return fmt.Sprintf("%d %d %d %d %d", out.Y, out.Cb, out.Cr, out.A, probe.calls)
}

func TestGoLeanNYCbCrAModel(t *testing.T) {
	dir := os.Getenv("GOLEAN_NYCBCRA_MODEL_DIR")
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
		if _, err := fmt.Fprintln(writer, goLeanNYCbCrAOutcome(tag, r, g, b, a)); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Flush(); err != nil {
		t.Fatal(err)
	}
}
