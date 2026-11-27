# Independent review of `compress/flate.hash4u`

**Verdict: FAITHFUL**

Reviewed author commit `f161012` against pinned Go 1.27.0 `src/compress/flate/deflate.go:335–337`. The multiplication wraps at 32 bits before the shift. The subtraction for the shift wraps in uint8, so a requested width above 32 yields zero.

Ten independent literal fixtures passed against the pinned native entry and compiled Lean. Saturating the underflow instead changed four fixtures. The runnable reviewer program checks the exact rejection counts. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/compress/flate/hash-words/review` from the repository root.

Author evidence passed 66,240 boundaries and 60,000 samples at seeds 17, 2718, and 65537. Scope is the typed integer operations. Hash-chain and compressor effects remain outside scope.
Evidence SHA-256: `3c4b63bf0f20967da5d9b4d724d43b3c1628c3ddcf8fa6440f5bb1725b9db3ce`
