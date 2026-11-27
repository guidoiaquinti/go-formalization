package rand_test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
	"testing"
)

var globalBoundedCalls = []func(uint64) uint64{
	func(n uint64) uint64 { return uint64(rand.Int64N(int64(n))) },
	rand.Uint64N,
	func(n uint64) uint64 { return uint64(rand.Int32N(int32(n))) },
	func(n uint64) uint64 { return uint64(rand.Uint32N(uint32(n))) },
	func(n uint64) uint64 { return uint64(rand.IntN(int(n))) },
	func(n uint64) uint64 { return uint64(rand.UintN(uint(n))) },
	func(n uint64) uint64 { return uint64(rand.N(signedByte(n))) },
	func(n uint64) uint64 { return uint64(rand.N(uint16(n))) },
}

func nativeGlobalBounded(call func(uint64) uint64, c boundedFixture) (wire string) {
	src := &boundedSource{words: c.Words}
	defer func() {
		if p := recover(); p != nil {
			status := fmt.Sprint(p)
			if p == exhausted {
				status = "callback"
			}
			wire = fmt.Sprintf("%s:0:%d", status, src.pulls)
		}
	}()
	var result uint64
	rand.GoLeanWithGlobalSource(src, func() {
		result = call(c.Bound)
	})
	return fmt.Sprintf("ok:%d:%d", result, src.pulls)
}

func TestGoLeanGlobalBounded(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_RAND_GLOBAL_BOUNDED_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewScanner(file)
	count := 0
	for reader.Scan() {
		var c boundedFixture
		if err := json.Unmarshal(reader.Bytes(), &c); err != nil {
			t.Fatal(err)
		}
		rows := make([]string, len(globalBoundedCalls))
		for i, call := range globalBoundedCalls {
			rows[i] = nativeGlobalBounded(call, c)
			if want := oracleBounded(boundedMethods[i+2], c); rows[i] != want {
				t.Fatalf("fixture %+v method %d: Go %q, oracle %q", c, i, rows[i], want)
			}
		}
		if got := strings.Join(rows, "|"); got != c.Want {
			t.Fatalf("fixture %+v: Go %q", c, got)
		}
		count++
	}
	if err := reader.Err(); err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("missing cases")
	}
	t.Logf("%d package-level bounded draws per instantiation matched Lean and independent oracle", count)
}
