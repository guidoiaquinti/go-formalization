package strconv_test

import (
	"bufio"
	"encoding/json"
	"fmt"
	inner "internal/strconv"
	"math/big"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
)

type formatAPI struct {
	uintText   func(uint64, int) string
	intText    func(int64, int) string
	itoa       func(int) string
	uintAppend func([]byte, uint64, int) []byte
	intAppend  func([]byte, int64, int) []byte
}

func formatWire(raw []byte) string {
	if raw == nil {
		return "nil"
	}
	fields := make([]string, len(raw))
	for i, v := range raw {
		fields[i] = strconv.Itoa(int(v))
	}
	return "slice:" + strings.Join(fields, ",")
}
func appendWire(panicState bool, output, old []byte, reuse bool) string {
	return fmt.Sprintf("%t/%s/%s/%t", panicState, formatWire(output), formatWire(old), reuse)
}

func observeText(call func() string) (out string) {
	defer func() {
		if p := recover(); p != nil {
			if p != "strconv: illegal AppendInt/FormatInt base" {
				panic(p)
			}
			out = "panic"
		}
	}()
	return formatWire([]byte(call()))
}
func observeAppend(cells []byte, lenVisible int, call func([]byte) []byte) (observed string) {
	old := slices.Clone(cells)
	defer func() {
		if p := recover(); p != nil {
			if p != "strconv: illegal AppendInt/FormatInt base" {
				panic(p)
			}
			observed = appendWire(true, nil, old, false)
		}
	}()
	out := call(old[:lenVisible:len(old)])
	return appendWire(false, out, old, len(old) > 0 && &out[0] == &old[0])
}
func oracleAppend(cells []byte, n int, text string, valid, negative bool) string {
	old := slices.Clone(cells)
	output := slices.Clone(cells[:n])
	if negative {
		if n < len(old) {
			old[n] = '-'
		}
		output = append(output, '-')
	}
	if !valid {
		return appendWire(true, nil, old, false)
	}
	if negative {
		text = strings.TrimPrefix(text, "-")
	}
	reuse := len(text) <= len(old)-len(output)
	if reuse {
		copy(old[len(output):], text)
	}
	output = append(output, text...)
	return appendWire(false, output, old, reuse)
}
func TestGoLeanIntegerFormat(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_STRCONV_FORMAT_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	rows := bufio.NewScanner(f)
	count := 0
	internal := formatAPI{inner.FormatUint, inner.FormatInt, inner.Itoa, inner.AppendUint, inner.AppendInt}
	public := formatAPI{strconv.FormatUint, strconv.FormatInt, strconv.Itoa, strconv.AppendUint, strconv.AppendInt}
	for rows.Scan() {
		var c struct {
			Word         uint64
			Base, Length int
			Cells        []byte
			Want         string
		}
		if err := json.Unmarshal(rows.Bytes(), &c); err != nil {
			t.Fatal(err)
		}
		valid := c.Base >= 2 && c.Base <= 36
		ut, it, uw, iw := "", "", "panic", "panic"
		if valid {
			ut = new(big.Int).SetUint64(c.Word).Text(c.Base)
			it = big.NewInt(int64(c.Word)).Text(c.Base)
			uw, iw = formatWire([]byte(ut)), formatWire([]byte(it))
		}
		expect := []string{uw, iw, formatWire([]byte(big.NewInt(int64(int(c.Word))).Text(10))),
			oracleAppend(c.Cells, c.Length, ut, valid, false), oracleAppend(c.Cells, c.Length, it, valid, int64(c.Word) < 0)}
		var all []string
		for _, api := range []formatAPI{internal, public} {
			got := []string{
				observeText(func() string { return api.uintText(c.Word, c.Base) }), observeText(func() string { return api.intText(int64(c.Word), c.Base) }),
				observeText(func() string { return api.itoa(int(c.Word)) }),
				observeAppend(c.Cells, c.Length, func(dst []byte) []byte { return api.uintAppend(dst, c.Word, c.Base) }),
				observeAppend(c.Cells, c.Length, func(dst []byte) []byte { return api.intAppend(dst, int64(c.Word), c.Base) }),
			}
			if !slices.Equal(got, expect) {
				t.Fatalf("%+v: native %v, oracle %v", c, got, expect)
			}
			all = append(all, got...)
		}
		if strings.Join(all, "|") != c.Want {
			t.Fatalf("%+v: native %v", c, all)
		}
		count++
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	t.Logf("integer format: %d pinned observations per declaration", count)
}
