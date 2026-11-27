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

func TestGoLeanBuilderCopyCheck(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_STRINGS_COPY_CHECK_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	var count int
	for {
		var mode, capacity int
		var text, want string
		n, err := fmt.Fscan(reader, &mode, &capacity, &text, &want)
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
		b := strings.GoLeanBuilderCopyCheckSetup(mode, capacity, string(payload))
		for attempt := range 2 {
			if got := goLeanCopyCheckResult(b); got != want {
				t.Fatalf("mode %d, capacity %d, payload %q, attempt %d: Go=%q want=%q",
					mode, capacity, payload, attempt, got, want)
			}
		}
		count++
	}
	if count == 0 {
		t.Fatal("no Builder.copyCheck comparison rows")
	}
	t.Logf("private Builder.copyCheck matched %d rows, two calls per row", count)
}

func goLeanCopyCheckResult(b *strings.Builder) (result string) {
	status := "ok"
	var prior *strings.Builder
	if b != nil {
		_, prior, _, _ = strings.GoLeanBuilderCopyCheckObserve(b)
	}
	defer func() {
		if p := recover(); p != nil {
			_, runtimePanic := p.(runtime.Error)
			switch {
			case b == nil && runtimePanic:
				result = "nil"
				return
			case p == "strings: illegal use of non-zero Builder copied by value":
				status = "copied"
			default:
				panic(p)
			}
		}
		id, addr, capacity, contents := strings.GoLeanBuilderCopyCheckObserve(b)
		if prior != nil && addr != prior {
			panic("copyCheck changed an established address")
		}
		payload := make([]string, len(contents))
		for i := range len(contents) {
			payload[i] = strconv.Itoa(int(contents[i]))
		}
		result = fmt.Sprintf("%s:%s:%d:%s", status, id, capacity, strings.Join(payload, ","))
	}()
	strings.GoLeanBuilderCopyCheck(b)
	return
}
