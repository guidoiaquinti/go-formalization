package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Family                 int
	Hi, Lo                 uint64
	Port                   uint16
	Bits, Length           int
	Zone, Backing, Payload []byte
	Want                   string
}

func csv(raw []byte) string {
	fields := make([]string, len(raw))
	for i, v := range raw {
		fields[i] = strconv.Itoa(int(v))
	}
	return strings.Join(fields, ",")
}

func appended(reused, nilSlice bool, output, original []byte) string {
	return fmt.Sprintf("%t:%t:%s:%s", reused, nilSlice, csv(output), csv(original))
}

func ip(family int, hi, lo uint64, zone []byte) string {
	return fmt.Sprintf("%d:%d:%d:%s", family, hi, lo, csv(zone))
}

func main() {
	const hi = 0x0102030405060708
	const lo = 0x1112131415161718
	const v4 = 0xffff01020304
	base := []byte{1, 2, 3, 4, 5, 6, 7, 8, 17, 18, 19, 20, 21, 22, 23, 24}
	a := append(slices.Clone(base), 255, 0)
	p := append(slices.Clone(a), 52, 18)
	q := append(slices.Clone(base), 65)
	staged := append([]byte{91}, base[:8]...)
	old := append([]byte{91}, bytes.Repeat([]byte{73}, 8)...)
	cases := []fixture{
		{6, hi, lo, 0x1234, 65, 1, []byte{255, 0}, old, []byte{1, 2, 3, 4, 33}, strings.Join([]string{
			"18", csv(a), appended(false, false, append([]byte{91}, a...), staged), "false:" + ip(6, hi, lo, []byte{255, 0}),
			csv(p), appended(false, false, append([]byte{91}, p...), staged), "false:" + ip(6, hi, lo, []byte{255, 0}) + ":4660",
			csv(q), appended(false, false, append([]byte{91}, q...), staged), "true:" + ip(4, 0, v4, nil) + ":-1",
		}, "|")},
		{4, 0, v4, 0x1234, 24, 1, nil, []byte{91, 73, 73, 73, 73}, []byte{1, 2, 3, 4, 52, 18}, strings.Join([]string{
			"4", "1,2,3,4", "true:false:91,1,2,3,4:91,1,2,3,4", "false:" + ip(4, 0, v4, nil),
			"1,2,3,4,52,18", "false:false:91,1,2,3,4,52,18:91,1,2,3,4", "true:" + ip(4, 0, v4, nil) + ":4660",
			"1,2,3,4,24", "false:false:91,1,2,3,4,24:91,1,2,3,4", "false:" + ip(4, 0, v4, nil) + ":24",
		}, "|")},
		{0, 0, 0, 0, -1, 0, nil, nil, append(slices.Clone(base), 128, 52, 18), strings.Join([]string{
			"0", "", "true:true::", "true:" + ip(6, hi, lo, []byte{128, 52, 18}),
			"0,0", "false:false:0,0:", "true:" + ip(6, hi, lo, []byte{128}) + ":4660",
			"255", "false:false:255:", "true:" + ip(6, hi, lo, nil) + ":18",
		}, "|")},
	}
	reject := func(row, field int, output string) {
		if strings.Split(cases[row].Want, "|")[field] == output {
			panic(fmt.Sprintf("mutation survived on row %d field %d", row, field))
		}
	}
	reject(0, 2, appended(false, false, append([]byte{91}, a...), old))
	reject(0, 0, strconv.Itoa(len(base)))
	reject(1, 4, "1,2,3,4,18,52")
	reject(0, 7, csv(append(slices.Clone(a), 65)))
	reject(0, 9, "false:"+ip(6, hi, lo, nil)+":65")
	reject(1, 9, "true:"+ip(4, 0, v4, nil)+":18")
	reject(2, 3, "true:"+ip(6, hi, lo, []byte{0, 0, 0}))
	fmt.Println("detected one-step append, dropped size annotation, big-endian port, preserved prefix zone, rejected invalid bits, non-atomic decode and borrowed zone mutations")

	file, err := os.CreateTemp("", "review-netip-binary-")
	must(err)
	defer os.Remove(file.Name())
	encoder := json.NewEncoder(file)
	var protocol bytes.Buffer
	for _, c := range cases {
		tag := 0
		if c.Backing == nil {
			tag = 1
		}
		fmt.Fprint(&protocol, c.Family, " ", c.Hi, " ", c.Lo, " ", c.Port, " ", c.Bits, " ", len(c.Zone), " ", c.Length, " ", len(c.Backing), " ", tag)
		for _, part := range [][]byte{c.Zone, c.Backing, c.Payload} {
			for _, b := range part {
				fmt.Fprint(&protocol, " ", b)
			}
		}
		fmt.Fprintln(&protocol)
		must(encoder.Encode(c))
	}
	must(file.Close())
	for i, row := range check.LeanRows("godiff-netip-binary-codec", len(cases), &protocol) {
		if row != cases[i].Want {
			panic(fmt.Sprintf("row %d got %q want %q", i, row, cases[i].Want))
		}
	}
	must(os.Setenv("GOLEAN_NETIP_BINARY_CASES", file.Name()))
	must(check.TestInPackage("net/netip", "TestGoLeanBinaryCodec", "difftest/net/netip/binary-codec/testdata/binary_codec_test.go"))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
