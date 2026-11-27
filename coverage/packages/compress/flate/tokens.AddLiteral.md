# Independent review of `compress/flate.tokens.AddLiteral`

**Verdict: FAITHFUL**

Evidence SHA-256: `3d9843c55d21c80aab34a6ad8cb185dc207f02e589d3c171952ebce3088f1b3f`

I reviewed evidence commit `10d3a12` against pinned Go 1.27.0 `src/compress/flate/token.go`. The source writes at the old uint16 count, increments the selected literal histogram and then increments the count. A count of 65,535 and saturated selected bin both wrap, while other backing words and histograms retain their contents.

Pinned native Go and Lean agreed on 280 fixed rows and 20,000 rows per seed for 17, 2718 and 65537 on darwin/arm64 using `go run ./difftest/compress/flate/token-state`. Native fixed rows inspect the full token storage and sampled rows inspect touched words and sentinels. Lean observes the reported subset. Both compare the modeled state with the independent oracle, and native replay checks complete histograms. The independent native state assertions and three separate Lean rows in `go run ./difftest/compress/flate/token-state/review` reject ten wrong outcomes. Scope covers nonnil receivers, disjoint emitted bytes and reslicing to capacity. The view aliases the same receiver's backing. Concurrent writes and nil receivers are outside scope. [Independent control](../../../../difftest/compress/flate/token-state/review/main.go)
