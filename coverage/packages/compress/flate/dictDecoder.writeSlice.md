# Independent review of `compress/flate.dictDecoder.writeSlice`

**Verdict: FAITHFUL**

The returned view begins at the writer, including when the visible suffix is empty. The theorem checks the origin. The runnable overlay writes through visible and spare cells and checks that a receiver-backing mutation is visible in the returned view.

The scoped Lean axiom audit passed and the pinned Go 1.27.0 and Lean comparison replayed 991 fixed rows and 60,000 sampled rows at seeds 17, 2718 and 65537 on darwin/arm64. The author scope covers nonnil decoder state with ordered cursors, legal buffer geometry, and in-range mark counts. The model represents slice origin and geometry, with backing identity checked in native Go. Panic text and concurrent mutation are outside scope.

Seven literal controls check both pinned native Go and Lean, including nil versus nonnil empty views. Clipping write capacity to the visible suffix disagrees on 3 rows. Wrapping the flush at capacity or allowing a write into spare capacity each disagrees on 2. A separate same-package overlay checks bidirectional aliasing and failed-write atomicity on one stateful trace.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/compress/flate/dictionary-state/review
```

Evidence SHA-256: `41ce726d23262fce35e56c6f7c98a5fbdeecbe813680506c074aff96264c5e27`
