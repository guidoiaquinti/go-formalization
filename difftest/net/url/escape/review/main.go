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
	Raw  []byte
	Mode byte
	Want string
}

func bytesWire(value string) string {
	encoded := make([]string, len(value))
	for i := range len(value) {
		encoded[i] = strconv.Itoa(int(value[i]))
	}
	return strings.Join(encoded, ",")
}

func fixed(input string, mode byte, class, escaped, query, path string) fixture {
	return fixture{[]byte(input), mode, strings.Join([]string{class,
		bytesWire(escaped), bytesWire(query), bytesWire(path)}, "|")}
}

func main() {
	cases := []fixture{
		fixed("", 0, "", "", "", ""),
		fixed("A", 0, "true:true:10", "%41", "A", "A"),
		fixed("Ff", 128, "false:true:15,false:true:15", "Ff", "Ff", "Ff"),
		fixed("g", 128, "true:false:16", "%67", "g", "g"),
		fixed(" ", 32, "true:false:0", "+", "+", "%20"),
		fixed(" ", 33, "true:false:0", "%20", "+", "%20"),
		fixed("+", 2, "false:false:11", "+", "%2B", "+"),
		fixed("/", 1, "false:false:15", "/", "%2F", "%2F"),
		fixed("\xff", 2, "true:false:42", "%FF", "%FF", "%FF"),
	}
	mutations := map[int]fixture{
		3: fixed("g", 128, "false:true:16", "g", "g", "g"),
		5: fixed(" ", 33, "true:false:0", "+", "+", "%20"),
		6: fixed("+", 2, "true:false:11", "%2B", "%2B", "%2B"),
		8: fixed("\xff", 2, "true:false:42", "%ff", "%ff", "%ff"),
	}
	var input bytes.Buffer
	for _, f := range cases {
		fmt.Fprint(&input, f.Mode)
		for _, b := range f.Raw {
			fmt.Fprintf(&input, " %d", b)
		}
		fmt.Fprintln(&input)
	}
	rows := check.LeanRows("godiff-net-url-escape", len(cases), &input)
	file, err := os.CreateTemp("", "golean-url-escape-review")
	must(err)
	defer os.Remove(file.Name())
	enc := json.NewEncoder(file)
	rejected := 0
	for i, f := range cases {
		if rows[i] != f.Want {
			log.Fatalf("fixture %d: Lean %q want %q", i, rows[i], f.Want)
		}
		must(enc.Encode(f))
		if broken, ok := mutations[i]; ok && broken.Want != rows[i] {
			rejected++
		}
	}
	must(file.Close())
	must(os.Setenv("GOLEAN_URL_ESCAPE_CASES", file.Name()))
	const overlay = "difftest/net/url/escape/testdata/"
	must(check.TestInPackage("net/url", "TestGoLeanEscape", overlay+"bridge_test.go", overlay+"escape_test.go"))
	if rejected != 4 {
		log.Fatalf("expected four rejected mutation witnesses, got %d", rejected)
	}
	fmt.Println("Nine independent rows matched pinned native Go and Lean; four mutation witnesses rejected.")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
