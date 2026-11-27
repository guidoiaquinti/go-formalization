# Independent review of `compress/flate.dictDecoder.writeByte`

**Verdict: FAITHFUL**

The write checks the visible length, including on a slice with spare capacity. The theorem establishes failure at the end and the independent native overlay checks that a failed write preserves both cursors, the full flag, slice geometry, alias and all bytes. A successful control confirms the specific byte and writer update.

The scoped Lean axiom audit passed and the pinned Go 1.27.0 and Lean comparison replayed 991 fixed rows and 60,000 sampled rows at seeds 17, 2718 and 65537 on darwin/arm64. The author scope covers nonnil decoder state with ordered cursors, legal buffer geometry, and in-range mark counts. The model represents slice origin and geometry, with backing identity checked in native Go. Panic text and concurrent mutation are outside scope.

Seven literal controls check both pinned native Go and Lean, including nil versus nonnil empty views. Clipping write capacity to the visible suffix disagrees on 3 rows. Wrapping the flush at capacity or allowing a write into spare capacity each disagrees on 2. A separate same-package overlay checks bidirectional aliasing and failed-write atomicity on one stateful trace.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/compress/flate/dictionary-state/review
```

Evidence SHA-256: `445d0e8d9fe74aefb65fe04aaec713dd2cda8702c326bf2ccca7e7d108e0c44d`
