# Independent review of `bytes.Reader.UnreadRune`

**Verdict: FAITHFUL**

Evidence SHA-256: `cfa5940a3f70c90f97c4359b7d0ff60944a2f563d884ee8ada1e8ca98c2f03c0`

Go 1.27.0 checks cursor zero before the previous-rune marker in the modeled nonnegative cursor domain. A negative marker at a positive cursor gives the previous-operation error. Both failures preserve the state and create a fresh native error. With a nonnegative marker at a positive cursor, Go restores that exact index and clears the marker, including synthetic markers past the slice. Nil receiver dereference was checked in the native overlay. All 64,099 author inputs matched, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Bytes/ReaderRune.lean) · [Comparison](../../../difftest/bytes/reader-rune/check.go)

Nine independent paired native and Lean rows passed. Rewinding by a single byte failed two fixtures. Checking the marker before cursor zero failed two. The public control also confirmed that a successful rewind consumes the marker. Synthetic negative cursors and concurrency are outside scope. [Reviewer control](../../../difftest/bytes/reader-rune/review/main.go)
