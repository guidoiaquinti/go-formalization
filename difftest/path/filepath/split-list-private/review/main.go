package main

import (
	"bytes"
	"encoding/hex"
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
	fixtures := []struct {
		input                string
		unix, windows, plan9 []string
	}{
		{"", []string{}, []string{}, []string{}},
		{";;", []string{";;"}, []string{"", "", ""}, []string{";;"}},
		{"\x00:\x00", []string{"\x00", "\x00"}, []string{"\x00:\x00"}, []string{"", ":", ""}},
		{"\"a;b\"\";c\";", []string{"\"a;b\"\";c\";"}, []string{"a;b;c", ""}, []string{"\"a;b\"\";c\";"}},
		{"\";\";\"", []string{"\";\";\""}, []string{";", ""}, []string{"\";\";\""}},
		{"\\\";\xff\x00:\"", []string{"\\\";\xff\x00", "\""}, []string{"\\;\xff\x00:"}, []string{"\\\";\xff", ":\""}},
		{"\"\xc0\xaf;\xef\xbf\xbd\";\xff", []string{"\"\xc0\xaf;\xef\xbf\xbd\";\xff"}, []string{"\xc0\xaf;\xef\xbf\xbd", "\xff"}, []string{"\"\xc0\xaf;\xef\xbf\xbd\";\xff"}},
	}
	var input, native bytes.Buffer
	for _, row := range fixtures {
		fmt.Fprint(&input, "0")
		for _, b := range []byte(row.input) {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
		values := row.unix
		if filepath.ListSeparator == ';' {
			values = row.windows
		} else if filepath.ListSeparator == 0 {
			values = row.plan9
		}
		fmt.Fprint(&native, "x", hex.EncodeToString([]byte(row.input)))
		for _, part := range values {
			fmt.Fprint(&native, " x", hex.EncodeToString([]byte(part)))
		}
		fmt.Fprintln(&native)
	}
	model := exec.Command("lean/.lake/build/bin/godiff-path-filepath-split-list-private")
	model.Stdin = &input
	output, err := model.Output()
	if err != nil {
		log.Fatal(err)
	}
	rows := strings.Split(strings.TrimSuffix(string(output), "\n"), "\n")
	if len(rows) != len(fixtures) {
		log.Fatalf("unexpected Lean row count %d", len(rows))
	}
	for i, c := range fixtures {
		if want := wire(c.unix) + "|" + wire(c.windows) + "|" + wire(c.plan9); rows[i] != want {
			log.Fatalf("fixture %d: Lean=%q want %q", i, rows[i], want)
		}
	}
	testPrivate(native.Bytes())
	fmt.Printf("private splitList passed %d literal direct native and three-profile Lean fixtures.\n", len(fixtures))
}

func testPrivate(cases []byte) {
	scratch, err := os.MkdirTemp("", "golean-split-list-control-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(scratch)
	name := filepath.Join(scratch, "cases")
	if err := os.WriteFile(name, cases, 0600); err != nil {
		log.Fatal(err)
	}
	data := "difftest/path/filepath/split-list-private/testdata"
	replacements := make(map[string]string)
	for _, test := range []string{"bridge_test.go", "split_list_test.go"} {
		source, err := filepath.Abs(filepath.Join(data, test))
		if err != nil {
			log.Fatal(err)
		}
		replacements[filepath.Join(runtime.GOROOT(), "src", "path", "filepath", "golean_control_"+test)] = source
	}
	overlay, err := json.Marshal(struct{ Replace map[string]string }{replacements})
	if err != nil {
		log.Fatal(err)
	}
	file := filepath.Join(scratch, "overlay.json")
	if err := os.WriteFile(file, overlay, 0600); err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "test", "-count=1", "-overlay", file, "-run", "^TestGoLeanPrivateSplitList$", "path/filepath")
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN=local", "GOROOT="+runtime.GOROOT(), "GODEBUG=goindex=0", "GOMODCACHE="+filepath.Join(scratch, "modcache"), "GOOS="+runtime.GOOS, "GOARCH="+runtime.GOARCH, "GOLEAN_PRIVATE_SPLIT_LIST_CASES="+name)
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("private overlay: %v\n%s", err, out)
	}
}

func wire(parts []string) string {
	encoded := make([]string, len(parts))
	for i, part := range parts {
		values := make([]string, len(part))
		for j := range len(part) {
			values[j] = strconv.Itoa(int(part[j]))
		}
		encoded[i] = strings.Join(values, ",")
	}
	return fmt.Sprintf("slice:%d:%s", len(parts), strings.Join(encoded, ";"))
}
