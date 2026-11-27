# Independent review of `index/suffixarray.freq_8_32`

**Verdict: FAITHFUL**

Evidence SHA-256: `bd29071742fc568612300a0feadedcd9f77afb495b8f8af9a23e3d7c23d075c8`

I checked the pinned Go 1.27.0 declaration and the compiled Lean observation against the native private call. The helper consults the first frequency slot for cache validity. Otherwise it clears and counts into the supplied frequency slice, or aliases the bucket when frequency is nil. Counts before a bounds failure survive. Byte variants extend the slice to the fixed 256-symbol alphabet. Integer variants use the retained view length.

I ran `go run ./difftest/index/suffixarray/buckets/review` on darwin/arm64. Eight independent literal rows across all four profiles passed. The rows distinguish using a cached count table, separate alias and nonalias behavior, preserving completed counts on a bounds panic, the inclusive and exclusive prefix totals, and modular accumulation for both widths. The author comparison covers 2,908 fixed rows and 60,000 seeded rows per declaration at seeds 17, 2718 and 65537. [Reviewer control](../../../../difftest/index/suffixarray/buckets/review/main.go)

Scope covers sequential access and disjoint frequency and bucket allocations, or a nil frequency input. The model observes returned views and retained backing. Arbitrarily overlapping nonnil buffers are outside scope.
