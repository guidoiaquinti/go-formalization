package zip

import (
	"math"
	"strings"
	"testing"
)

func TestGoLeanWriterConfigReview(t *testing.T) {
	cw := &countWriter{}
	w := &Writer{cw: cw, comment: "keep"}
	w.SetOffset(0)
	w.SetOffset(math.MinInt64)
	func() {
		defer func() {
			if got := recover(); got != "zip: SetOffset called after data was written" {
				t.Fatalf("second offset panic %v", got)
			}
		}()
		w.SetOffset(0)
	}()
	if w.cw != cw || cw.count != math.MinInt64 || w.comment != "keep" {
		t.Fatal("nonzero offset guard changed state")
	}
	exact := strings.Repeat("a", 65533) + "µ"
	over := exact + "x"
	if err := w.SetComment(exact); err != nil || w.comment != exact {
		t.Fatal("65535 bytes rejected or changed")
	}
	first, second := w.SetComment(over), w.SetComment(over)
	if first == nil || second == nil || first == second || first.Error() != "zip: Writer.Comment too long" || second.Error() != first.Error() ||
		w.comment != exact || w.cw != cw || cw.count != math.MinInt64 {
		t.Fatal("oversize guard mutated state or error behavior")
	}
	var nilWriter *Writer
	if err := nilWriter.SetComment(over); err == nil || err.Error() != first.Error() {
		t.Fatal("comment length must be checked before nil receiver access")
	}
}
