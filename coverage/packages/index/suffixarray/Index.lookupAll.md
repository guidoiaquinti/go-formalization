# Independent review of `index/suffixarray.Index.lookupAll`

**Verdict: FAITHFUL**

Evidence SHA-256: `0988282f06b7c37fa1974f71267f4dc6f27ac07d12cba100cd31c3b1d1ae0409`

I checked the pinned Go 1.27.0 lookup code against the compiled Lean model and exact native overlay. The private method finds the contiguous prefix interval with two binary searches. Its returned view points into the original suffix array and retains spare capacity after the match interval.

I ran `go run ./difftest/index/suffixarray/lookup/review` on darwin/arm64. Six independent literal rows matched pinned Go and compiled Lean. The controls distinguish re-sorting returned positions, clipping the private view capacity, ignoring a positive limit, and dereferencing a nil receiver before the zero-limit guard. The author replay passed 5,772 boundaries and 60,000 samples at seeds 17, 2718 and 65537. [Reviewer control](../../../../difftest/index/suffixarray/lookup/review/main.go)

The scope covers correctly formed suffix arrays in both internal offset widths, raw bytes and signed limits, plus the nil and zero index cases. Invalid manually corrupted tables are outside scope.
