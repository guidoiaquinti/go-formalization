package time_test

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"
)

type fixture struct {
	Wall                uint64
	Ext, Duration, Half int64
	Loc                 int
	Lean                string
}

func quotient(wall uint64, ext, d int64) (int64, int64) {
	nanos := int64(wall % 1073741824)
	sec := ext
	if wall >= 1<<63 {
		sec = 59453308800 + int64(wall/1073741824%(1<<33))
	}
	negative := sec < 0
	if negative {
		sec = -sec
		if nanos != 0 {
			sec--
			nanos = 1000000000 - nanos
		}
	}
	var q, r int64
	switch {
	case d < 1000000000 && 1000000000%(2*d) == 0:
		q, r = nanos/d, nanos%d
	case d%1000000000 == 0:
		unit := d / 1000000000
		q, r = sec/unit, sec%unit*1000000000+nanos
	default:
		n := new(big.Int).SetUint64(uint64(sec))
		n.Mul(n, big.NewInt(1000000000))
		n.Add(n, new(big.Int).SetUint64(uint64(nanos)))
		var whole, remainder big.Int
		whole.QuoRem(n, big.NewInt(d), &remainder)
		q, r = int64(whole.Bit(0)), remainder.Int64()
	}
	q &= 1
	if negative && r != 0 {
		q, r = q^1, d-r
	}
	return q, r
}

func TestGoLeanClockRound(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_TIME_CLOCK_ROUND"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	in := json.NewDecoder(f)
	locations := []*time.Location{nil, time.UTC, time.Local, time.FixedZone("test", 1530)}
	wire := func(value time.Time) string {
		wall, ext, p := time.GoLeanArithmeticFields(value)
		for id, loc := range locations {
			if p == loc {
				return fmt.Sprintf("%d,%d,%d", wall, ext, id)
			}
		}
		t.Fatal("location pointer changed")
		return ""
	}
	for row := 0; ; row++ {
		var c fixture
		err := in.Decode(&c)
		if err == io.EOF && row > 0 {
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		value := time.GoLeanArithmeticTime(c.Wall, c.Ext, locations[c.Loc])
		stripped := time.GoLeanStripMono(value)
		d := time.Duration(c.Duration)
		trunc, round := stripped, stripped
		divResult := "skip"
		if d > 0 {
			q, r := quotient(c.Wall, c.Ext, c.Duration)
			gotQ, gotR := time.GoLeanClockDiv(value, d)
			if gotQ != int(q) || gotR != time.Duration(r) {
				t.Fatalf("row %d: native quotient %d,%d want %d,%d", row, gotQ, gotR, q, r)
			}
			divResult = fmt.Sprintf("%d,%d", q, r)
			trunc = stripped.Add(time.Duration(-r))
			delta := c.Duration - r
			if uint64(r)*2 < uint64(c.Duration) {
				delta = -r
			}
			round = stripped.Add(time.Duration(delta))
		}
		half := uint64(c.Half)*2 < uint64(c.Duration)
		got := strings.Join([]string{wire(value.Truncate(d)), wire(value.Round(d)), divResult, fmt.Sprint(time.GoLeanLessThanHalf(time.Duration(c.Half), d))}, "|")
		want := strings.Join([]string{wire(trunc), wire(round), divResult, fmt.Sprint(half)}, "|")
		if got != want || c.Lean != want {
			t.Fatalf("row %d: Go %q Lean %q want %q", row, got, c.Lean, want)
		}
	}
}
