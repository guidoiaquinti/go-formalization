package bytealg_test

import (
	"bufio"
	"fmt"
	"internal/bytealg"
	"internal/cpu"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanBytealgIndexAMD64Profile(t *testing.T) {
	a, b := bytealg.GoLeanPinnedIndexInit(false), bytealg.GoLeanPinnedIndexInit(true)
	if a != 31 || b != 63 {
		t.Fatalf("pinned amd64 initializer selected %d and %d", a, b)
	}
	if runtime.GOARCH == "amd64" {
		want := a
		if cpu.X86.HasAVX2 {
			want = b
		}
		if bytealg.MaxLen != want {
			t.Fatalf("native amd64 MaxLen=%d want=%d", bytealg.MaxLen, want)
		}
	}
	file, err := os.Open(os.Getenv("GOLEAN_BYTEALG_AMD64_PROFILE_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows := bufio.NewScanner(file)
	count := 0
	for rows.Scan() {
		left, right, ok := strings.Cut(rows.Text(), " ")
		if !ok {
			t.Fatalf("invalid fixture")
		}
		n, err := strconv.Atoi(left)
		if err != nil {
			t.Fatal(err)
		}
		want, err := strconv.Atoi(right)
		if err != nil {
			t.Fatal(err)
		}
		got := bytealg.GoLeanPinnedCutover(n)
		if got != want {
			t.Fatalf("pinned Cutover(%d)=%d want=%d", n, got, want)
		}
		if runtime.GOARCH == "amd64" && bytealg.Cutover(n) != want {
			t.Fatalf("native Cutover(%d)=%d want=%d", n, bytealg.Cutover(n), want)
		}
		count++
	}
	if err := rows.Err(); err != nil || count == 0 {
		t.Fatalf("read %d cases: %v", count, err)
	}
	fmt.Printf("amd64 source-body check matched %d Cutover rows and both initializer branches; native AMD64=%t\n", count, runtime.GOARCH == "amd64")
}
