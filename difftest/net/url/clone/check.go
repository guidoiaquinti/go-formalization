package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	present, force, omit bool
	user                 int
	fields               [10]string
}

func asInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func render(c fixture, freshURL, freshUser bool) string {
	if !c.present {
		return "nil"
	}
	fields := make([]string, 0, 15)
	for _, s := range c.fields {
		b := []byte(s)
		values := make([]string, len(b))
		for i, v := range b {
			values[i] = strconv.Itoa(int(v))
		}
		fields = append(fields, strings.Join(values, ","))
	}
	return strings.Join(append(fields, strconv.FormatBool(c.force), strconv.FormatBool(c.omit),
		strconv.Itoa(c.user), strconv.FormatBool(freshURL), strconv.FormatBool(freshUser)), "|")
}

func runGo(c fixture) string {
	var original *url.URL
	if c.present {
		f := c.fields
		original = &url.URL{Scheme: f[0], Opaque: f[1], Host: f[2], Path: f[3], Fragment: f[4],
			RawQuery: f[5], RawPath: f[6], RawFragment: f[7], ForceQuery: c.force, OmitHost: c.omit}
		switch c.user {
		case 1:
			original.User = url.User(f[8])
		case 2:
			original.User = url.UserPassword(f[8], f[9])
		}
	}
	clone := original.Clone()
	if clone == nil {
		return "nil"
	}
	actual := fixture{present: true, force: clone.ForceQuery, omit: clone.OmitHost,
		fields: [10]string{clone.Scheme, clone.Opaque, clone.Host, clone.Path, clone.Fragment,
			clone.RawQuery, clone.RawPath, clone.RawFragment}}
	if clone.User != nil {
		actual.user = 1
		actual.fields[8] = clone.User.Username()
		if password, set := clone.User.Password(); set {
			actual.user = 2
			actual.fields[9] = password
		}
	}
	freshURL, freshUser := clone != original, clone.User == nil || clone.User != original.User
	return render(actual, freshURL, freshUser)
}

func compare(label string, cases []fixture) {
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&input, asInt(c.present), " ", c.user, " ", asInt(c.force), " ", asInt(c.omit))
		for _, s := range c.fields {
			fmt.Fprint(&input, " ", len(s))
		}
		for _, s := range c.fields {
			for _, b := range []byte(s) {
				fmt.Fprint(&input, " ", b)
			}
		}
		fmt.Fprintln(&input)
	}
	for i, lean := range check.LeanRows("godiff-url-clone", len(cases), &input) {
		c := cases[i]
		want := render(c, true, true)
		if goResult := runGo(c); lean != want || goResult != want {
			log.Fatalf("%s row %d %+v: Go=%q Lean=%q want=%q", label, i, c, goResult, lean, want)
		}
	}
	fmt.Printf("%s: %d pinned Go and Lean clone results agree\n", label, len(cases))
}

func main() {
	boundaries := []fixture{{}}
	for user := range 3 {
		for flags := range 4 {
			for _, name := range []string{"", "user", "\x00\xff", "🙂"} {
				c := fixture{present: true, user: user, force: flags&1 != 0, omit: flags&2 != 0,
					fields: [10]string{"scheme", "opaque", "host", "path", "fragment", "rawquery", "rawpath", "rawfragment"}}
				if user != 0 {
					c.fields[8] = name
				}
				if user == 2 && flags&1 != 0 {
					c.fields[9] = name
				}
				boundaries = append(boundaries, c)
			}
		}
	}
	compare("boundaries", boundaries)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			c := fixture{present: rng.Intn(9) != 0, user: rng.Intn(3), force: rng.Intn(2) != 0, omit: rng.Intn(2) != 0}
			limit := 8
			if c.user != 0 {
				limit++
			}
			if c.user == 2 {
				limit++
			}
			for j := range limit {
				b := make([]byte, rng.Intn(10))
				rng.Read(b)
				c.fields[j] = string(b)
			}
			cases[i] = c
		}
		compare(fmt.Sprintf("seed=%d", seed), cases)
	}
}
