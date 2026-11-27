package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Words                           []int64
	A, B, Limit, Mode, FailAt, Mask int
	Want                            string
}

func wire(status string, values []int64, trace string) string {
	fields := make([]string, len(values))
	for i, value := range values {
		fields[i] = fmt.Sprint(value)
	}
	return status + "/" + strings.Join(fields, ",") + "/" + trace
}

func ascending(n int) []int64 {
	data := make([]int64, n+2)
	data[0], data[n+1] = -91, 92
	for i := 1; i <= n; i++ {
		data[i] = int64(i)
	}
	return data
}

func descending(n int) []int64 {
	data := ascending(n)
	for i := 1; i <= n; i++ {
		data[i] = int64(n - i + 1)
	}
	return data
}

func main() {
	reversedTrace := "38,39;37,39;37,38;26,27;25,27;25,26;14,15;13,15;13,14;26,38;14,38;14,26;2,1"
	cases := []fixture{
		{descending(12), 1, 13, 0, 1, 0, 0, wire("ok", ascending(12), "") + "|" + wire("panic:0", descending(12), "11,12")},
		{descending(13), 1, 14, 0, 1, 1, 0, wire("ok", ascending(13), "") + "|" + wire("panic:1", descending(13), "2,1;8,2")},
		{ascending(13), 1, 14, 1, 1, 3, 0, wire("ok", ascending(13), "") + "|" + wire("panic:3", ascending(13), "7,4;10,7;7,4;2,1")},
		{descending(50), 1, 51, 1, 1, 12, 0, wire("ok", ascending(50), "") + "|" + wire("panic:12", ascending(50), reversedTrace)},
		{[]int64{-91, 92}, 1, 1, 0, 0, -1, 0, "ok/-91,92/|ok/-91,92/"},
	}
	short := cases[0]
	root := (short.B - short.A - 1) / 2
	heapFirst := fmt.Sprintf("%d,%d", short.Words[short.A+root], short.Words[short.A+2*root+1])
	if heapFirst == "11,12" {
		panic("heap-before-insertion mutation survived")
	}
	wide := cases[1]
	q := (wide.B - wide.A) / 4
	selectFirst := fmt.Sprintf("%d,%d", wide.Words[wide.A+2*q], wide.Words[wide.A+q])
	if selectFirst == "2,1" {
		panic("zero-budget pivot selection mutation survived")
	}
	noReverse := wire("panic:12", descending(50), reversedTrace)
	if strings.HasSuffix(cases[3].Want, noReverse) {
		panic("decreasing-hint no-reversal mutation survived")
	}
	fmt.Printf("rejected heap-first operand pair %s, zero-budget pivot pair %s and omitted reversal\n", heapFirst, selectFirst)

	f, err := os.CreateTemp("", "review-pdqsort-")
	must(err)
	defer os.Remove(f.Name())
	encoder := json.NewEncoder(f)
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&input, strconv.IntSize, " ", c.A, " ", c.B, " ", c.Limit, " ", c.Mode, " ", c.FailAt, " ", c.Mask)
		for _, word := range c.Words {
			fmt.Fprint(&input, " ", word)
		}
		fmt.Fprintln(&input)
		must(encoder.Encode(c))
	}
	must(f.Close())
	for i, row := range check.LeanRows("godiff-slices-pdqsort", len(cases), &input) {
		if row != cases[i].Want {
			panic(fmt.Sprintf("fixture %d got %q want %q", i, row, cases[i].Want))
		}
	}
	must(os.Setenv("GOLEAN_SLICES_PDQSORT_CASES", f.Name()))
	const data = "difftest/slices/pdqsort/testdata/"
	must(check.TestInPackage("slices", "TestGoLeanPdqsort", data+"bridge_test.go", data+"pdqsort_test.go"))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
