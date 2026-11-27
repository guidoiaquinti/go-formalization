package bytealg_test

import (
	"encoding/json"
	"internal/bytealg"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanRabinKarp(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_RABIN_KARP_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	for i := 0; ; i++ {
		var row struct {
			Input, Separator []byte
			Hashes           [4]uint32
			First, Last      string
		}
		if err := decoder.Decode(&row); err == io.EOF {
			if i == 0 {
				t.Fatal("missing Rabin-Karp cases")
			}
			t.Logf("%d inputs matched both slice and string instantiations of all four declarations", i)
			break
		} else if err != nil {
			t.Fatal(err)
		}
		s, sep := string(row.Input), string(row.Separator)
		h := row.Hashes
		want := [8]uint32{h[0], h[1], h[0], h[1], h[2], h[3], h[2], h[3]}
		if got := bytealg.GoLeanRabinKarpHashes(row.Separator, sep); got != want {
			t.Fatalf("row %d separator=%x: hash got %v want %v", i, row.Separator, got, want)
		}
		for j, call := range []func() int{
			func() int { return bytealg.GoLeanRabinKarpFirstSlice(row.Input, row.Separator) },
			func() int { return bytealg.GoLeanRabinKarpFirstString(s, sep) },
			func() int { return bytealg.GoLeanRabinKarpLastSlice(row.Input, row.Separator) },
			func() int { return bytealg.GoLeanRabinKarpLastString(s, sep) },
		} {
			expect := row.First
			if j >= 2 {
				expect = row.Last
			}
			if got := observe(call); got != expect {
				t.Fatalf("row %d function %d input=%x sep=%x: got %s want %s", i, j, row.Input, row.Separator, got, expect)
			}
		}
	}
}

func observe(call func() int) (out string) {
	defer func() {
		if value := recover(); value != nil {
			err, ok := value.(runtime.Error)
			if !ok || !strings.HasPrefix(err.Error(), "runtime error: index out of range") {
				panic(value)
			}
			out = "panic"
		}
	}()
	return strconv.Itoa(call())
}
