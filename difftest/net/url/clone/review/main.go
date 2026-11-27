package main

import (
	"bytes"
	"log"
	"net/url"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	var missing *url.URL
	if missing.Clone() != nil {
		log.Fatal("nil receiver yielded a URL")
	}

	noUser := &url.URL{}
	if clone := noUser.Clone(); clone == nil || clone == noUser || clone.User != nil {
		log.Fatal("zero URL did not become a new object with nil userinfo")
	}

	noPassword := &url.URL{User: new(url.Userinfo)}
	if clone := noPassword.Clone(); clone == nil || clone == noPassword || clone.User == nil ||
		clone.User == noPassword.User || clone.User.Username() != "" {
		log.Fatal("present empty userinfo was not copied")
	} else if _, set := clone.User.Password(); set {
		log.Fatal("absent password became present")
	}

	original := &url.URL{Host: "\xff", RawPath: "%2f", ForceQuery: true, User: url.UserPassword("", "")}
	clone := original.Clone()
	if clone == nil || clone == original || clone.User == nil || clone.User == original.User {
		log.Fatal("original and clone must have different URL and userinfo pointers")
	}
	if password, set := clone.User.Password(); password != "" || !set {
		log.Fatal("present empty password was lost")
	}
	clone.Scheme, clone.Host, clone.RawPath, clone.ForceQuery = "other", "other", "other", false
	clone.User = url.User("other")
	if original.Scheme != "" || original.Host != "\xff" || original.RawPath != "%2f" ||
		!original.ForceQuery || original.User.Username() != "" {
		log.Fatal("mutating the clone altered the original")
	}

	input := bytes.NewBufferString("0 0 0 0 0 0 0 0 0 0 0 0 0 0\n" +
		"1 0 0 0 0 0 0 0 0 0 0 0 0 0\n" +
		"1 1 0 0 0 0 0 0 0 0 0 0 0 0\n" +
		"1 2 1 0 0 0 1 0 0 0 3 0 0 0 255 37 50 102\n")
	wants := []string{
		"nil",
		"||||||||||false|false|0|true|true",
		"||||||||||false|false|1|true|true",
		"||255||||37,50,102||||true|false|2|true|true",
	}
	actual := check.LeanRows("godiff-url-clone", len(wants), input)
	for i, row := range actual {
		if row != wants[i] {
			log.Fatalf("row %d = %q, want %q", i, row, wants[i])
		}
	}
	with := func(row, field int, replacement string) string {
		values := strings.Split(actual[row], "|")
		values[field] = replacement
		return strings.Join(values, "|")
	}
	mutants := []struct {
		row   int
		value string
	}{
		{0, actual[1]},            // allocate for a nil receiver
		{1, with(1, 13, "false")}, // return the original URL
		{2, with(2, 14, "false")}, // share the original userinfo
		{2, with(2, 12, "0")},     // erase present but empty userinfo
		{3, with(3, 12, "1")},     // erase presence of an empty password
		{3, with(3, 6, "")},       // omit the raw path
	}
	for _, m := range mutants {
		if m.value == actual[m.row] {
			log.Fatalf("row %d did not reject mutant %q", m.row, m.value)
		}
	}
	log.Print("four independent native and Lean fixtures passed and rejected six wrong observations")
}
