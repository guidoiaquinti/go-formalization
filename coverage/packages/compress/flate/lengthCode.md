# Independent review of `compress/flate.lengthCode`

**Verdict: FAITHFUL**

Evidence SHA-256: `8b7913f510d70a4e50cf34a2a7a48f7032620f49f1b40b4ac0017e65a9d3ce02`

I checked evidence commit `a028017` against pinned Go 1.27.0 `src/compress/flate/token.go`. Pinned Go directly indexes the length table with a byte. The model identifies the same intervals from their starting points, with code 28 reserved for byte 255. All 256 inputs are checked independently. The registered theorems cover the first paired interval and the last byte. The fixed control rejects both advancing at byte 9 and merging byte 255 into the preceding interval.

I ran `go run ./difftest/compress/flate/match-codes` on darwin/arm64. The pinned private Go comparison and Lean matched all 65,536 low-16-bit offsets, with all length bytes represented, plus 20,000 cases per seed for 17, 2718 and 65537. Eight independent literal rows passed with `go run ./difftest/compress/flate/match-codes/review` and rejected six deliberately incorrect observations. Match construction and compression output are outside scope. [Independent control](../../../../difftest/compress/flate/match-codes/review/main.go)
