# Independent review of `compress/flate.emitLiterals`

**Verdict: FAITHFUL**

Evidence SHA-256: `4ab4228cd83a8595ac481077b4c6ab17682f13c198bf22551bd6cba91a32ff71`

I reviewed evidence commit `10d3a12` against pinned Go 1.27.0 `src/compress/flate/token.go`. The source emits in input order using the same updates as AddLiteral. The fixed control writes the last slot then wraps into slots zero and one and counts a repeated byte twice. An empty input preserves the complete receiver.

Pinned native Go and Lean agreed on 280 fixed rows and 20,000 rows per seed for 17, 2718 and 65537 on darwin/arm64 using `go run ./difftest/compress/flate/token-state`. Native fixed rows inspect the full token storage and sampled rows inspect touched words and sentinels. Lean observes the reported subset. Both compare the modeled state with the independent oracle, and native replay checks complete histograms. The independent native state assertions and three separate Lean rows in `go run ./difftest/compress/flate/token-state/review` reject ten wrong outcomes. Scope covers nonnil receivers, disjoint emitted bytes and reslicing to capacity. The view aliases the same receiver's backing. Concurrent writes and nil receivers are outside scope. [Independent control](../../../../difftest/compress/flate/token-state/review/main.go)
