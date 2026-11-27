package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	want := []string{
		"BigEndian", "binary.BigEndian",
		"LittleEndian", "binary.LittleEndian",
		"NativeEndian", "binary.NativeEndian",
	}
	direct := []string{
		binary.BigEndian.String(), binary.BigEndian.GoString(),
		binary.LittleEndian.String(), binary.LittleEndian.GoString(),
		binary.NativeEndian.String(), binary.NativeEndian.GoString(),
	}
	quoted := []string{
		fmt.Sprintf("%q", binary.BigEndian), fmt.Sprintf("%#v", binary.BigEndian),
		fmt.Sprintf("%q", binary.LittleEndian), fmt.Sprintf("%#v", binary.LittleEndian),
		fmt.Sprintf("%q", binary.NativeEndian), fmt.Sprintf("%#v", binary.NativeEndian),
	}
	out, err := exec.Command("lean/.lake/build/bin/godiff-binary-endian-names").Output()
	if err != nil {
		log.Fatal(err)
	}
	lean := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	if len(lean) != len(want) {
		log.Fatalf("got %d rows, want %d", len(lean), len(want))
	}
	for i, name := range want {
		formatted := name
		if i%2 == 0 {
			formatted = strconv.Quote(name)
		}
		if direct[i] != name || quoted[i] != formatted || lean[i] != name {
			log.Fatalf("method %d: Go=%q fmt=%q Lean=%q want %q/%q", i, direct[i], quoted[i], lean[i], name, formatted)
		}
	}
	nameLoss, packageLoss := 0, 0
	for _, name := range want {
		if name != strings.ReplaceAll(name, "NativeEndian", "LittleEndian") {
			nameLoss++
		}
		if name != strings.TrimPrefix(name, "binary.") {
			packageLoss++
		}
	}
	if nameLoss != 2 || packageLoss != 3 {
		log.Fatalf("mutation counts: native=%d prefix=%d", nameLoss, packageLoss)
	}
	fmt.Println("Six singleton Go and Lean method results passed, including quoted and Go-syntax formatting.")
}
