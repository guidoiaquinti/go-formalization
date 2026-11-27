# Independent review of `regexp/syntax.appendClass`

**Verdict: FAITHFUL**

Reviewed author commit `765c316` against pinned Go 1.27.0. The independent fixtures check both the empty class identity and a merge through the next-to-last destination pair.

Pinned Go and Lean replayed 5,040 fixed and 60,000 sampled inputs, using seeds 17, 2718, and 65537. Of these, 55,695 rows exercised valid indices for Less and Swap. Four independent literal fixtures exercised the registered properties and rejected eight deliberately wrong outcomes. The scope uses pair-aligned storage and separated, disjoint source ranges for operations that consume a clean class. Growth capacity and zero-capacity backing identity are outside scope.

Run the [reviewer control](../../../../difftest/regexp/syntax/range-class/review/main.go) with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/regexp/syntax/range-class/review`.

Evidence SHA-256: `6153b51f090f970359710930d64addbea8385e21eb55e78a926f7e0a9888067e`
