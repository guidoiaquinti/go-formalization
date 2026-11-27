package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type observation struct {
	mode, creator, attrs, descriptor, dos, unix, toUnix uint32
}

type fixture struct {
	creator, attrs, mode, flags, slash uint32
	want                               observation
}

func main() {
	const (
		dir     = uint32(fs.ModeDir)
		link    = uint32(fs.ModeSymlink)
		char    = uint32(fs.ModeDevice | fs.ModeCharDevice)
		special = uint32(fs.ModeSetuid | fs.ModeSetgid | fs.ModeSticky)
	)
	cases := []fixture{
		{0xffab, 17, dir | link | 0o444, 8, 1,
			observation{dir, 0x3ab, 0x81240011, 1, dir | 0o555, 17, 0x8124}},
		{0x3ec, 0xafa40010, char | uint32(fs.ModeSetgid) | 0o640, 9, 0,
			observation{link | special | 0o644, 0x3ec, 0x25a00000, 1, dir | 0o777, 16, 0x25a0}},
		{0x13f1, 0x20000001, uint32(fs.ModeNamedPipe) | 0o022, 16, 1,
			observation{dir | char, 0x3f1, 0x10120001, 0, 0o444, 1, 0x1012}},
		{0x0007, 0xf001, 0o200, 0x8000, 0,
			observation{0o444, 0x307, 0x80800000, 0, 0o444, 1, 0x8080}},
	}
	native, err := os.CreateTemp("", "golean-zip-mode-review-")
	must(err)
	defer os.Remove(native.Name())
	var input bytes.Buffer
	for _, c := range cases {
		data := fmt.Sprintf("%d %d %d %d %d", c.creator, c.attrs, c.mode, c.flags, c.slash)
		fmt.Fprintln(&input, data)
		_, err := fmt.Fprintln(native, data, render(c.want))
		must(err)
	}
	must(native.Close())

	var rejected [6]int
	for i, got := range check.LeanRows("godiff-zip-file-mode", len(cases), &input) {
		c := cases[i]
		if want := render(c.want); got != want {
			log.Fatalf("row %d: Lean %q, want %q", i, got, want)
		}
		for j, mutant := range mutants(c) {
			if render(mutant) != got {
				rejected[j]++
			}
		}
	}
	if want := [6]int{2, 4, 2, 3, 1, 1}; rejected != want {
		log.Fatalf("mutation rejections %v, want %v", rejected, want)
	}
	must(os.Setenv("GOLEAN_ZIP_MODE_ROWS", native.Name()))
	must(check.TestInPackage("archive/zip", "TestGoLeanFileMode", "difftest/archive/zip/file-mode/testdata/file_mode_test.go"))
	fmt.Printf("ZIP metadata: %d pinned native and Lean rows, mutation rejections %v\n", len(cases), rejected)
}

func render(o observation) string {
	return fmt.Sprintf("%d %d %d %d %d %d %d", o.mode, o.creator, o.attrs, o.descriptor, o.dos, o.unix, o.toUnix)
}

func mutants(c fixture) [6]observation {
	o := c.want
	wrong := [6]observation{o, o, o, o, o, o}
	if c.slash != 0 {
		wrong[0].mode &^= uint32(fs.ModeDir)
	}
	wrong[1].creator &= 0xff00
	wrong[2].descriptor = 0
	if c.attrs&1 != 0 {
		wrong[3].dos |= 0o222
	}
	wrong[4].unix |= c.attrs & 0xf000
	if c.mode&uint32(fs.ModeDir) != 0 {
		wrong[5].toUnix = o.toUnix&0x0fff | 0x4000
	}
	return wrong
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
