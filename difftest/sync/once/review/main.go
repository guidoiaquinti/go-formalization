package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Tag, Mode, Second int
	Payload, X, Y     int64
	Want              string
}

func main() {
	const poisoned = "panic:42:true:true,ok:false:true"
	const missing = "nilCallback:false:true,ok:false:true"
	const complete = "ok:false:true,ok:false:true"
	const normal = "ok:true:true,ok:false:true"
	const nilPanic = "nilPanic:true:true,ok:false:true"
	const noReceiver = "nilReceiver:false:nil,nilReceiver:false:nil"
	joint := func(plain, f, v, vv string) string {
		return strings.Join([]string{plain, plain, f, v, vv}, "|")
	}
	cases := []fixture{
		{1, 2, 1, 42, -2, 3, joint(poisoned, "panic:42:true,panic:42:false", "panic:42:true,panic:42:false", "panic:42:true,panic:42:false")},
		{1, 0, 2, 42, -2, 3, joint(missing, "nilCallback:false,nilCallback:false", "nilCallback:false,nilCallback:false", "nilCallback:false,nilCallback:false")},
		{2, 0, 2, 42, -2, 3, joint(complete, "nilCallback:false,nilCallback:false", "nilCallback:false,nilCallback:false", "nilCallback:false,nilCallback:false")},
		{0, 1, 1, 0, -2, 3, joint(noReceiver, "ok:true,ok:false", "-2:true,-2:false", "-2/3:true,-2/3:false")},
		{1, 1, 3, 0, -2, 3, joint(normal, "ok:true,ok:false", "-2:true,-2:false", "-2/3:true,-2/3:false")},
		{1, 3, 1, 0, -2, 3, joint(nilPanic, "nilPanic:true,nilPanic:false", "nilPanic:true,nilPanic:false", "nilPanic:true,nilPanic:false")},
	}
	var input bytes.Buffer
	file, err := os.CreateTemp("", "review-sync-once-")
	must(err)
	defer os.Remove(file.Name())
	enc := json.NewEncoder(file)
	for _, c := range cases {
		fmt.Fprintln(&input, c.Tag, c.Mode, c.Second, c.Payload, c.X, c.Y)
		must(enc.Encode(c))
	}
	must(file.Close())
	must(os.Setenv("GOLEAN_SYNC_ONCE_CASES", file.Name()))
	must(check.TestInPackage("sync", "TestGoLeanSyncOnce", "difftest/sync/once/testdata/bridge_test.go",
		"difftest/sync/once/testdata/once_test.go"))

	rows := check.LeanRows("godiff-sync-once", len(cases), &input)
	for i, row := range rows {
		if row != cases[i].Want {
			log.Fatalf("row %d got %q, want %q", i, row, cases[i].Want)
		}
	}
	reject := func(row, column int, mutant string) {
		if strings.Split(rows[row], "|")[column] == mutant {
			log.Fatalf("row %d column %d did not reject %q", row, column, mutant)
		}
	}
	reject(0, 0, "panic:42:true:false,ok:true:true")           // retry after panic
	reject(1, 1, "nilCallback:false:false,panic:42:true:true") // fail to publish completion on nil call
	reject(2, 0, "nilCallback:false:true,panic:42:true:true")  // evaluate callback after completion
	reject(0, 2, "panic:42:true,ok:false")                     // swallow cached panic
	reject(5, 2, "nilPanic:true,panic:0:false")                // replace panic(nil) outcome
	reject(4, 3, "-2:true,-2:true")                            // recompute cached value
	reject(4, 4, "-2/3:true,-2/0:false")                       // lose second cached value
	log.Print("six independent pinned Go and Lean cases passed and rejected seven wrong observations")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
