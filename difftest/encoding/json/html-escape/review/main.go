package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Present      bool
	Backing, Src []byte
	Length       int
	Want         string
}

func wire(data []byte) string {
	if data == nil {
		return "nil"
	}
	fields := make([]string, len(data))
	for i, v := range data {
		fields[i] = strconv.Itoa(int(v))
	}
	return "slice:" + strings.Join(fields, ",")
}
func row(present bool, old string, length int, input, output, retained string, reused bool) fixture {
	pub := "panic"
	if present {
		pub = wire([]byte(output))
	}
	return fixture{present, []byte(old), []byte(input), length, fmt.Sprintf("%s/%s/%t|%s", wire([]byte(output)), wire([]byte(retained)), reused, pub)}
}

func main() {
	cases := []fixture{
		{true, nil, nil, 0, "nil/nil/false|slice:"},
		row(false, "", 0, "<", "\\u003c", "", false),
		row(true, "Pxxx", 1, "ab<c", "Pab\\u003cc", "Pabx", false),
		row(true, "Pxxxxxxx", 1, "ab<c", "Pab\\u003cc", "Pabxxxxx", false),
		row(true, "Pxxxxxxxxxx", 1, "ab<c", "Pab\\u003cc", "Pab\\u003ccx", true),
		row(true, "", 0, "\xe2\x80\xa8\xe2\x80\xa9\xe2\x80\xaa", "\\u2028\\u2029\xe2\x80\xaa", "", false),
		row(true, "", 0, "\xff\xe2\x80<>&", "\xff\xe2\x80\\u003c\\u003e\\u0026", "", false),
		row(true, "", 0, "\\u003c \" &", "\\u003c \" \\u0026", "", false),
	}
	mutations := map[int]string{
		0: "slice:/nil/false|slice:",
		1: wire([]byte("\\u003c")) + "/slice:/false|slice:",
		2: wire([]byte("Pab\\u003cc")) + "/" + wire([]byte("Pxxx")) + "/false|" + wire([]byte("Pab\\u003cc")),
		3: wire([]byte("Pab\\u003cc")) + "/" + wire([]byte("Pab\\u003")) + "/false|" + wire([]byte("Pab\\u003cc")),
		5: wire([]byte("\\u2028\\u2029\\u202a")) + "/slice:/false|" + wire([]byte("\\u2028\\u2029\\u202a")),
		6: wire([]byte("\xef\xbf\xbd\\u003c\\u003e\\u0026")) + "/slice:/false|" + wire([]byte("\xef\xbf\xbd\\u003c\\u003e\\u0026")),
	}
	var input bytes.Buffer
	for _, c := range cases {
		p, backed := 0, 0
		if c.Present {
			p = 1
		}
		if c.Backing != nil {
			backed = 1
		}
		fmt.Fprintf(&input, "%d %d %d %d", p, backed, c.Length, len(c.Backing))
		for _, part := range [][]byte{c.Backing, c.Src} {
			for _, v := range part {
				fmt.Fprintf(&input, " %d", v)
			}
		}
		fmt.Fprintln(&input)
	}
	rejected := 0
	for i, got := range check.LeanRows("godiff-encoding-json-html-escape", len(cases), &input) {
		if got != cases[i].Want {
			log.Fatalf("row %d: Lean %q want %q", i, got, cases[i].Want)
		}
		if alt, ok := mutations[i]; ok && alt != got {
			rejected++
		}
	}
	if rejected != len(mutations) {
		log.Fatalf("rejected %d altered outcomes", rejected)
	}
	f, err := os.CreateTemp("", "golean-json-html-review-")
	must(err)
	defer os.Remove(f.Name())
	rows := json.NewEncoder(f)
	for _, c := range cases {
		must(rows.Encode(c))
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_JSON_HTML_CASES", f.Name()))
	const dir = "difftest/encoding/json/html-escape/testdata/"
	must(check.TestInPackage("encoding/json", "TestGoLeanHTMLEscape", dir+"bridge_test.go", dir+"html_test.go"))
	fmt.Printf("JSON HTML escape: %d paired fixed cases, %d altered outcomes rejected\n", len(cases), rejected)
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
