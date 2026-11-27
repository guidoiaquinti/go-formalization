package strings_test

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanLastIndexFunc(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_STRINGS_LAST_INDEX_FUNC_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	var count int
	for {
		var truth bool
		var kind int
		var text, want string
		n, err := fmt.Fscan(reader, &truth, &kind, &text, &want)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		payload, err := hex.DecodeString(text[1:])
		if err != nil {
			t.Fatal(err)
		}
		if got := goLeanLastIndexFuncResult(string(payload), kind, truth, false); got != want {
			t.Fatalf("input %q, predicate %d, truth %t: Go=%q want=%q", payload, kind, truth, got, want)
		}
		if truth {
			if got := goLeanLastIndexFuncResult(string(payload), kind, truth, true); got != want {
				t.Fatalf("input %q, predicate %d: public Go=%q want=%q", payload, kind, got, want)
			}
		}
		count++
	}
	if count == 0 {
		t.Fatal("no private lastIndexFunc comparison rows")
	}
	t.Logf("private lastIndexFunc matched %d comparison rows", count)
}

func goLeanLastIndexFuncResult(s string, kind int, truth, public bool) (result string) {
	var calls []rune
	var predicate func(rune) bool
	if kind != 0 {
		predicate = func(r rune) bool {
			calls = append(calls, r)
			switch kind {
			case 1:
				return false
			case 2:
				return true
			case 3:
				return r == '\uFFFD'
			case 4:
				return r == '🚀'
			default:
				return r%7 == 2
			}
		}
	}
	defer func() {
		if p := recover(); p != nil {
			if _, ok := p.(runtime.Error); ok && kind == 0 {
				result = "panic"
			} else {
				panic(p)
			}
		}
	}()
	var index int
	if public {
		index = strings.LastIndexFunc(s, predicate)
	} else {
		index = strings.GoLeanLastIndexFunc(s, predicate, truth)
	}
	row := strconv.AppendInt(nil, int64(index), 10)
	row = append(row, '|')
	for i, r := range calls {
		if i != 0 {
			row = append(row, ',')
		}
		row = strconv.AppendInt(row, int64(r), 10)
	}
	return string(row)
}
