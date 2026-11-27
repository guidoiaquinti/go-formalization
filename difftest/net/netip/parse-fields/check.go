package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	value       string
	start, stop int
	fields      [4]byte
}

func main() {
	f, err := os.CreateTemp("", "golean-netip-fields-")
	must(err)
	defer os.Remove(f.Name())
	out := bufio.NewWriter(f)
	var fixed []fixture
	for _, text := range []string{"", "0", "1.2.3.4", "255.255.255.255", "256.1.2.3", "1.02.3.4", "0.0.0.0", "1.2.3", "1.2.3.", "1.2.3.4.",
		".1.2.3", "1..2.3", "1.2.3.4.5", "1.2.3.4 extra", "1.2.3.\xff", "1.256x.3.4", "1.0000000", "1.2.00.4",
		"host:8", "[]:8", "[::1%zone]:8", "[:colon]:port", "[host:8", "host:", ":8", ":", "\x00:abc\xff", "x::y", "[x]more:0"} {
		for lo := 0; lo <= len(text); lo++ {
			for hi := lo; hi <= len(text); hi++ {
				fixed = append(fixed, fixture{text, lo, hi, [4]byte{91, 92, 93, 94}})
			}
		}
	}
	compare("fixed", fixed, out)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			rows := make([]fixture, 5000)
			for i := range rows {
				var text string
				switch i % 4 {
				case 0:
					text = fmt.Sprintf("%d.%d.%d.%d", rng.Intn(300), rng.Intn(300), rng.Intn(300), rng.Intn(300))
				case 1:
					text = fmt.Sprintf("[host:%x]:%d", rng.Uint32(), rng.Intn(100000))
				case 2, 3:
					raw := make([]byte, rng.Intn(50))
					rng.Read(raw)
					if i%4 == 3 {
						for j := range raw {
							const alphabet = "0500123456789.....:::[] "
							raw[j] = alphabet[int(raw[j])%len(alphabet)]
						}
					}
					text = string(raw)
				}
				lo, hi := 0, len(text)
				if i%3 == 0 {
					lo = rng.Intn(hi + 1)
					hi = lo + rng.Intn(hi-lo+1)
				}
				var fields [4]byte
				rng.Read(fields[:])
				rows[i] = fixture{text, lo, hi, fields}
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), rows, out)
		}
	}
	must(out.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_NETIP_PARSE_FIELDS", f.Name()))
	must(check.TestInPackage("net/netip", "TestGoLeanParseFields", "difftest/net/netip/parse-fields/testdata/parse_fields_test.go"))
}

func compare(label string, rows []fixture, out *bufio.Writer) {
	var input bytes.Buffer
	for _, r := range rows {
		fmt.Fprint(&input, r.start, " ", r.stop)
		for _, b := range r.fields {
			fmt.Fprint(&input, " ", b)
		}
		for _, b := range []byte(r.value) {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
	}
	for i, result := range check.LeanRows("godiff-netip-parse-fields", len(rows), &input) {
		r := rows[i]
		_, err := fmt.Fprintf(out, "%d %d %d %d %d %d %q\t%s\n",
			r.start, r.stop, r.fields[0], r.fields[1], r.fields[2], r.fields[3], r.value, result)
		must(err)
	}
	fmt.Printf("netip parse fields %s: %d Lean rows\n", label, len(rows))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
