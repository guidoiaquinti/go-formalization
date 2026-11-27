package main

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	fixtures := []struct {
		path                 string
		unix, windows, plan9 bool
	}{
		{"", false, false, false}, {"/", true, false, true}, {"/root", true, false, true},
		{"\\root", false, false, false}, {"C:/root", false, true, false}, {"C:\\root", false, true, false},
		{"C:relative", false, false, false}, {"C:", false, false, false},
		{"//host/share", true, true, true}, {"\\\\host\\share", false, true, false},
		{"#d/name", false, false, true}, {"#", false, false, true}, {"x#d", false, false, false},
		{"\xff:/\x00", false, true, false},
	}
	var rows bytes.Buffer
	driveRelative, rootWithoutVolume, ignoreDevice := 0, 0, 0
	for _, f := range fixtures {
		fmt.Fprint(&rows, "0")
		for i := range len(f.path) {
			fmt.Fprint(&rows, " ", f.path[i])
		}
		fmt.Fprintln(&rows)
		native := f.unix
		switch runtime.GOOS {
		case "windows":
			native = f.windows
		case "plan9":
			native = f.plan9
		}
		if got := filepath.IsAbs(f.path); got != native {
			log.Fatalf("Go IsAbs(%x)=%t want %t", f.path, got, native)
		}
		if (f.windows || len(f.path) > 1 && f.path[1] == ':') != f.windows {
			driveRelative++
		}
		if (f.windows || strings.HasPrefix(f.path, "/") || strings.HasPrefix(f.path, "\\")) != f.windows {
			rootWithoutVolume++
		}
		if f.unix != f.plan9 {
			ignoreDevice++
		}
	}
	for i, got := range check.LeanRows("godiff-path-filepath-is-abs", len(fixtures), &rows) {
		f := fixtures[i]
		want := strconv.FormatBool(f.unix) + "|" + strconv.FormatBool(f.windows) + "|" + strconv.FormatBool(f.plan9)
		if got != want {
			log.Fatalf("Lean profiles IsAbs(%x)=%s want %s", f.path, got, want)
		}
	}
	if driveRelative != 2 || rootWithoutVolume != 3 || ignoreDevice != 2 {
		log.Fatalf("mutation differences drive=%d root=%d device=%d, want 2/3/2", driveRelative, rootWithoutVolume, ignoreDevice)
	}
	fmt.Println("IsAbs passed 14 fixtures in native Go and all Lean profiles; mutations failed 2/3/2")
}
