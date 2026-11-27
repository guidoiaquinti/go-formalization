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
	Present                             bool
	Type, Nested                        int
	Value, Struct, Field, Source, Inner []byte
	Want                                string
}

func wire(s string) string {
	fields := make([]string, len(s))
	for i, b := range []byte(s) {
		fields[i] = strconv.Itoa(int(b))
	}
	return strings.Join(fields, ",")
}

func observation(value, typ, invalid, marshal string, id int, unmarshal string) string {
	return strings.Join([]string{wire(value), wire("json: unsupported value: " + value),
		wire(typ), wire(invalid), marshal, strconv.Itoa(id), unmarshal, strconv.Itoa(id)}, "|")
}

func main() {
	cases := []fixture{
		{true, 2, 1, []byte("number\xff"), []byte("Item"), []byte("data.01"), nil, []byte("detail\x00"),
			observation("number\xff", "json: unsupported type: *int", "json: Unmarshal(nil *int)",
				wire("json: error calling MarshalJSON for type *int: detail\x00"), 1,
				wire("json: cannot unmarshal number\xff into Item.data.01 of type *int: detail\x00"))},
		{true, 1, 1, []byte("number"), []byte("Item"), []byte("data.1."), []byte("MarshalText"), []byte("cause"),
			observation("number", "json: unsupported type: int", "json: Unmarshal(non-pointer int)",
				wire("json: error calling MarshalText for type int: cause"), 1,
				wire("json: cannot unmarshal number into Go struct field Item.data.1. of type int: cause"))},
		{true, 3, 0, []byte("string"), nil, nil, nil, nil,
			observation("string", "json: unsupported type: []string", "json: Unmarshal(non-pointer []string)",
				"nilpanic", 0, wire("json: cannot unmarshal string into Go value of type []string"))},
		{true, 2, 2, []byte("string"), nil, nil, nil, nil,
			observation("string", "json: unsupported type: *int", "json: Unmarshal(nil *int)",
				"nestedpanic", 2, "nestedpanic")},
		{true, 0, 2, nil, nil, nil, nil, nil,
			strings.Join([]string{"", wire("json: unsupported value: "), "nilpanic",
				wire("json: Unmarshal(nil)"), "nilpanic", "2", "nilpanic", "2"}, "|")},
		{false, 0, 0, nil, nil, nil, nil, nil, strings.TrimSuffix(strings.Repeat("nilpanic|", 8), "|")},
	}
	var input bytes.Buffer
	truth, err := os.CreateTemp("", "review-json-compat-errors-")
	must(err)
	defer os.Remove(truth.Name())
	enc := json.NewEncoder(truth)
	for _, c := range cases {
		present := 0
		if c.Present {
			present = 1
		}
		typeInfo := c.Type
		typeText := []string{"", "int", "*int", "[]string"}[c.Type]
		if typeInfo == 3 {
			typeInfo = 1
		}
		fmt.Fprint(&input, present, " ", typeInfo, " ", c.Nested)
		for _, raw := range [][]byte{[]byte(typeText), c.Inner, c.Value, c.Struct, c.Field, c.Source} {
			fmt.Fprint(&input, " ", len(raw))
			for _, b := range raw {
				fmt.Fprint(&input, " ", b)
			}
		}
		fmt.Fprintln(&input)
		must(enc.Encode(c))
	}
	must(truth.Close())
	must(os.Setenv("GOLEAN_JSON_COMPAT_ERRORS", truth.Name()))
	must(check.TestInPackage("encoding/json", "TestGoLeanCompatErrors", "difftest/encoding/json/compat-errors/testdata/errors_test.go"))
	rows := check.LeanRows("godiff-json-compat-errors", len(cases), &input)
	for i, row := range rows {
		if row != cases[i].Want {
			log.Fatalf("row %d got %q, want %q", i, row, cases[i].Want)
		}
	}
	for _, m := range []struct {
		row, col int
		value    string
	}{
		{0, 0, wire("number�")},
		{0, 1, wire("json: unsupported value: 117")},
		{4, 2, wire("json: unsupported type: <nil>")},
		{0, 3, wire("json: Unmarshal(non-pointer *int)")},
		{0, 4, wire("json: error calling  for type *int: detail\x00")},
		{4, 4, "nestedpanic"}, // invoke nested callback before type rendering
		{3, 5, "0"},
		{0, 6, wire("json: cannot unmarshal number\xff into Go struct field Item.data.01 of type *int: detail\x00")},
		{1, 6, wire("json: cannot unmarshal number into Item.data.1. of type int: cause")},
		{4, 6, "nestedpanic"},
		{3, 7, "0"},
	} {
		if strings.Split(rows[m.row], "|")[m.col] == m.value {
			log.Fatalf("row %d column %d accepted mutation %q", m.row, m.col, m.value)
		}
	}
	log.Print("six independent pinned Go and Lean rows passed and rejected eleven incorrect outcomes")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
