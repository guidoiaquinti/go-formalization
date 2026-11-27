# Independent review of `bytes.Reader.Reset`

**Verdict: FAITHFUL**

Evidence SHA-256: `8955c1b0acd2c7a54b5225c325d744a1d3df82ccf04e6033e5dd30fa813decc1`

The Go 1.27.0 source in `src/bytes/reader.go` overwrites the receiver with the supplied slice descriptor, cursor zero and previous-rune marker -1. The model preserves the scoped slice bytes, nilness, backing capacity and state observation.

The author comparison passed 3,170 boundaries and 60,000 sampled cases using seeds 17, 2718 and 65537. The shared reviewer control passed five literal Lean fixtures and native public state transitions on darwin/arm64. I reset after a successful rune read, verified the old unread state was cleared, and observed bytes from the new backing even after the old backing changed. I also checked the nil-receiver panic. Retaining the old state failed four Lean fixtures.

The model covers legal slice geometry and nonnegative int64 cursor states. Native execution on other architectures was outside this review.
