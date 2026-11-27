package main

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	profile  int
	path     string
	positive []string
	want     bool
}

func main() {
	cases := []fixture{
		{0, "", nil, false},
		{0, "#dev/leaf", nil, true},
		{0, "./a/../leaf", nil, true},
		{0, "\\..\\leaf", nil, true},
		{0, "a/../..//\x00", nil, false},
		{0, "a/../\xff", nil, true},
		{0, "/a/../b", nil, false},
		{0, "..leaf", nil, true},
		{2, "#", nil, false},
		{2, "#dev/leaf", nil, false},
		{2, "./#dev", nil, true},
		{2, "a/../#dev", nil, true},
		{2, "a/../../#dev", nil, false},
		{1, "a\\../leaf", nil, true},
		{1, "a/..\\..//leaf", nil, false},
		{1, "first/second:leaf", nil, false},
		{1, "first\\cOm³/..", nil, false},
		{1, "first\\cOm⁴/..", nil, true},
		{1, "CoNoUt$", nil, false},
		{1, "CoNoUt$foo", nil, true},
		{1, "dir//PrN .log/..", nil, true},
		{1, "dir//PrN .log/..", []string{"PrN .log"}, false},
		{1, "dir//PrN .log/..", []string{"PrN", "prn .log"}, true},
		{1, "dir//PrN .log\x00/..", []string{"PrN .log\x00"}, true},
		{1, "dir//NUL/..", nil, false},
		{1, "dir//ordinary/..", []string{"ordinary"}, true},
		{1, "dir//LPT¹", nil, false},
		{1, "dir//LPT².ext", []string{"LPT².ext"}, false},
		{1, "dir//LPT².ext", nil, true},
	}
	var request bytes.Buffer
	for _, test := range cases {
		fmt.Fprintf(&request, "%d %d", test.profile, len(test.positive))
		for _, part := range test.positive {
			fmt.Fprintf(&request, " %d", len(part))
			writeBytes(&request, part)
		}
		writeBytes(&request, test.path)
		fmt.Fprintln(&request)
	}
	nativeProfile := 0
	switch runtime.GOOS {
	case "plan9":
		nativeProfile = 2
	case "windows":
		nativeProfile = 1
	}
	native := 0
	for i, got := range check.LeanRows("godiff-path-filepath-is-local", len(cases), &request) {
		test := cases[i]
		if got != strconv.FormatBool(test.want) {
			log.Fatalf("profile %d path %q devices %q: Lean=%s want=%t",
				test.profile, test.path, test.positive, got, test.want)
		}
		if test.profile == nativeProfile && test.profile != 1 {
			if got := filepath.IsLocal(test.path); got != test.want {
				log.Fatalf("native %q: got %t want %t", test.path, got, test.want)
			}
			native++
		}
	}
	distinct := func(a, b int) {
		if cases[a].want == cases[b].want {
			log.Fatalf("incorrect variant escaped fixtures %d and %d", a, b)
		}
	}
	// Each pair changes only a relevant input detail.
	distinct(21, 20) // Ignore the device callback for extended names.
	distinct(21, 22) // Trim or fold the callback argument.
	distinct(21, 23) // Consult the callback despite an embedded NUL.
	distinct(16, 17) // Reject all or none of the superscripts.
	distinct(1, 9)   // Share one Unix and Plan 9 absolute-path rule.
	fmt.Printf("filepath.IsLocal: %d native checks and %d Lean profile fixtures, five incorrect variants distinguished\n", native, len(cases))
}

func writeBytes(out *bytes.Buffer, raw string) {
	for i := range len(raw) {
		fmt.Fprintf(out, " %d", raw[i])
	}
}
