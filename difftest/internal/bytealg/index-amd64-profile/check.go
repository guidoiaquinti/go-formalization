package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	if strconv.IntSize != 64 {
		fmt.Println("64-bit signed arithmetic is required for the amd64 source-body comparison.")
		return
	}
	dir, err := os.MkdirTemp("", "golean-bytealg-amd64-profile-")
	must(err)
	defer os.RemoveAll(dir)
	pinned := filepath.Join(runtime.GOROOT(), "src", "internal", "bytealg", "index_amd64.go")
	extracted := filepath.Join(dir, "index_profile_test.go")
	must(os.WriteFile(extracted, sourceFixture(pinned), 0600))
	file, err := os.Create(filepath.Join(dir, "cases"))
	must(err)
	writer := bufio.NewWriter(file)
	var edges []int64
	for n := int64(-1024); n < 1024; n++ {
		edges = append(edges, n)
	}
	for step := range int64(32) {
		edges = append(edges, math.MinInt64+step, math.MaxInt64-step)
	}
	compare("boundaries", edges, writer)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]int64, 20_000)
		for i := range cases {
			cases[i] = int64(rng.Uint64())
		}
		compare(fmt.Sprintf("seed=%d", seed), cases, writer)
	}
	must(writer.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_BYTEALG_AMD64_PROFILE_CASES", file.Name()))
	const observer = "difftest/internal/bytealg/index-amd64-profile/testdata/profile_test.go"
	must(check.TestInPackage("internal/bytealg", "TestGoLeanBytealgIndexAMD64Profile", extracted, observer))
}

func sourceFixture(path string) []byte {
	src, err := os.ReadFile(path)
	must(err)
	set := token.NewFileSet()
	file, err := parser.ParseFile(set, path, src, 0)
	must(err)
	body := func(name string) string {
		for _, decl := range file.Decls {
			f, ok := decl.(*ast.FuncDecl)
			if ok && f.Name.Name == name && f.Recv == nil && f.Body != nil {
				return string(src[set.Position(f.Body.Lbrace).Offset+1 : set.Position(f.Body.Rbrace).Offset])
			}
		}
		log.Fatalf("missing pinned function %s", name)
		return ""
	}
	return fmt.Appendf(nil, `package bytealg
func GoLeanPinnedCutover(n int) int {%s}
func GoLeanPinnedIndexInit(avx2 bool) int {
	var cpu struct { X86 struct { HasAVX2 bool } }
	cpu.X86.HasAVX2 = avx2
	var MaxLen int
	%s
	return MaxLen
}
`, body("Cutover"), body("init"))
}

func compare(label string, cases []int64, writer *bufio.Writer) {
	var input bytes.Buffer
	for _, n := range cases {
		fmt.Fprintln(&input, n)
	}
	for i, got := range check.LeanRows("godiff-internal-bytealg-index-amd64-profile", len(cases), &input) {
		n := cases[i]
		// Signed shift rounds down. Correct its negative remainder to get
		// the source's truncation toward zero independently of division.
		adjusted := n + 16
		want := adjusted >> 3
		if adjusted < 0 && adjusted&7 != 0 {
			want++
		}
		if expected := fmt.Sprintf("%d,31,63", want); got != expected {
			log.Fatalf("%s n=%d: Lean=%s want=%s", label, n, got, expected)
		}
		_, err := fmt.Fprintln(writer, n, want)
		must(err)
	}
	fmt.Printf("amd64 index profile %s: %d Lean rows agree.\n", label, len(cases))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
