package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	alphabet := make([]string, 256)
	for i := range alphabet {
		alphabet[i] = strconv.Itoa(i)
	}
	identity := "bytes:" + strings.Join(alphabet, ",")
	alphabet[255] = "120"
	firstWins := "bytes:" + strings.Join(alphabet, ",")
	cases := []struct{ input, want string }{
		{"0 1 0 0", "nil-receiver"},
		{"1 1 0 1 2 1 97 98 88", "nil:true:single:97,98|88"},
		{"1 0 0 1 2 1 97 98 88", "nil:false:single:97,98|88"},
		{"1 1 1 1 2 1 97 98 88", "nil:true:" + identity},
		{"1 0 1 1 2 1 97 98 88", "nil:false:" + identity},
		{"1 1 0 2 1 1 255 120 1 1 255 121", "nil:true:" + firstWins},
	}
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&input, c.input)
	}
	model := exec.Command("lean/.lake/build/bin/godiff-strings-replacer-build-once")
	model.Stdin = &input
	out, err := model.Output()
	if err != nil {
		log.Fatal(err)
	}
	rows := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	if len(rows) != len(cases) {
		log.Fatalf("unexpected row count %d", len(rows))
	}
	for i, c := range cases {
		if rows[i] != c.want {
			log.Fatalf("fixture %d: Lean=%q want %q", i, rows[i], c.want)
		}
	}
	testFile := "difftest/strings/replacer-build-once/review/testdata/build_once_test.go"
	if len(os.Args) == 2 {
		testFile = os.Args[1]
	}
	checkNative(testFile)
	fmt.Printf("buildOnce reviewer control: %d Lean fixtures and private native lifecycle checks passed.\n", len(cases))
}

func checkNative(testFile string) {
	scratch, err := os.MkdirTemp("", "golean-build-once-control-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(scratch)
	replacements := make(map[string]string)
	for _, source := range []string{testFile, filepath.Join(filepath.Dir(testFile), "bridge_test.go")} {
		path, err := filepath.Abs(source)
		if err != nil {
			log.Fatal(err)
		}
		replacements[filepath.Join(runtime.GOROOT(), "src", "strings", "golean_build_once_control_"+filepath.Base(path))] = path
	}
	overlay, err := json.Marshal(struct{ Replace map[string]string }{replacements})
	if err != nil {
		log.Fatal(err)
	}
	file := filepath.Join(scratch, "overlay.json")
	if err := os.WriteFile(file, overlay, 0600); err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "test", "-count=1", "-overlay", file, "-run", "^TestGoLeanBuildOnceControl$", "-v", "strings")
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN=local", "GOROOT="+runtime.GOROOT(), "GODEBUG=goindex=0",
		"GOMODCACHE="+filepath.Join(scratch, "modcache"), "GOOS="+runtime.GOOS, "GOARCH="+runtime.GOARCH)
	out, err := cmd.CombinedOutput()
	if err != nil || !bytes.Contains(out, []byte("--- PASS: TestGoLeanBuildOnceControl")) {
		log.Fatalf("private native control: %v\n%s", err, out)
	}
}
