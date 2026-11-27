package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	profile                    int
	failure                    bool
	path, query, operand, want string
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--deleted-cwd" {
		must(os.Chdir(os.Args[2]))
		must(os.Remove(os.Args[2]))
		must(os.Unsetenv("PWD"))
		wd, wdErr := os.Getwd()
		value, err := filepath.Abs("./relative")
		if wdErr == nil {
			if value != filepath.Join(wd, "relative") || err != nil {
				log.Fatalf("relative Abs=%q,%v; Getwd=%q", value, err, wd)
			}
		} else if value != "" || err == nil || err.Error() != wdErr.Error() {
			log.Fatalf("relative Abs=%q,%v; Getwd=%v", value, err, wdErr)
		}
		value, err = filepath.Abs("/abs/../leaf")
		if value != "/leaf" || err != nil {
			log.Fatalf("absolute Abs=%q,%v with deleted cwd", value, err)
		}
		if wdErr == nil {
			fmt.Print("ok:", wd)
		} else {
			fmt.Print("err:", wdErr.Error())
		}
		return
	}

	cases := []fixture{
		{0, true, "/a//../leaf", "", "getwd failed", wire("ok:", "/leaf")},
		{0, true, "relative", "", "\xff missing", wire("err:", "\xff missing")},
		{0, false, "", "", "/cwd/dir", wire("ok:", "/cwd/dir")},
		{0, false, "../\xff", "", "/cwd/dir", wire("ok:", "/cwd/\xff")},
		{0, false, "#device//x/../leaf", "", "/cwd", wire("ok:", "/cwd/#device/leaf")},
		{2, true, "#device//x/../leaf", "", "getwd failed", wire("ok:", "#device/leaf")},
		{2, true, "./#device", "", "getwd failed", wire("err:", "getwd failed")},
		{2, false, "../leaf", "", "#device/dir", wire("ok:", "#device/leaf")},
		{1, false, "", ".", `C:/work/./dir/..`, wire("ok:", `C:\work`)},
		{1, true, "", ".", "\x00 invalid", wire("err:", "\x00 invalid")},
		{1, true, `C:\ready`, `C:\ready`, "os failure", wire("err:", "os failure")},
		{1, false, `C:relative`, `C:relative`, `C:\base\..\leaf`, wire("ok:", `C:\leaf`)},
		{1, false, `/a/../b`, `/a/../b`, `D:\a\..\b`, wire("ok:", `D:\b`)},
	}
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		dir, err := os.MkdirTemp("", "golean-abs-cwd-")
		must(err)
		defer os.RemoveAll(dir)
		executable, err := os.Executable()
		must(err)
		cmd := exec.Command(executable, "--deleted-cwd", dir)
		cmd.Stderr = os.Stderr
		raw, err := cmd.Output()
		must(err)
		status, operand, ok := strings.Cut(string(raw), ":")
		if !ok || (status != "ok" && status != "err") {
			log.Fatalf("invalid deleted-cwd observation %q", raw)
		}
		failed := status == "err"
		want := wire("ok:", filepath.Join(operand, "relative"))
		if failed {
			want = wire("err:", operand)
		}
		cases = append(cases,
			fixture{0, failed, "./relative", "", operand, want},
			fixture{0, failed, "/abs/../leaf", "", operand, wire("ok:", "/leaf")})
		fmt.Printf("Native deleted-cwd probe matched Getwd outcome %s for relative Abs and cleaned the absolute path\n", status)

	}

	var input bytes.Buffer
	for _, c := range cases {
		failure := 0
		if c.failure {
			failure = 1
		}
		fmt.Fprintf(&input, "%d %d %d %d", c.profile, failure, len(c.path), len(c.query))
		for _, raw := range []string{c.path, c.query, c.operand} {
			for i := range len(raw) {
				fmt.Fprintf(&input, " %d", raw[i])
			}
		}
		fmt.Fprintln(&input)
	}
	for i, got := range check.LeanRows("godiff-path-filepath-abs", len(cases), &input) {
		if got != cases[i].want {
			log.Fatalf("fixture %d: Lean=%s want=%s", i, got, cases[i].want)
		}
	}
	if cases[0].want == wire("err:", cases[0].operand) ||
		cases[10].want == wire("ok:", cases[10].path) ||
		cases[8].want == wire("ok:", cases[8].operand) ||
		cases[4].want == cases[5].want {
		log.Fatal("one incorrect variant escaped")
	}
	fmt.Printf("filepath.Abs: %d Lean fixed cases; four incorrect variants distinguished\n", len(cases))
}

func wire(prefix, raw string) string {
	values := make([]string, len(raw))
	for i := range len(values) {
		values[i] = strconv.Itoa(int(raw[i]))
	}
	return prefix + strings.Join(values, ",")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
