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

type grayModelProbe struct {
	r, g, b, a uint32
	calls      int
	panics     bool
}

func (p *grayModelProbe) RGBA() (uint32, uint32, uint32, uint32) {
	p.calls++
	if p.panics {
		panic(p)
	}
	return p.r, p.g, p.b, p.a
}

func grayModelOutcome(tag int, r, g, b, a uint32) (out string) {
	probe := &grayModelProbe{r: r, g: g, b: b, a: a, panics: tag == 2}
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
		c = Gray{uint8(r)}
	case 1, 2:
		c = probe
	case 3:
	default:
		panic("unknown gray input")
	}
	result := grayModel(c).(Gray)
	return fmt.Sprintf("%d %d", result.Y, probe.calls)
}

func TestGoLeanGrayModel(t *testing.T) {
	dir := os.Getenv("GOLEAN_GRAYMODEL_DIR")
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
		var red, green, blue, alpha uint32
		n, err := fmt.Fscan(r, &tag, &red, &green, &blue, &alpha)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprintln(w, grayModelOutcome(tag, red, green, blue, alpha))
	}
	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}
}
