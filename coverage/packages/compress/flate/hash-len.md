# Independent review of `compress/flate.hashLen`

**Verdict: FAITHFUL**

Reviewed author commit `f161012` against pinned Go 1.27.0 `src/compress/flate/deflatefast.go:161–176`. Count three uses its own 32-bit multiplier. Counts five through eight shift a 64-bit wrapped product and then retain its low 32 bits. Other counts select the four-byte branch.

Ten independent literal fixtures passed against the pinned native entry and compiled Lean. Returning zero for unrecognized counts changed one fixture. Selecting the high 32 bits of the wide shifted result changed three. The runnable reviewer program checks the exact rejection counts. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/compress/flate/hash-words/review` from the repository root.

Author evidence passed 66,240 boundaries and 60,000 samples at seeds 17, 2718, and 65537. Scope is the typed integer operations. Hash-chain and compressor effects remain outside scope.
Evidence SHA-256: `2794fcf52f05bb3d19d55cfab37c01c5bb7b3357df5e39d4c8a54c9847dfcbd6`
