package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"reflect"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	edges := []string{"", "abc", "\x00\xff", "猫", "🚀猫"}
	for b := range 256 {
		edges = append(edges, string([]byte{byte(b)}))
	}
	check.CompareSamples("godiff-strings-new-reader", edges, func(rng *rand.Rand, _ int) string {
		data := make([]byte, rng.Intn(97))
		rng.Read(data)
		return string(data)
	}, func(s string) string {
		return describe("0", s)
	}, actual, func(s string) string {
		return describe("0 -1", s)
	})
}

func actual(s string) string {
	reader := strings.NewReader(s)
	if reader == nil {
		log.Fatal("NewReader returned nil")
	}
	// Read the initial fields before the public read changes them.
	fields := reflect.ValueOf(reader).Elem()
	cursor := fields.FieldByName("i").Int()
	previousRune := fields.FieldByName("prevRune").Int()
	contents, err := io.ReadAll(reader)
	if err != nil {
		log.Fatalf("reading from new reader: %v", err)
	}
	return describe(fmt.Sprintf("%d %d", cursor, previousRune), string(contents))
}

func describe(prefix, data string) string {
	row := []byte(prefix)
	for i := range len(data) {
		row = strconv.AppendUint(append(row, ' '), uint64(data[i]), 10)
	}
	return string(row)
}
