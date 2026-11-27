# Independent review of `compress/flate.dictDecoder.readFlush`

**Verdict: FAITHFUL**

The view starts at the old reader and ends at the old writer with spare backing capacity preserved. The source moves the reader to the writer and resets both cursors with full set at the end of the visible window. The theorem establishes the wrap state. Independent native controls check both alias directions as well as the payload and wrap.

The scoped Lean axiom audit passed and the pinned Go 1.27.0 and Lean comparison replayed 991 fixed rows and 60,000 sampled rows at seeds 17, 2718 and 65537 on darwin/arm64. The author scope covers nonnil decoder state with ordered cursors, legal buffer geometry, and in-range mark counts. The model represents slice origin and geometry, with backing identity checked in native Go. Panic text and concurrent mutation are outside scope.

Seven literal controls check both pinned native Go and Lean, including nil versus nonnil empty views. Clipping write capacity to the visible suffix disagrees on 3 rows. Wrapping the flush at capacity or allowing a write into spare capacity each disagrees on 2. A separate same-package overlay checks bidirectional aliasing and failed-write atomicity on one stateful trace.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/compress/flate/dictionary-state/review
```

Evidence SHA-256: `58916ed859834e1f72a6056dea6237dbed3f178c6d88f76ca3c5b18d782d76cd`
