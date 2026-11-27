package strings_test

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanSplitSeqPrivate(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_SPLIT_SEQ_PRIVATE_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	rows := json.NewDecoder(f)
	for i := 0; ; i++ {
		var row struct {
			Input, Sep          []byte
			Save, First, Second int
			Want                string
		}
		if err := rows.Decode(&row); err == io.EOF {
			if i == 0 {
				t.Fatal("missing iterator fixtures")
			}
			t.Logf("%d configured iterator traces agree", i)
			break
		} else if err != nil {
			t.Fatal(err)
		}
		seq := strings.GoLeanSplitSeqPrivate(string(row.Input), string(row.Sep), row.Save)
		run := func(budget int) (result string) {
			var parts []string
			defer func() {
				if fault := recover(); fault != nil {
					err, ok := fault.(runtime.Error)
					if !ok || !strings.HasPrefix(err.Error(), "runtime error: slice bounds out of range") {
						panic(fault)
					}
					result = "p:" + strings.Join(parts, ";")
				}
			}()
			if budget > 0 {
				seq(func(part string) bool {
					encoded := []byte{'f'}
					for j := range len(part) {
						if j != 0 {
							encoded = append(encoded, ',')
						}
						encoded = strconv.AppendUint(encoded, uint64(part[j]), 10)
					}
					parts = append(parts, string(encoded))
					return len(parts) < budget
				})
			}
			return "o:" + strings.Join(parts, ";")
		}
		result := fmt.Sprintf("%s|%s|%s|%s", run(row.First), run(row.Second), run(len(row.Input)+2), run(1))
		if result != row.Want {
			t.Fatalf("row %d input=%x sep=%x save=%d: Go=%s want=%s", i, row.Input, row.Sep, row.Save, result, row.Want)
		}
	}
}
