package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	cases := []struct {
		args []string
		want string
	}{
		{nil, "true:true:nil:"},
		{[]string{}, "true:true:nil:"},
		{[]string{"a"}, "strings.NewReplacer: odd argument count"},
		{[]string{"a", "b", "c"}, "strings.NewReplacer: odd argument count"},
		{[]string{"", ""}, "true:true:slice:0:|0:"},
		{[]string{"ab", "z"}, "true:true:slice:2:97,98|1:122"},
		{[]string{"a", "b", "a", "c"}, "true:true:slice:1:97|1:98|1:97|1:99"},
		{[]string{"\xff\x00", "\x80"}, "true:true:slice:2:255,0|1:128"},
	}
	f, err := os.CreateTemp("", "golean-review-newreplacer-")
	must(err)
	defer os.Remove(f.Name())
	enc := json.NewEncoder(f)
	var model bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&model, len(c.args))
		var args [][]byte
		if c.args != nil {
			args = make([][]byte, len(c.args))
		}
		for i, arg := range c.args {
			fmt.Fprintf(&model, " %d", len(arg))
			for j := range len(arg) {
				fmt.Fprintf(&model, " %d", arg[j])
			}
			args[i] = []byte(arg)
		}
		fmt.Fprintln(&model)
		must(enc.Encode(struct {
			Args [][]byte
			Want string
		}{args, c.want}))
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_NEW_REPLACER_CASES", f.Name()))
	const data = "difftest/strings/new-replacer/testdata/"
	must(check.TestInPackage("strings", "TestGoLeanNewReplacer", data+"bridge_test.go", data+"new_replacer_test.go"))
	for i, got := range check.LeanRows("godiff-strings-new-replacer", len(cases), &model) {
		if got != cases[i].want {
			log.Fatalf("Lean row %d: got %q want %q", i, got, cases[i].want)
		}
	}
	if cases[5].want == "true:true:slice:3:97,98,0|2:122,0" ||
		cases[5].want == "false:false:slice:2:97,98|1:122" ||
		cases[1].want == "true:true:slice:" {
		log.Fatal("constructor mutation survived")
	}
	fmt.Printf("NewReplacer passed %d independent Go and Lean observations\n", len(cases))
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
