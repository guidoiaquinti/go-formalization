# Independent review of `bytes.Reader.Size`

**Verdict: FAITHFUL**

Evidence SHA-256: `1d1b1d3980914c656c1230d716366810ff361ea92ee624c6c39ef6ac2babb554`

The Go 1.27.0 source in `src/bytes/reader.go` returns int64(len(r.s)) without changing the cursor or previous-rune marker. The model preserves the scoped slice bytes, nilness, backing capacity and state observation.

The author comparison passed 3,170 boundaries and 60,000 sampled cases using seeds 17, 2718 and 65537. The shared reviewer control passed five literal Lean fixtures and native public state transitions on darwin/arm64. I checked the total size after reading a rune, verified that UnreadRune still succeeds, checked the total at the maximum signed cursor, and checked the nil-receiver panic. Reporting unread length failed two Lean fixtures.

The model covers legal slice geometry and nonnegative int64 cursor states. Native execution on other architectures was outside this review.
