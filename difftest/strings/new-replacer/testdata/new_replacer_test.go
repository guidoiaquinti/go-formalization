package strings_test

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func observe(args []string) (result string) {
	defer func() {
		if failure := recover(); failure != nil {
			if text, ok := failure.(string); ok && text == "strings.NewReplacer: odd argument count" {
				result = text
			} else {
				panic(failure)
			}
		}
	}()
	r := strings.NewReplacer(args...)
	// Expose a slice alias if the constructor retained the caller's slots.
	for i := range args {
		args[i] += "\x00"
	}
	oldnew, once, noImplementation := strings.GoLeanInspectNewReplacer(r)
	storage := "nil"
	if oldnew != nil {
		storage = "slice"
	}
	fields := make([]string, len(oldnew))
	for i, arg := range oldnew {
		wire := strconv.AppendInt(nil, int64(len(arg)), 10)
		wire = append(wire, ':')
		for j := range len(arg) {
			if j != 0 {
				wire = append(wire, ',')
			}
			wire = strconv.AppendUint(wire, uint64(arg[j]), 10)
		}
		fields[i] = string(wire)
	}
	return fmt.Sprintf("%t:%t:%s:%s", reflect.ValueOf(once).Elem().IsZero(), noImplementation, storage, strings.Join(fields, "|"))
}

func TestGoLeanNewReplacer(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_NEW_REPLACER_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows := json.NewDecoder(file)
	for i := 0; ; i++ {
		var row struct {
			Args [][]byte
			Want string
		}
		if err := rows.Decode(&row); err == io.EOF {
			if i == 0 {
				t.Fatal("missing constructor fixtures")
			}
			t.Logf("%d native constructor observations agree", i)
			break
		} else if err != nil {
			t.Fatal(err)
		}
		var args []string
		if row.Args != nil {
			args = make([]string, len(row.Args))
		}
		for j, raw := range row.Args {
			args[j] = string(raw)
		}
		if got := observe(args); got != row.Want {
			t.Fatalf("row %d original arguments=%x: Go=%s want=%s", i, row.Args, got, row.Want)
		}
	}
}
