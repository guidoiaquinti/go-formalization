# Independent review of `compress/flate.dictDecoder.writeMark`

**Verdict: FAITHFUL**

The operation advances only the writer by an in-range count. The theorem checks the resulting writer directly. Fixed controls preserve the read cursor, full flag and entire backing when advancing to or stopping before the end.

The scoped Lean axiom audit passed and the pinned Go 1.27.0 and Lean comparison replayed 991 fixed rows and 60,000 sampled rows at seeds 17, 2718 and 65537 on darwin/arm64. The author scope covers nonnil decoder state with ordered cursors, legal buffer geometry, and in-range mark counts. The model represents slice origin and geometry, with backing identity checked in native Go. Panic text and concurrent mutation are outside scope.

Seven literal controls check both pinned native Go and Lean, including nil versus nonnil empty views. Clipping write capacity to the visible suffix disagrees on 3 rows. Wrapping the flush at capacity or allowing a write into spare capacity each disagrees on 2. A separate same-package overlay checks bidirectional aliasing and failed-write atomicity on one stateful trace.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/compress/flate/dictionary-state/review
```

Evidence SHA-256: `4749f5c71efe122e56b8800e25e52e788ab5f9580e6829a82cd7e992c35ac9eb`
