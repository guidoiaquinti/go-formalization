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

type fixture struct {
	Present              bool
	Mode                 byte
	Decoded, Raw, Update []byte
	Want                 string
}

func encoded(s string) string {
	fields := make([]string, len(s))
	for i := range fields {
		fields[i] = strconv.Itoa(int(s[i]))
	}
	return strings.Join(fields, ",")
}

func ok(s string) string             { return "ok:" + encoded(s) }
func set(decoded, raw string) string { return "ok/" + encoded(decoded) + "/" + encoded(raw) }
func failed(payload, decoded, raw string) string {
	return "escape:" + encoded(payload) + "/" + encoded(decoded) + "/" + encoded(raw)
}

func fixtureOf(present bool, mode byte, decoded, raw, update string, valid bool, outputs ...string) fixture {
	return fixture{present, mode, []byte(decoded), []byte(raw), []byte(update),
		strconv.FormatBool(valid) + "|" + strings.Join(outputs, "|")}
}

func main() {
	cases := []fixture{
		fixtureOf(true, 1, "*", "%2a", "%2a", true, ok("%2a"), ok("%2a"), set("*", "%2a"), set("*", "%2a")),
		fixtureOf(true, 1, "*", "%zz", "%", true, ok("*"), ok("*"), failed("%", "*", "%zz"), failed("%", "*", "%zz")),
		fixtureOf(true, 1, "/a?b", "/a?b", "/a?b", false, ok("/a%3Fb"), ok("/a?b"), set("/a?b", "/a?b"), set("/a?b", "")),
		fixtureOf(true, 1, "/", "%20", "%2f", true, ok("/"), ok("/"), set("/", "%2f"), set("/", "%2f")),
		fixtureOf(true, 1, "[", "[", "%FF%00", true, ok("["), ok("["), set("\xff\x00", ""), set("\xff\x00", "")),
		fixtureOf(true, 1, " ", " ", "+", false, ok("%20"), ok("%20"), set("+", ""), set("+", "")),
		fixtureOf(false, 1, "", "%", "+%1", true, "panic", "panic", "escape:"+encoded("%1")+"/nil", "escape:"+encoded("%1")+"/nil"),
		fixtureOf(false, 1, "", "", "%20", true, "panic", "panic", "panic", "panic"),
	}
	var rows bytes.Buffer
	for _, c := range cases {
		tag := 0
		if c.Present {
			tag = 1
		}
		fmt.Fprint(&rows, tag, " ", c.Mode, " ", len(c.Decoded), " ", len(c.Raw))
		for _, block := range [][]byte{c.Decoded, c.Raw, c.Update} {
			for _, b := range block {
				fmt.Fprint(&rows, " ", b)
			}
		}
		fmt.Fprintln(&rows)
	}
	if len(os.Args) == 2 && os.Args[1] == "rows" {
		_, _ = os.Stdout.Write(rows.Bytes())
		return
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-net-url-path-fragment")
	cmd.Stdin = &rows
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	results := strings.Split(strings.TrimSuffix(string(output), "\n"), "\n")
	if len(results) != len(cases) {
		log.Fatalf("got %d Lean fixtures, want %d", len(results), len(cases))
	}
	for i, got := range results {
		if got != cases[i].Want {
			log.Fatalf("fixture %d Lean=%q want %q", i, got, cases[i].Want)
		}
	}
	checkNative(cases)
	fmt.Printf("net/url path and fragment reviewer control: %d pinned private native and Lean fixtures passed.\n", len(cases))
}

func checkNative(cases []fixture) {
	dir, err := os.MkdirTemp("", "golean-url-path-control-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	caseFile := filepath.Join(dir, "cases.json")
	data, err := encodeCases(cases)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(caseFile, data, 0600); err != nil {
		log.Fatal(err)
	}
	replacements := make(map[string]string)
	for _, name := range []string{"bridge_test.go", "path_fragment_test.go"} {
		path, err := filepath.Abs(filepath.Join("difftest/net/url/path-fragment/testdata", name))
		if err != nil {
			log.Fatal(err)
		}
		replacements[filepath.Join(runtime.GOROOT(), "src", "net", "url", "golean_control_"+name)] = path
	}
	overlay, err := json.Marshal(struct{ Replace map[string]string }{replacements})
	if err != nil {
		log.Fatal(err)
	}
	file := filepath.Join(dir, "overlay.json")
	if err := os.WriteFile(file, overlay, 0600); err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "test", "-count=1", "-overlay", file, "-run", "^TestGoLeanPathFragment$", "-v", "net/url")
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN=local", "GOROOT="+runtime.GOROOT(), "GODEBUG=goindex=0",
		"GOMODCACHE="+filepath.Join(dir, "modcache"), "GOOS="+runtime.GOOS, "GOARCH="+runtime.GOARCH,
		"GOLEAN_URL_PATH_FRAGMENT_CASES="+caseFile)
	out, err := cmd.CombinedOutput()
	if err != nil || !bytes.Contains(out, []byte("--- PASS: TestGoLeanPathFragment")) {
		log.Fatalf("private native control: %v\n%s", err, out)
	}
}

func encodeCases(cases []fixture) ([]byte, error) {
	var data bytes.Buffer
	encoder := json.NewEncoder(&data)
	for _, c := range cases {
		if err := encoder.Encode(c); err != nil {
			return nil, err
		}
	}
	return data.Bytes(), nil
}
