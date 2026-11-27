# Independent review of `regexp/syntax.ranges.Less`

**Verdict: FAITHFUL**

Reviewed author commit `765c316` against pinned Go 1.27.0. Equal lower endpoints sort by descending upper endpoint. The control checks the pair with upper endpoint 11 before the pair ending in 9.

Pinned Go and Lean replayed 5,040 fixed and 60,000 sampled inputs, using seeds 17, 2718, and 65537. Of these, 55,695 rows exercised valid indices for Less and Swap. Four independent literal fixtures exercised the registered properties and rejected eight deliberately wrong outcomes. The scope uses pair-aligned storage and separated, disjoint source ranges for operations that consume a clean class. Growth capacity and zero-capacity backing identity are outside scope.

Run the [reviewer control](../../../../difftest/regexp/syntax/range-class/review/main.go) with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/regexp/syntax/range-class/review`.

Evidence SHA-256: `9aaae6e03188babd0e5259b0f18a19b31031429177b0a397de95c6aa3b9e4a2e`
