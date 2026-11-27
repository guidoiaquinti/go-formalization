package flate

import "testing"

func TestReviewGoLeanTokenState(t *testing.T) {
	var base tokens
	for i := range base.tokens {
		base.tokens[i] = token(0x10000 + i)
	}
	base.nFilled = 9
	base.litHist[5], base.extraHist[0], base.offHist[0] = 65535, 65535, 7
	unchanged := base
	unchanged.Reset()
	if unchanged != base {
		t.Fatal("zero-count reset touched state")
	}
	view := unchanged.Slice()
	if view == nil || len(view) != 0 || cap(view) != 65536 {
		t.Fatalf("zero-count view geometry: len=%d cap=%d", len(view), cap(view))
	}
	all := view[:cap(view)]
	all[0] = 13
	if unchanged.tokens[0] != 13 {
		t.Fatal("view write did not reach receiver")
	}
	unchanged.AddLiteral(5)
	if all[0] != 5 || len(view) != 0 {
		t.Fatal("old view did not retain length and observe receiver write")
	}

	base.n = 65535
	reset := base
	reset.Reset()
	if reset.n != 0 || reset.nFilled != 0 || reset.litHist != [256]uint16{} ||
		reset.extraHist != [32]uint16{} || reset.offHist != [32]uint16{} || reset.tokens != base.tokens {
		t.Fatal("nonzero reset did not clear state while keeping full backing")
	}
	literal := base
	literal.AddLiteral(5)
	if literal.n != 0 || literal.litHist[5] != 0 || literal.tokens[65535] != 5 ||
		literal.tokens[0] != base.tokens[0] || literal.extraHist != base.extraHist ||
		literal.offHist != base.offHist || literal.nFilled != base.nFilled {
		t.Fatal("literal append did not wrap or changed unrelated state")
	}
	end := base
	end.AddEOB()
	if end.n != 0 || end.extraHist[0] != 0 || end.tokens[65535] != 256 ||
		end.litHist != base.litHist || end.offHist != base.offHist || end.nFilled != base.nFilled {
		t.Fatal("EOB did not wrap or changed unrelated state")
	}
	sequence := base
	emitLiterals(&sequence, []byte{5, 7, 5})
	wantWords := base.tokens
	wantWords[65535], wantWords[0], wantWords[1] = 5, 7, 5
	wantHistogram := base.litHist
	wantHistogram[5], wantHistogram[7] = 1, 1
	if sequence.n != 2 || sequence.tokens != wantWords || sequence.litHist != wantHistogram ||
		sequence.extraHist != base.extraHist || sequence.offHist != base.offHist || sequence.nFilled != base.nFilled {
		t.Fatal("emit changed unrelated state or missed wrapped words and bins")
	}
	empty := base
	emitLiterals(&empty, nil)
	if empty != base {
		t.Fatal("emitting an empty slice changed state")
	}
}
