# Independent review of `compress/flate.tokens.Reset`

**Verdict: FAITHFUL**

Evidence SHA-256: `aac4c4d821f0f05813009cf62e5f2fc7dea35cd37679dc1bf706d6965cafede9`

I reviewed evidence commit `10d3a12` against pinned Go 1.27.0 `src/compress/flate/token.go`. The source selects its no-op solely from zero token count, so a zero count retains nonzero histograms and filled count. For nonzero count the control checks all bins clear and every token word is retained.

Pinned native Go and Lean agreed on 280 fixed rows and 20,000 rows per seed for 17, 2718 and 65537 on darwin/arm64 using `go run ./difftest/compress/flate/token-state`. Native fixed rows inspect the full token storage and sampled rows inspect touched words and sentinels. Lean observes the reported subset. Both compare the modeled state with the independent oracle, and native replay checks complete histograms. The independent native state assertions and three separate Lean rows in `go run ./difftest/compress/flate/token-state/review` reject ten wrong outcomes. Scope covers nonnil receivers, disjoint emitted bytes and reslicing to capacity. The view aliases the same receiver's backing. Concurrent writes and nil receivers are outside scope. [Independent control](../../../../difftest/compress/flate/token-state/review/main.go)
