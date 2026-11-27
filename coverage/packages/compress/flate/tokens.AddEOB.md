# Independent review of `compress/flate.tokens.AddEOB`

**Verdict: FAITHFUL**

Evidence SHA-256: `274828c33d02690e76413989308726373ae64e5e8b024fc8da03c4450e1e6667`

I reviewed evidence commit `10d3a12` against pinned Go 1.27.0 `src/compress/flate/token.go`. The source writes token 256 at the old count and increments extra bin zero and count as uint16. Control starts with a saturated extra bin and verifies that literal and offset histograms remain unchanged.

Pinned native Go and Lean agreed on 280 fixed rows and 20,000 rows per seed for 17, 2718 and 65537 on darwin/arm64 using `go run ./difftest/compress/flate/token-state`. Native fixed rows inspect the full token storage and sampled rows inspect touched words and sentinels. Lean observes the reported subset. Both compare the modeled state with the independent oracle, and native replay checks complete histograms. The independent native state assertions and three separate Lean rows in `go run ./difftest/compress/flate/token-state/review` reject ten wrong outcomes. Scope covers nonnil receivers, disjoint emitted bytes and reslicing to capacity. The view aliases the same receiver's backing. Concurrent writes and nil receivers are outside scope. [Independent control](../../../../difftest/compress/flate/token-state/review/main.go)
