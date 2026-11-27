# Independent review of `bytes.NewReader`

**Verdict: FAITHFUL**

Evidence SHA-256: `d6ef7842b974acc844120ca38c0fe677db83d8edc6cba64fd929fc7ace09a17a`

The Go 1.27.0 source in `src/bytes/reader.go` copies the input slice descriptor into a new Reader with cursor zero and previous-rune marker -1. The model preserves the scoped slice bytes, nilness, backing capacity and state observation.

The author comparison passed 3,170 boundaries and 60,000 sampled cases using seeds 17, 2718 and 65537. The shared reviewer control passed five literal Lean fixtures and native public state transitions on darwin/arm64. I changed an input byte after construction and observed it through an actual rune read. Nil and empty nonnil inputs retain their respective slice shapes in the overlay. Initializing the marker to zero failed five Lean fixtures.

The model covers legal slice geometry and nonnegative int64 cursor states. Native execution on other architectures was outside this review.
