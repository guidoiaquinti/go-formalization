# Independent review of `index/suffixarray.Index.Lookup`

**Verdict: FAITHFUL**

Evidence SHA-256: `35675c3e58d005005ef7a1566ad3b7e00e32111fbfce8848752d5bea20bef302`

I checked the pinned Go 1.27.0 lookup code against the compiled Lean model and exact native overlay. The public method returns nil for empty needles, misses and a zero limit. A positive limit truncates to the first suffix-order matches, while a negative limit retains all. The output is an independent copy.

I ran `go run ./difftest/index/suffixarray/lookup/review` on darwin/arm64. Six independent literal rows matched pinned Go and compiled Lean. The controls distinguish re-sorting returned positions, clipping the private view capacity, ignoring a positive limit, and dereferencing a nil receiver before the zero-limit guard. The author replay passed 5,772 boundaries and 60,000 samples at seeds 17, 2718 and 65537. [Reviewer control](../../../../difftest/index/suffixarray/lookup/review/main.go)

The scope covers correctly formed suffix arrays in both internal offset widths, raw bytes and signed limits, plus the nil and zero index cases. Invalid manually corrupted tables are outside scope.
