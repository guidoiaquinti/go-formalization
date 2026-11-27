# Independent review of `index/suffixarray.freq_32`

**Verdict: FAITHFUL**

Evidence SHA-256: `5bb5b78ee3cde2406c91f068ab4b04e986bc2e1afbd75fd17c13357a304198dd`

I checked the pinned Go 1.27.0 declaration and the compiled Lean observation against the native private call. The helper consults the first frequency slot for cache validity. Otherwise it clears and counts into the supplied frequency slice, or aliases the bucket when frequency is nil. Counts before a bounds failure survive. Byte variants extend the slice to the fixed 256-symbol alphabet. Integer variants use the retained view length.

I ran `go run ./difftest/index/suffixarray/buckets/review` on darwin/arm64. Eight independent literal rows across all four profiles passed. The rows distinguish using a cached count table, separate alias and nonalias behavior, preserving completed counts on a bounds panic, the inclusive and exclusive prefix totals, and modular accumulation for both widths. The author comparison covers 2,908 fixed rows and 60,000 seeded rows per declaration at seeds 17, 2718 and 65537. [Reviewer control](../../../../difftest/index/suffixarray/buckets/review/main.go)

Scope covers sequential access and disjoint frequency and bucket allocations, or a nil frequency input. The model observes returned views and retained backing. Arbitrarily overlapping nonnil buffers are outside scope.
