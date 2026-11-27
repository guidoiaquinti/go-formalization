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

func TestGoLeanBuilderGrowPrivate(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_STRINGS_GROW_PRIVATE_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	var count int
	for {
		var mode, capacity, request int
		var text, want string
		n, err := fmt.Fscan(reader, &mode, &capacity, &request, &text, &want)
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
		if got := goLeanPrivateGrowResult(b, request); got != want {
			t.Fatalf("mode %d, cap %d, request %d, payload %q: Go=%q want=%q",
				mode, capacity, request, payload, got, want)
		}
		count++
	}
	if count == 0 {
		t.Fatal("no private Builder.grow comparison rows")
	}
	t.Logf("private Builder.grow matched %d rows", count)
}

func goLeanPrivateGrowResult(b *strings.Builder, request int) (result string) {
	status := "ok"
	var prior *strings.Builder
	if b != nil {
		_, prior, _, _ = strings.GoLeanBuilderCopyCheckObserve(b)
	}
	defer func() {
		if p := recover(); p != nil {
			e, ok := p.(runtime.Error)
			switch {
			case b == nil && ok:
				result = "nil"
				return
			case ok && strings.HasPrefix(e.Error(), "runtime error: makeslice: len out of range"):
				status = "allocation"
			case ok && strings.HasPrefix(e.Error(), "runtime error: slice bounds out of range"):
				status = "bounds"
			default:
				panic(p)
			}
		}
		id, addr, capacity, contents := strings.GoLeanBuilderCopyCheckObserve(b)
		if addr != prior {
			panic("grow changed an established address")
		}
		payload := make([]string, len(contents))
		for i := range len(contents) {
			payload[i] = strconv.Itoa(int(contents[i]))
		}
		result = fmt.Sprintf("%s:%s:%d:%s", status, id, capacity, strings.Join(payload, ","))
	}()
	strings.GoLeanBuilderGrowPrivate(b, request)
	return
}
