package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	file, err := os.CreateTemp("", "golean-new-replacer-")
	must(err)
	defer os.Remove(file.Name())
	buffer := bufio.NewWriter(file)
	fixture := json.NewEncoder(buffer)

	edges := [][]string{
		nil, {}, {""}, {"", ""}, {"", "", ""}, {"a", "A"},
		{"a", "A", "a", "Z"}, {"ab", "A", "a", "Z"},
		{"\xff\x00", "é", "🚀", "\xe2\x82"}, {"abc", "d", "e"},
	}
	for b := range 256 {
		raw := string([]byte{byte(b)})
		edges = append(edges, []string{raw}, []string{raw, raw + "\x00"}, []string{"", raw, raw, ""})
	}
	compare("boundaries", edges, fixture)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([][]string, 5_000)
			for i := range cases {
				args := make([]string, rng.Intn(19))
				for j := range args {
					raw := make([]byte, rng.Intn(17))
					rng.Read(raw)
					args[j] = string(raw)
				}
				cases[i] = args
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), cases, fixture)
		}
	}
	must(buffer.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_NEW_REPLACER_CASES", file.Name()))
	const data = "difftest/strings/new-replacer/testdata/"
	must(check.TestInPackage("strings", "TestGoLeanNewReplacer", data+"bridge_test.go", data+"new_replacer_test.go"))
}

func expected(args []string) string {
	switch {
	case len(args)%2 != 0:
		return "strings.NewReplacer: odd argument count"
	case len(args) == 0:
		return "true:true:nil:"
	default:
		fields := make([]string, len(args))
		for i, arg := range args {
			digits := make([]string, len(arg))
			for j := range len(arg) {
				digits[j] = strconv.Itoa(int(arg[j]))
			}
			fields[i] = strconv.Itoa(len(arg)) + ":" + strings.Join(digits, ",")
		}
		return "true:true:slice:" + strings.Join(fields, "|")
	}
}

func compare(label string, cases [][]string, fixture *json.Encoder) {
	var input bytes.Buffer
	for _, args := range cases {
		fmt.Fprint(&input, len(args))
		for _, arg := range args {
			fmt.Fprintf(&input, " %d", len(arg))
			for i := range len(arg) {
				fmt.Fprintf(&input, " %d", arg[i])
			}
		}
		fmt.Fprintln(&input)
	}
	for i, got := range check.LeanRows("godiff-strings-new-replacer", len(cases), &input) {
		args := cases[i]
		want := expected(args)
		if got != want {
			log.Fatalf("%s row %d args=%q: Lean=%s want=%s", label, i, args, got, want)
		}
		var raw [][]byte
		if args != nil {
			raw = make([][]byte, len(args))
		}
		for j, arg := range args {
			raw[j] = []byte(arg)
		}
		must(fixture.Encode(struct {
			Args [][]byte
			Want string
		}{raw, want}))
	}
	fmt.Printf("strings.NewReplacer %s: %d Lean states agree.\n", label, len(cases))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
