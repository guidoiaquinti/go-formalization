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
	Config        string
	Raw           []byte
	Params, Start int
	Present       bool
	Want          string
}

func mk(config, raw string, params, start int, present bool, want ...string) fixture {
	return fixture{config, []byte(raw), params, start, present, strings.Join(want, "|")}
}

func main() {
	cases := []fixture{
		mk("", "", 10000, 0, true, "true:false", "ok/nil", "ok/map:", "map:"),
		mk("", "", 10001, 1, false, "false:false", "ok/map:", "ok/map:", "panic"),
		mk("2", "a=x&", 2, 2, true, "false:true", "limit/map:1:97=list:3:111,108,100/1:122=nil", "limit/map:", "map:"),
		mk("0", "a=1&&a=2", 10001, 2, true, "true:true",
			"ok/map:1:97=list:3:111,108,100;1:49;1:50/1:122=nil",
			"ok/map:1:97=list:1:49;1:50", "map:1:97=list:1:49;1:50"),
		mk("", "%&a=1&ignored;=&a=2&b=%zz&c=x+y", 1, 1, true, "true:false",
			"semicolon/map:1:97=list:1:49;1:50/1:99=list:3:120,32,121",
			"semicolon/map:1:97=list:1:49;1:50/1:99=list:3:120,32,121",
			"map:1:97=list:1:49;1:50/1:99=list:3:120,32,121"),
		mk("", "%X=%Y&%zz=v&a%3Bb=%FF", 1, 1, true, "true:false",
			"escape:2:37,88/map:3:97,59,98=list:1:255",
			"escape:2:37,88/map:3:97,59,98=list:1:255", "map:3:97,59,98=list:1:255"),
		mk("9223372036854775808", "x", 10001, 0, true, "false:false", "panic",
			"ok/map:1:120=list:0:", "map:1:120=list:0:"),
		mk("-2", "a", -3, 0, false, "true:false", "limit/nil", "limit/map:", "panic"),
	}
	var rows bytes.Buffer
	for _, c := range cases {
		present := 0
		if c.Present {
			present = 1
		}
		fmt.Fprint(&rows, strconv.IntSize, " ", c.Params, " ", c.Start, " ", present, " ", len(c.Config))
		for _, block := range [][]byte{[]byte(c.Config), c.Raw} {
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
	cmd := exec.Command("lean/.lake/build/bin/godiff-net-url-query")
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
	fmt.Printf("net/url query reviewer control: %d pinned private native and Lean fixtures passed.\n", len(cases))
}

func checkNative(cases []fixture) {
	dir, err := os.MkdirTemp("", "golean-url-query-control-")
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
	for _, name := range []string{"bridge_test.go", "query_test.go"} {
		path, err := filepath.Abs(filepath.Join("difftest/net/url/query/testdata", name))
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
	cmd := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "test", "-count=1", "-overlay", file, "-run", "^TestGoLeanURLQuery$", "-v", "net/url")
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN=local", "GOROOT="+runtime.GOROOT(), "GODEBUG=goindex=0",
		"GOMODCACHE="+filepath.Join(dir, "modcache"), "GOOS="+runtime.GOOS, "GOARCH="+runtime.GOARCH,
		"GOLEAN_URL_QUERY_CASES="+caseFile)
	out, err := cmd.CombinedOutput()
	if err != nil || !bytes.Contains(out, []byte("--- PASS: TestGoLeanURLQuery")) {
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
