package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	windows              bool
	base, target, result string
	fail                 bool
}

const (
	bareUNC         = `\\host\share`
	rootedUNC       = bareUNC + `\`
	divergenceProbe = "--probe-windows-rel"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == divergenceProbe {
		fmt.Println("ready")
		result, err := filepath.Rel(bareUNC, rootedUNC)
		log.Fatalf("Rel returned unexpectedly: %q, %v", result, err)
	}
	observeDivergence()

	cases := []fixture{
		{false, "/", "/x/../y/", "y", false},
		{false, "a/../.", "./b", "b", false},
		{false, "/a/b/c", "/a/b", "..", false},
		{false, "/a/b/c", "/a/d", "../../d", false},
		{false, "/a//b", "/a/b/", ".", false},
		{false, "../same/a", "../same/b", "../b", false},
		{false, "../../same/a", "../same/b", "", true},
		{false, "x/./a/..", "./x/\x00/../\xff", "\xff", false},
		{false, "/a/../x", "./y//", "", true},
		{false, "../../raw\xff", ".././\x00/", "", true},
		{false, "x/K", "x/k", "../k", false},
		{false, "x/\xff", "x/�", "../�", false},
		{true, "C:/a//b/.", "c:\\A\\d", "..\\d", false},
		{true, "C:base", "c:base\\x", "x", false},
		{true, "C:\\base", "C:base\\leaf", "", true},
		{true, "C:/base", "D:/BASE/.", "", true},
		{true, "\\\\server\\share", "//SERVER/share/x", "x", false},
		{true, "\\\\server\\share\\x", "//SERVER/share", "", true},
		{true, "\\\\server\\share\\x\\y", "//SERVER/share/x/z", "..\\z", false},
		{true, "\\\\?\\C:\\a", "//?/c:/A/b", "b", false},
		{true, "\\a\\b", "\\A\\b\\c", "c", false},
		{true, "x\\K\\a", "X/k/B", "..\\B", false},
		{true, "x\\\xff", "X\\�\\a", "a", false},
		{true, "x\\\xff\\q", "x\\\x80\\q", ".", false},
		{true, rootedUNC, bareUNC, "", true},
		{true, bareUNC, bareUNC, ".", false},
	}
	var rows bytes.Buffer
	for _, f := range cases {
		profile := 0
		if f.windows {
			profile = 1
		}
		fmt.Fprintf(&rows, "%d %d", profile, len(f.base))
		for _, raw := range []string{f.base, f.target} {
			for i := range len(raw) {
				fmt.Fprintf(&rows, " %d", raw[i])
			}
		}
		fmt.Fprintln(&rows)
	}
	var native, cleanError, foldUnix, unrooted int
	for i, got := range check.LeanRows("godiff-path-filepath-rel", len(cases), &rows) {
		f := cases[i]
		want := wire("ok:", f.result)
		if f.fail {
			want = wire("err:", errorText(f.base, f.target))
		}
		if got != want {
			log.Fatalf("profile Windows=%t Rel(%q,%q): Lean=%s want=%s", f.windows, f.base, f.target, got, want)
		}
		if f.windows == (filepath.Separator == '\\') {
			actual, err := filepath.Rel(f.base, f.target)
			if f.fail {
				if actual != "" || err == nil || err.Error() != errorText(f.base, f.target) {
					log.Fatalf("native Rel(%q,%q)=%q,%v", f.base, f.target, actual, err)
				}
			} else if err != nil || actual != f.result {
				log.Fatalf("native Rel(%q,%q)=%q,%v", f.base, f.target, actual, err)
			}
			native++
		}
		if !f.windows {
			if f.fail && errorText(f.base, f.target) != errorText(path.Clean(f.base), path.Clean(f.target)) {
				cleanError++
			}
			if !f.fail && f.result != "." && strings.EqualFold(f.base, f.target) {
				foldUnix++
			}
			if f.fail && strings.HasPrefix(f.base, "/") != strings.HasPrefix(f.target, "/") {
				unrooted++
			}
		}
	}
	if cleanError != 2 || foldUnix != 2 || unrooted != 1 {
		log.Fatalf("unexpected witness counts: cleaned error=%d Unix folding=%d root parity=%d", cleanError, foldUnix, unrooted)
	}
	fmt.Printf("filepath.Rel: %d native and %d Lean fixed cases; cleaned error paths, Unix folding, and omitted root parity caught on %d/%d/%d rows\n",
		native, len(cases), cleanError, foldUnix, unrooted)
}

func errorText(base, target string) string {
	return "Rel: can't make " + target + " relative to " + base
}

func wire(prefix, value string) string {
	fields := make([]string, len(value))
	for i := range len(value) {
		fields[i] = strconv.Itoa(int(value[i]))
	}
	return prefix + strings.Join(fields, ",")
}

func observeDivergence() {
	// In the pinned Windows source the full cleaned paths differ. Their
	// volumes compare equal, the bare base receives a root separator,
	// and the target keeps its original separator after volume removal.
	// The component loop compares two empty slices before the separators;
	// it then advances both start positions to the end of the strings.
	const base, target = `\`, `\`
	b0, bi, t0, ti := 0, 0, 0, 0
	for iteration := range 4 {
		for bi < len(base) && base[bi] != '\\' {
			bi++
		}
		for ti < len(target) && target[ti] != '\\' {
			ti++
		}
		if !strings.EqualFold(target[t0:ti], base[b0:bi]) {
			log.Fatal("component scan exited unexpectedly")
		}
		if bi < len(base) {
			bi++
		}
		if ti < len(target) {
			ti++
		}
		b0, t0 = bi, ti
		if b0 != len(base) || t0 != len(target) || bi != b0 || ti != t0 {
			log.Fatalf("iteration %d did not reach the exhausted fixed point", iteration)
		}
	}

	var request bytes.Buffer
	fmt.Fprintf(&request, "1 %d", len(bareUNC))
	for _, input := range []string{bareUNC, rootedUNC} {
		for i := range len(input) {
			fmt.Fprintf(&request, " %d", input[i])
		}
	}
	fmt.Fprintln(&request)
	if got := check.LeanRows("godiff-path-filepath-rel", 1, &request)[0]; got != "diverges" {
		log.Fatalf("exhausted UNC scan: Lean=%s, want diverges", got)
	}
	if filepath.Separator == '\\' {
		boundedWindowsProbe()
		fmt.Println("native Windows UNC divergence observed in bounded subprocess")
	}
	fmt.Println("Windows UNC cursor fixed point and Lean divergence observation agree")
}

func boundedWindowsProbe() {
	executable, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command(executable, divergenceProbe)
	cmd.Stderr = os.Stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	ready := make(chan bool, 1)
	go func() {
		scanner := bufio.NewScanner(stdout)
		ready <- scanner.Scan() && scanner.Text() == "ready"
	}()
	select {
	case ok := <-ready:
		if !ok {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
			log.Fatal("Windows subprocess failed to reach Rel")
		}
	case <-time.After(10 * time.Second):
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		log.Fatal("Windows subprocess did not start in time")
	}
	exited := make(chan error, 1)
	go func() { exited <- cmd.Wait() }()
	select {
	case err := <-exited:
		log.Fatalf("Windows Rel subprocess exited unexpectedly: %v", err)
	case <-time.After(300 * time.Millisecond):
		if err := cmd.Process.Kill(); err != nil {
			log.Fatalf("Windows Rel subprocess stopped before cancellation: %v", err)
		}
		<-exited
	}
}
