package main

import (
	"bytes"
	"fmt"
	"log"
	"math/big"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	input32, input64 string
	value            string
	raw              uint64
	want32, want64   string
}

func main() {
	cases := []fixture{
		{"32 1 1 9223372036854775808 0 2147483648", "64 1 1 9223372036854775808 9223372036854775808",
			"-9223372036854775808", 1 << 63,
			"true:0,2147483648|true:0,2147483648|false:0,2147483648|-1|-9223372036854775808|9223372036854775808|true|false|64",
			"true:9223372036854775808|true:9223372036854775808|false:9223372036854775808|-1|-9223372036854775808|9223372036854775808|true|false|64"},
		{"32 1 0 18446744073709551615 0 2147483648", "64 1 0 18446744073709551615 9223372036854775808",
			"9223372036854775808", ^uint64(0),
			"true:1|true:1|false:4294967295,4294967295|1|-9223372036854775808|9223372036854775808|false|true|64",
			"true:1|true:1|false:18446744073709551615|1|-9223372036854775808|9223372036854775808|false|true|64"},
		{"32 1 1 4294967296 1 0 1", "64 1 1 4294967296 1 1",
			"-18446744073709551617", 1 << 32,
			"false:0,1|false:0,1|false:0,1|-1|-1|1|false|false|65",
			"false:4294967296|false:4294967296|false:4294967296|-1|-1|1|false|false|65"},
		{"32 1 0 0", "64 1 0 0", "0", 0,
			"false:|false:|false:|0|0|0|true|true|0",
			"false:|false:|false:|0|0|0|true|true|0"},
		{"32 0 0 0", "64 0 0 0", "", 0,
			"false:|panic|panic|panic|panic|panic|panic|panic|panic",
			"false:|panic|panic|panic|panic|panic|panic|panic|panic"},
	}
	var actual64 []string
	for _, width := range []int{32, 64} {
		var input bytes.Buffer
		wants := make([]string, 0, len(cases))
		for _, c := range cases {
			row, want := c.input32, c.want32
			if width == 64 {
				row, want = c.input64, c.want64
			}
			fmt.Fprintln(&input, row)
			wants = append(wants, want)
			if width == strconv.IntSize {
				if got := direct(c); got != want {
					log.Fatalf("native width %d for %s got %q, want %q", width, row, got, want)
				}
			}
		}
		actual := check.LeanRows("godiff-math-big-int-words", len(cases), &input)
		if width == 64 {
			actual64 = actual
		}
		for i, result := range actual {
			if result != wants[i] {
				log.Fatalf("Lean width %d at row %d got %q, want %q", width, i, result, wants[i])
			}
		}
	}

	mutations := []struct {
		row, col int
		output   string
	}{
		{0, 0, "false:9223372036854775808"}, // discard constructor sign
		{3, 1, "true:"},                     // retain prior setter sign for zero
		{1, 2, "true:18446744073709551615"}, // retain prior unsigned setter sign
		{3, 3, "1"},                         // treat zero as positive
		{2, 4, "1"},                         // omit negation on signed conversion
		{2, 5, "18446744073709551615"},      // apply sign on unsigned conversion
		{0, 6, "false"},                     // reject negative signed boundary
		{1, 6, "true"},                      // accept positive signed overflow
		{0, 7, "true"},                      // ignore negative sign on unsigned eligibility
		{2, 8, "1"},                         // count only the top limb
	}
	for _, m := range mutations {
		if strings.Split(actual64[m.row], "|")[m.col] == m.output {
			log.Fatalf("row %d column %d accepted wrong output %q", m.row, m.col, m.output)
		}
	}
	log.Print("five native and ten Lean cases passed and rejected ten wrong observations")
}

func direct(c fixture) string {
	show := func(value *big.Int) string {
		var words []string
		for _, word := range value.Bits() {
			words = append(words, strconv.FormatUint(uint64(word), 10))
		}
		return fmt.Sprintf("%t:%s", value.Sign() < 0, strings.Join(words, ","))
	}
	fresh := big.NewInt(int64(c.raw))
	other := big.NewInt(int64(c.raw))
	if fresh == other {
		log.Fatal("NewInt reused another pointer")
	}
	if words, otherWords := fresh.Bits(), other.Bits(); len(words) > 0 &&
		len(otherWords) > 0 && &words[0] == &otherWords[0] {
		log.Fatal("NewInt shared backing words")
	}
	observe := func(f func() any) (result string) {
		defer func() {
			if p := recover(); p != nil {
				if err, ok := p.(runtime.Error); !ok || err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
					panic(p)
				}
				result = "panic"
			}
		}()
		return fmt.Sprint(f())
	}
	var x, signed, unsigned *big.Int
	if c.value != "" {
		var ok bool
		x, ok = new(big.Int).SetString(c.value, 10)
		if !ok {
			log.Fatalf("invalid fixture integer %q", c.value)
		}
		signed, unsigned = big.NewInt(-7), big.NewInt(-7)
	}
	return strings.Join([]string{show(fresh),
		observe(func() any {
			if signed.SetInt64(int64(c.raw)) != signed {
				log.Fatal("SetInt64 returned a different receiver")
			}
			return show(signed)
		}),
		observe(func() any {
			if unsigned.SetUint64(c.raw) != unsigned {
				log.Fatal("SetUint64 returned a different receiver")
			}
			return show(unsigned)
		}),
		observe(func() any { return x.Sign() }),
		observe(func() any { return x.Int64() }),
		observe(func() any { return x.Uint64() }),
		observe(func() any { return x.IsInt64() }),
		observe(func() any { return x.IsUint64() }),
		observe(func() any { return x.BitLen() }),
	}, "|")
}
