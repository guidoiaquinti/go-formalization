# Independent review of `compress/flate.bulkHash4`

**Verdict: FAITHFUL**

Evidence SHA-256: `8e00a8782b37b05db59a3932684364221f20205cb77c707bbc520825ef77fcbc`

I checked evidence commit `4fc8602` against Go 1.27.0 `src/compress/flate/deflate.go` and `src/compress/flate/load_store.go`. The source starts with the visible length guard, then updates the little-endian rolling word using each subsequent visible byte. Writes obey the destination length. The independent six-to-one row verifies that the first hash survives the subsequent bounds failure and the spare backing remains untouched. The five-to-two row checks that an extra source capacity byte does not create another window. The proofs cover the three-byte no-op and first write for a four-byte view.

Pinned native Go and Lean agreed on 6,006 fixed cases and 20,000 cases for each of seeds 17, 2718 and 65537 on darwin/arm64 using `go run ./difftest/compress/flate/hash-window`. Six independent literal rows passed with `go run ./difftest/compress/flate/hash-window/review` and rejected seven wrong observations. Scope covers disjoint source and destination in legal slice geometry, the full destination backing and the bounds failure class. The wording of panic messages and concurrent mutation are outside scope. [Independent control](../../../../difftest/compress/flate/hash-window/review/main.go)
