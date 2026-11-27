# Independent review of `compress/flate.dictDecoder.tryWriteCopy`

**Verdict: FAITHFUL**

The fast method checks for source wrap and insufficient visible output space before copying. Both rejection theorems establish exact state preservation. A nonwrapping request that fits the visible window uses forward overlapping copies and advances by the full request. The controls check the exact fit and a too-long request with spare backing.

The scope is nonnil decoder state with legal ordered cursors and positive distances within available history. Counts are nonnegative and feasible without native int overflow. The pinned Go overlay and Lean observe counts, cursors, full flags, and complete backing including bytes beyond the visible window. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 7,644 fixed cases and 60,000 samples at seeds 17, 2718 and 65537 on darwin/arm64. Seven independent literal controls passed. Reading only the initial backing during overlap disagrees on 3 cases. Letting a wrapped request perform the ordinary copy in the fast method disagrees on 3, and letting a nonfitting request truncate in the fast method disagrees on 1.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/compress/flate/dictionary-copy/review
```

Evidence SHA-256: `39e86afb8f14b8013b706e7ef66be4571eeb9b17257af1e2f7c72bde044101b9`
