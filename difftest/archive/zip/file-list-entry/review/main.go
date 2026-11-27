package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Name, Other                   string
	Present, Directory, Duplicate bool
	Seconds                       int64
	Nanos, Offset                 int
	Want                          string
}

func wire(raw string) string {
	fields := make([]string, len(raw))
	for i, digit := range []byte(raw) {
		fields[i] = strconv.Itoa(int(digit))
	}
	return strings.Join(fields, ",")
}

func expected(directory, elem string, slash bool, order int, time string) string {
	return fmt.Sprintf("%s|%s|%t|%d|%s|0|2147484013|2147483648|true|true|%s|true|%s",
		wire(directory), wire(elem), slash, order, wire(elem), time, wire("d "+elem+"/"))
}

func main() {
	const zero = "-62135596800|0"
	const nonUTCConverted = "1234|567"
	cases := []fixture{
		{Name: "a//", Other: "a/z", Directory: false, Want: expected("a", "", true, -1, zero)},
		{Name: "d/\xff\x00/", Other: "d/\x80", Present: true, Duplicate: true, Seconds: 1234,
			Nanos: 567, Offset: -7 * 3600, Want: expected("d", "\xff\x00", true, 1, nonUTCConverted)},
		{Name: "x/", Other: "x", Directory: true, Want: expected(".", "x", true, 0, zero)},
		{Name: "b/a", Other: "a/z", Want: expected("b", "a", false, 1, zero)},
		{Name: "/", Other: "", Want: expected(".", "", true, 0, zero)},
	}
	var input bytes.Buffer
	truth, err := os.CreateTemp("", "review-zip-file-list-")
	must(err)
	defer os.Remove(truth.Name())
	encoder := gob.NewEncoder(truth)
	for _, c := range cases {
		present := 0
		if c.Present {
			present = 1
		}
		fmt.Fprint(&input, len(c.Name), " ", present, " ", c.Seconds, " ", c.Nanos)
		for _, text := range []string{c.Name, c.Other} {
			for _, b := range []byte(text) {
				fmt.Fprint(&input, " ", b)
			}
		}
		fmt.Fprintln(&input)
		must(encoder.Encode(c))
	}
	must(truth.Close())
	must(os.Setenv("GOLEAN_ZIP_FILE_LIST_ROWS", truth.Name()))
	must(check.TestInPackage("archive/zip", "TestGoLeanZipFileList", "difftest/archive/zip/file-list-entry/testdata/entry_test.go"))
	rows := check.LeanRows("godiff-zip-file-list-entry", len(cases), &input)
	for i, row := range rows {
		if row != cases[i].Want {
			log.Fatalf("row %d got %q, want %q", i, row, cases[i].Want)
		}
	}
	for _, m := range []struct {
		row, col int
		wrong    string
	}{
		{0, 0, wire(".")},         // strip every trailing slash
		{1, 3, "-1"},              // compare signed bytes
		{2, 3, "1"},               // use trailing slash as a sort key
		{3, 3, "-1"},              // sort by element before directory
		{0, 4, wire("a")},         // use a path-cleaning basename
		{1, 5, "9900"},            // use backing file size
		{1, 6, "2147484159"},      // use 0777 directory permissions
		{1, 7, "2147484013"},      // return permissions in Type
		{0, 8, "false"},           // consult the internal directory flag
		{1, 9, "false"},           // keep a backing object in Sys
		{0, 10, "0"},              // replace absent-file zero time with Unix epoch
		{1, 11, "0"},              // lose fractional time
		{1, 12, "false"},          // allocate a distinct info object
		{1, 13, wire("d �\x00/")}, // replace malformed name bytes
	} {
		if strings.Split(rows[m.row], "|")[m.col] == m.wrong {
			log.Fatalf("row %d column %d accepted wrong output %q", m.row, m.col, m.wrong)
		}
	}
	log.Print("five independent pinned Go and Lean rows passed and rejected fourteen wrong observations")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
