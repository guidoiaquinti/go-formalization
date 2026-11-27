# Independent review of `compress/flate.offsetCode`

**Verdict: FAITHFUL**

Evidence SHA-256: `a9415583e8f23e301a8fcf107eaea8183e0652adc5ae1dee2f9129e4fd15dc92`

I checked evidence commit `a028017` against pinned Go 1.27.0 `src/compress/flate/token.go`. Go uses the short table for offsets below 256. The long branch shifts by seven and narrows to a byte before selecting a second table. The theorem covers the branch boundary and the long-branch period of 32,768. Beyond the DEFLATE window, these results describe the pinned helper's representation only. For example, offset 32,768 returns 14, since the implementation selects the long branch before applying the narrowed index.

I ran `go run ./difftest/compress/flate/match-codes` on darwin/arm64. The pinned private Go comparison and Lean matched all 65,536 low-16-bit offsets, with all length bytes represented, plus 20,000 cases per seed for 17, 2718 and 65537. Eight independent literal rows passed with `go run ./difftest/compress/flate/match-codes/review` and rejected six deliberately incorrect observations. Match construction and compression output are outside scope. [Independent control](../../../../difftest/compress/flate/match-codes/review/main.go)
