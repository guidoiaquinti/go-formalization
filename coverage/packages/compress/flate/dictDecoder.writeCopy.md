# Independent review of `compress/flate.dictDecoder.writeCopy`

**Verdict: FAITHFUL**

The initial Go copy reads the ring suffix when the source wraps. The later copies can use newly written bytes, and the endpoint is limited by the visible window. Lean expresses the same sequence as forward byte copies at a circular distance. The theorem proves the clipped count. A full-ring fixture also checks that copying a window to itself advances the writer with unchanged bytes.

The scope is nonnil decoder state with legal ordered cursors and positive distances within available history. Counts are nonnegative and feasible without native int overflow. The pinned Go overlay and Lean observe counts, cursors, full flags, and complete backing including bytes beyond the visible window. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 7,644 fixed cases and 60,000 samples at seeds 17, 2718 and 65537 on darwin/arm64. Seven independent literal controls passed. Reading only the initial backing during overlap disagrees on 3 cases. Letting a wrapped request perform the ordinary copy in the fast method disagrees on 3, and letting a nonfitting request truncate in the fast method disagrees on 1.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/compress/flate/dictionary-copy/review
```

Evidence SHA-256: `4cc7e42f1ad6cdb208ed5f2b58a217dedcf4bf02b59ce33d290d14fe814f0437`
