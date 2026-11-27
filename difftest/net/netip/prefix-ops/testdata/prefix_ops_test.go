package netip

import (
	"bufio"
	"bytes"
	"cmp"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

type goLeanPrefixIP struct {
	Family int
	Hi, Lo uint64
	Zone   bool
}

func (a goLeanPrefixIP) cells() [16]byte {
	var cells [16]byte
	binary.BigEndian.PutUint64(cells[:8], a.Hi)
	binary.BigEndian.PutUint64(cells[8:], a.Lo)
	return cells
}
func (a goLeanPrefixIP) native() Addr {
	cells := a.cells()
	if a.Family == 4 {
		return AddrFrom4([4]byte(cells[12:]))
	}
	if a.Family == 6 {
		out := AddrFrom16(cells)
		if a.Zone {
			return out.WithZone(string([]byte{255, 0, 37}))
		}
		return out
	}
	return Addr{}
}
func (a goLeanPrefixIP) width() int {
	if a.Family == 4 {
		return 32
	}
	if a.Family == 6 {
		return 128
	}
	return 0
}
func (a goLeanPrefixIP) mask(n int) goLeanPrefixIP {
	if a.Family == 0 || n < 0 || n > a.width() {
		return goLeanPrefixIP{}
	}
	cells := a.cells()
	if a.Family == 4 {
		n += 96
	}
	for i := n; i < 128; i++ {
		cells[i/8] &^= byte(128 >> (i % 8))
	}
	return goLeanPrefixIP{a.Family, binary.BigEndian.Uint64(cells[:8]), binary.BigEndian.Uint64(cells[8:]), false}
}
func (a goLeanPrefixIP) cmp(b goLeanPrefixIP) int {
	if c := cmp.Compare(a.Family, b.Family); c != 0 {
		return c
	}
	x, y := a.cells(), b.cells()
	return bytes.Compare(x[:], y[:])
}
func goLeanPrefixGoWire(p Prefix) string {
	family := 0
	if p.ip.Is4() {
		family = 4
	}
	if p.ip.Is6() {
		family = 6
	}
	return fmt.Sprintf("%d,%d,%d,%d", family, p.ip.addr.hi, p.ip.addr.lo, p.Bits())
}
func goLeanPrefixOracleWire(a goLeanPrefixIP, n int) string {
	if a.Family == 0 || n < 0 || n > a.width() {
		return "0,0,0,-1"
	}
	a = a.mask(n)
	return fmt.Sprintf("%d,%d,%d,%d", a.Family, a.Hi, a.Lo, n)
}

func TestGoLeanPrefixOps(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_PREFIX_OPS_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		var c struct {
			A, B goLeanPrefixIP
			N, M int
			Want string
		}
		if err := json.Unmarshal(scanner.Bytes(), &c); err != nil {
			t.Fatal(err)
		}
		a, b := c.A.native(), c.B.native()
		p, q := PrefixFrom(a, c.N), PrefixFrom(b, c.M)
		result, err := a.Prefix(c.N)
		direct := goLeanPrefixGoWire(result)
		if err != nil {
			direct = "error:" + err.Error()
		}
		if result.ip.Zone() != "" || p.Masked().ip.Zone() != "" {
			t.Fatal("derived prefix retained a zone")
		}
		directWant := goLeanPrefixOracleWire(c.A, c.N)
		if c.N < 0 {
			directWant = "error:negative Prefix bits"
		} else if c.A.Family != 0 && c.N > c.A.width() {
			directWant = fmt.Sprintf("error:prefix length %d too large for IPv%d", c.N, c.A.Family)
		}
		av := c.A.Family != 0 && 0 <= c.N && c.N <= c.A.width()
		bv := c.B.Family != 0 && 0 <= c.M && c.M <= c.B.width()
		same := c.A.Family == c.B.Family
		contains := av && same && !c.B.Zone && c.A.mask(c.N) == c.B.mask(c.N)
		overlaps := av && bv && same && c.A.mask(min(c.N, c.M)) == c.B.mask(min(c.N, c.M))
		comparison := c.A.mask(c.N).cmp(c.B.mask(c.M))
		pb, qb := -1, -1
		if av {
			pb = c.N
		}
		if bv {
			qb = c.M
		}
		if comparison == 0 {
			comparison = cmp.Compare(pb, qb)
		}
		if comparison == 0 {
			comparison = c.A.cmp(c.B)
		}
		want := fmt.Sprintf("%s|%s|%t|%t|%t|%d", directWant, goLeanPrefixOracleWire(c.A, c.N), c.A.Family == 0, contains, overlaps, comparison)
		observed := fmt.Sprintf("%s|%s|%t|%t|%t|%d", direct, goLeanPrefixGoWire(p.Masked()), p.isZero(), p.Contains(b), p.Overlaps(q), p.Compare(q))
		if observed != want || observed != c.Want {
			t.Fatalf("%+v: native %s, oracle %s", c, observed, want)
		}
		if p.Compare(q) != -q.Compare(p) || p.Overlaps(q) != q.Overlaps(p) {
			t.Fatalf("asymmetric comparison or overlap: %s", observed)
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	t.Logf("prefix ops: %d pinned observations per declaration", count)
}
