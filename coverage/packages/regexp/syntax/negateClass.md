# Independent review of `regexp/syntax.negateClass`

**Verdict: FAITHFUL**

Reviewed author commit `765c316` against pinned Go 1.27.0. Negating the empty class yields a single inclusive zero-through-MaxRune range. The partial case checks rewriting within the old backing followed by growth.

Pinned Go and Lean replayed 5,040 fixed and 60,000 sampled inputs, using seeds 17, 2718, and 65537. Of these, 55,695 rows exercised valid indices for Less and Swap. Four independent literal fixtures exercised the registered properties and rejected eight deliberately wrong outcomes. The scope uses pair-aligned storage and separated, disjoint source ranges for operations that consume a clean class. Growth capacity and zero-capacity backing identity are outside scope.

Run the [reviewer control](../../../../difftest/regexp/syntax/range-class/review/main.go) with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/regexp/syntax/range-class/review`.

Evidence SHA-256: `f5179d4cbc9f3089b9113170af651641991935c3fb48a36f9b3e16ba7d4ed263`
