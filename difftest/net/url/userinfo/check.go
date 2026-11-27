package main

import (
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

type input struct {
	Tag                    int
	Name, Secret, Validate []byte
}
type fixture struct {
	input
	Want string
}

const alnum = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func escape(raw []byte) []byte {
	var out bytes.Buffer
	for _, b := range raw {
		if strings.IndexByte(alnum+"-_.~$&+,;=", b) >= 0 {
			out.WriteByte(b)
		} else {
			fmt.Fprintf(&out, "%%%02X", b)
		}
	}
	return out.Bytes()
}

func wire(raw []byte) string {
	fields := make([]string, len(raw))
	for i, b := range raw {
		fields[i] = strconv.Itoa(int(b))
	}
	return fmt.Sprintf("%d:%s", len(raw), strings.Join(fields, ","))
}

func expected(c input) string {
	name, secret := wire(c.Name), wire(c.Secret)
	username, password, encoded := "0:", "0::false", "0:"
	if c.Tag != 0 {
		username, password = name, secret+":"+strconv.FormatBool(c.Tag == 2)
		result := escape(c.Name)
		if c.Tag == 2 {
			result = append(append(result, ':'), escape(c.Secret)...)
		}
		encoded = wire(result)
	}
	valid := true
	for _, r := range string(c.Validate) {
		if !strings.ContainsRune(alnum+"-._:~!$&'()*+,;=%@", r) {
			valid = false
			break
		}
	}
	return strings.Join([]string{name + "/0::false", name + "/" + secret + ":true", username, password,
		encoded, strconv.FormatBool(valid)}, "|")
}

func batch(label string, cases []input, encoder *json.Encoder) {
	const shard = 1024
	for start := 0; start < len(cases); start += shard {
		current := cases[start:min(start+shard, len(cases))]
		var protocol bytes.Buffer
		rows := make([]fixture, len(current))
		for i, c := range current {
			fmt.Fprintf(&protocol, "%d %d %d", c.Tag, len(c.Name), len(c.Secret))
			for _, raw := range [][]byte{c.Name, c.Secret, c.Validate} {
				for _, b := range raw {
					fmt.Fprint(&protocol, " ", b)
				}
			}
			fmt.Fprintln(&protocol)
			rows[i] = fixture{c, expected(c)}
		}
		for i, row := range check.LeanRows("godiff-net-url-userinfo", len(current), &protocol) {
			if row != rows[i].Want {
				log.Fatalf("%s row %d %+v: got %q want %q", label, start+i, current[i], row, rows[i].Want)
			}
		}
		for _, r := range rows {
			must(encoder.Encode(r))
		}
	}
	fmt.Printf("URL userinfo %s: %d Lean rows agree.\n", label, len(cases))
}

func main() {
	var edges []input
	for i := range 1 << 16 {
		a, b := byte(i>>8), byte(i)
		edges = append(edges, input{i % 3, []byte{a}, []byte{b}, []byte{a, b}})
	}
	words := []string{"", "@", ":", "%", "%ab", "a:b", "a b", "u+v", "é", "\xff", "\x00", "\xe2\x82", "?!$&'[]"}
	for tag := range 3 {
		for _, left := range words {
			for _, right := range words {
				edges = append(edges, input{tag, []byte(left), []byte(right), []byte(left + right)})
			}
		}
	}
	f, err := os.CreateTemp("", "golean-userinfo-")
	must(err)
	defer os.Remove(f.Name())
	encoder := json.NewEncoder(f)
	batch("boundaries", edges, encoder)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]input, 20000)
		for i := range cases {
			raw := func() []byte {
				data := make([]byte, rng.Intn(40))
				if i%3 == 0 {
					const choices = "!abcXYZ09&+;:@?= %/\xff"
					for j := range data {
						data[j] = choices[rng.Intn(len(choices))]
					}
				} else {
					rng.Read(data)
				}
				return data
			}
			cases[i] = input{i % 3, raw(), raw(), raw()}
		}
		batch(fmt.Sprintf("seed=%d", seed), cases, encoder)
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_URL_USERINFO_CASES", f.Name()))
	const data = "difftest/net/url/userinfo/testdata/"
	must(check.TestInPackage("net/url", "TestGoLeanUserinfo", data+"bridge_test.go", data+"userinfo_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
