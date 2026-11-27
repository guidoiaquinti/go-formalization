package constant

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanScalars(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_CONSTANT_SCALARS"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	rows := 0
	for scanner.Scan() {
		input, want, found := strings.Cut(scanner.Text(), "\t")
		if !found {
			t.Fatal("missing expected row")
		}
		var code, variant int
		var flag bool
		if _, err := fmt.Sscan(input, &code, &flag, &variant); err != nil {
			t.Fatal(err)
		}
		u, ok := MakeUnknown().(unknownVal)
		if !ok {
			t.Fatal("MakeUnknown returned a different representation")
		}
		b, ok := MakeBool(flag).(boolVal)
		if !ok {
			t.Fatal("MakeBool returned a different representation")
		}
		var operand Value
		switch variant {
		case 0:
			operand = u
		case 1:
			operand = b
		case 3:
			operand = MakeString("x")
		case 4:
			operand = MakeInt64(int64(code))
		}
		observed := func() (result string) {
			defer func() {
				if fault := recover(); fault != nil {
					text, ok := fault.(string)
					if !ok {
						panic(fault)
					}
					result = "panic:" + text
				}
			}()
			return strconv.FormatBool(BoolVal(operand))
		}()
		got := fmt.Sprintf("unknownVal|boolVal:%t|%d|%s|%s|%d|%s|%s|%s|%s",
			bool(b), int(u.Kind()), u.String(), u.ExactString(), int(b.Kind()), b.String(),
			b.ExactString(), observed, Kind(code).String())
		if got != want {
			t.Fatalf("row %d (%s): Go %q want %q", rows, input, got, want)
		}
		rows++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if rows == 0 {
		t.Fatal("missing scalar observations")
	}
	t.Logf("%d pinned native observations", rows)
}
