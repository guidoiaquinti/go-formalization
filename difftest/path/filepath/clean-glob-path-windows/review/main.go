package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type answer struct {
	prefix  int
	cleaned string
}

func wire(out answer) string {
	var fields []string
	for i := range len(out.cleaned) {
		fields = append(fields, strconv.Itoa(int(out.cleaned[i])))
	}
	return fmt.Sprintf("%d:%s", out.prefix, strings.Join(fields, ","))
}

func main() {
	fixtures := []struct {
		path          string
		volume        int
		unix, windows answer
	}{
		{"", 0, answer{0, "."}, answer{0, "."}},
		{"/", 0, answer{1, "/"}, answer{1, "/"}},
		{"\\", 0, answer{0, ""}, answer{1, "\\"}},
		{"C:", 2, answer{0, "C"}, answer{2, "C:."}},
		{"C:/", 2, answer{0, "C:"}, answer{3, "C:/"}},
		{"C:\\leaf\\", 2, answer{0, "C:\\leaf"}, answer{2, "C:\\leaf"}},
		{"//h/s", 5, answer{0, "//h/"}, answer{4, "//h/"}},
		{"\xffé", 0, answer{0, "\xff\xc3"}, answer{0, "\xff\xc3"}},
	}
	var rows, native bytes.Buffer
	noSpecial, noClamp := 0, 0
	for _, f := range fixtures {
		fmt.Fprint(&rows, "0")
		for i := range len(f.path) {
			fmt.Fprint(&rows, " ", f.path[i])
		}
		fmt.Fprintln(&rows)
		want := f.unix
		if runtime.GOOS == "windows" {
			want = f.windows
		}
		fmt.Fprintf(&native, "x%x %d x%x\n", f.path, want.prefix, want.cleaned)
		if f.path != "" {
			end := len(f.path) - 1
			if (answer{min(f.volume, end), f.path[:end]}) != f.windows {
				noSpecial++
			}
			if f.volume == len(f.path) && len(f.path) > 2 && f.windows.prefix != f.volume {
				noClamp++
			}
		}
	}
	for i, got := range check.LeanRows("godiff-path-filepath-clean-glob-path-windows", len(fixtures), &rows) {
		f := fixtures[i]
		if want := wire(f.unix) + "|" + wire(f.windows); got != want {
			log.Fatalf("profiles for %x got %q, want %q", f.path, got, want)
		}
	}
	if noSpecial != 4 || noClamp != 1 {
		log.Fatalf("mutation failures %d/%d, want 4/1", noSpecial, noClamp)
	}
	dir, err := os.MkdirTemp("", "golean-review-clean-glob-windows-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, native.Bytes(), 0600); err != nil {
		log.Fatal(err)
	}
	if err := os.Setenv("GOLEAN_FILEPATH_CLEAN_GLOB_WINDOWS_CASES", path); err != nil {
		log.Fatal(err)
	}
	const testdata = "difftest/path/filepath/clean-glob-path-windows/testdata/"
	if err := check.TestInPackage("path/filepath", "TestGoLeanCleanGlobWindows", testdata+"bridge_test.go", testdata+"clean_glob_windows_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("cleanGlobPathWindows passed %d pinned native and Lean profile fixtures, always-chop and unclamped-volume mutants fail %d/%d\n",
		len(fixtures), noSpecial, noClamp)
}
