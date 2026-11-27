package zip

import (
	"encoding/json"
	"io"
	"os"
	"testing"
)

func TestGoLeanDetectUTF8(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_ZIP_DETECT_UTF8_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	dec := json.NewDecoder(file)
	count := 0
	for {
		var row struct {
			Raw            []byte
			Valid, Require bool
		}
		if err := dec.Decode(&row); err == io.EOF {
			break
		} else if err != nil {
			t.Fatal(err)
		}
		valid, require := detectUTF8(string(row.Raw))
		if valid != row.Valid || require != row.Require {
			t.Fatalf("%x: got (%t,%t), want (%t,%t)", row.Raw, valid, require, row.Valid, row.Require)
		}
		count++
	}
	if count == 0 {
		t.Fatal("no UTF-8 observations")
	}
	t.Logf("%d native observations", count)
}
