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
	Tag                    int
	Name, Secret, Validate []byte
	Want                   string
}

func serialize(raw string) string {
	fields := make([]string, len(raw))
	for i := range fields {
		fields[i] = strconv.Itoa(int(raw[i]))
	}
	return fmt.Sprintf("%d:%s", len(raw), strings.Join(fields, ","))
}
func row(tag int, name, secret, validate, escaped string, valid bool) fixture {
	n, p := serialize(name), serialize(secret)
	who, pass := "0:", "0::false"
	if tag != 0 {
		who, pass = n, p+":"+strconv.FormatBool(tag == 2)
	}
	return fixture{tag, []byte(name), []byte(secret), []byte(validate), strings.Join([]string{
		n + "/0::false", n + "/" + p + ":true", who, pass, serialize(escaped), strconv.FormatBool(valid),
	}, "|")}
}

func main() {
	cases := []fixture{
		row(0, "hidden", "password", "@%z+:", "", true),
		row(2, "root", "", "", "root:", true),
		row(1, "root", "retained", "a/b", "root", false),
		row(2, "\xff:@+ ", "\x00!/?%", "%zz@", "%FF%3A%40+%20:%00%21%2F%3F%25", true),
		row(1, "a+ b", "\xff", "\xff", "a+%20b", false),
		row(2, "", "ä", "ä", ":%C3%A4", false),
		row(2, "", "\x00", "\x00", ":%00", false),
	}
	var wire bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&wire, "%d %d %d", c.Tag, len(c.Name), len(c.Secret))
		for _, raw := range [][]byte{c.Name, c.Secret, c.Validate} {
			for _, v := range raw {
				fmt.Fprintf(&wire, " %d", v)
			}
		}
		fmt.Fprintln(&wire)
	}
	mutants := []struct {
		column      int
		alternative string
	}{
		{5, "false"},
		{4, serialize("root")},
		{3, "0::false"},
		{4, serialize("%FF%3A%40%2B+:%00%21%2F%3F%25")},
		{5, "true"},
		{5, "true"},
		{5, "true"},
	}
	rejects := 0
	for i, got := range check.LeanRows("godiff-net-url-userinfo", len(cases), &wire) {
		if got != cases[i].Want {
			log.Fatalf("row %d: Lean %q want %q", i, got, cases[i].Want)
		}
		if m := mutants[i]; strings.Split(got, "|")[m.column] != m.alternative {
			rejects++
		}
	}
	if rejects != len(mutants) {
		log.Fatalf("rejected %d of %d mutants", rejects, len(mutants))
	}
	f, err := os.CreateTemp("", "golean-url-userinfo-review-")
	must(err)
	defer os.Remove(f.Name())
	rows := json.NewEncoder(f)
	for _, c := range cases {
		must(rows.Encode(c))
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_URL_USERINFO_CASES", f.Name()))
	const dir = "difftest/net/url/userinfo/testdata/"
	must(check.TestInPackage("net/url", "TestGoLeanUserinfo", dir+"bridge_test.go", dir+"userinfo_test.go"))
	fmt.Printf("URL userinfo: %d independent native and Lean cases, %d altered results rejected\n", len(cases), rejects)
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
