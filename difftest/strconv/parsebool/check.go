package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func decimalBytes(s string) string {
	fields := make([]string, len(s))
	for i := range len(s) {
		fields[i] = strconv.Itoa(int(s[i]))
	}
	return strings.Join(fields, " ")
}

func main() {
	accepted := map[string]bool{
		"1": true, "t": true, "T": true, "true": true, "TRUE": true, "True": true,
		"0": false, "f": false, "F": false, "false": false, "FALSE": false, "False": false,
	}
	edges := []string{"", " true", "true ", "tRuE", "FaLsE", "false\x00", "tru", "false\n", "\xff"}
	for _, word := range []string{"true", "false"} {
		for mask := range 1 << len(word) {
			spelling := []byte(word)
			for i := range spelling {
				if mask&(1<<i) != 0 {
					spelling[i] -= 'a' - 'A'
				}
			}
			edges = append(edges, string(spelling))
		}
	}
	for b := range 256 {
		edges = append(edges, string([]byte{byte(b)}))
	}
	check.CompareSamples("godiff-parsebool", edges, func(r *rand.Rand, i int) string {
		if i%3 == 0 {
			return edges[r.Intn(len(edges))]
		}
		raw := make([]byte, r.Intn(33))
		r.Read(raw)
		return string(raw)
	}, func(s string) string {
		return "0 " + decimalBytes(s)
	}, func(s string) string {
		value, err := strconv.ParseBool(s)
		if err == nil {
			return fmt.Sprintf("%t ok", value)
		}
		if e, ok := err.(*strconv.NumError); ok && e.Err == strconv.ErrSyntax {
			return fmt.Sprintf("%t syntax %s [%s]", value, e.Func, decimalBytes(e.Num))
		}
		return fmt.Sprintf("%t unexpected error %T: %v", value, err, err)
	}, func(s string) string {
		if value, ok := accepted[s]; ok {
			return fmt.Sprintf("%t ok", value)
		}
		return "false syntax ParseBool [" + decimalBytes(s) + "]"
	})
}
