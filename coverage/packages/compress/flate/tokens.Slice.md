# Independent review of `compress/flate.tokens.Slice`

**Verdict: FAITHFUL**

Evidence SHA-256: `c492ac878a49c00e49530a51004ee8e7c0e9a3a402321f7136f363dc9ebed4ca`

I reviewed evidence commit `10d3a12` against pinned Go 1.27.0 `src/compress/flate/token.go`. The source returns a nonnil prefix of the fixed array with the count as length and capacity 65,536. The native control writes through a resliced zero-length view and separately verifies that the old view observes a later receiver append while retaining its old length.

Pinned native Go and Lean agreed on 280 fixed rows and 20,000 rows per seed for 17, 2718 and 65537 on darwin/arm64 using `go run ./difftest/compress/flate/token-state`. Native fixed rows inspect the full token storage and sampled rows inspect touched words and sentinels. Lean observes the reported subset. Both compare the modeled state with the independent oracle, and native replay checks complete histograms. The independent native state assertions and three separate Lean rows in `go run ./difftest/compress/flate/token-state/review` reject ten wrong outcomes. Scope covers nonnil receivers, disjoint emitted bytes and reslicing to capacity. The view aliases the same receiver's backing. Concurrent writes and nil receivers are outside scope. [Independent control](../../../../difftest/compress/flate/token-state/review/main.go)
