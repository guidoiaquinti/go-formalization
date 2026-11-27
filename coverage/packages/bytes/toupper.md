# Independent review of `bytes.ToUpper`

**Verdict: FAITHFUL**

Evidence SHA-256: `4e151eae624b7ecb47d8732384807cfe6703acf3d73fb29187f9f96a13de2861`

Reviewed author commit `ab89f24` against pinned Go 1.27.0 `src/bytes/bytes.go:679–706`. The ASCII fast path copies and clips capacity to length if it changes a byte. Unchanged ASCII is copied with allocator-dependent capacity. Non-ASCII uses the simple uppercase mapping and decodes each malformed byte as RuneError.

Eleven fixed byte fixtures and five Unicode mapping triplets passed against native Go and compiled Lean. The native observer checks nonnil results, exact capacity where modeled, and untouched source and spare cells after mutating every returned byte. A mutant enlarging changed ASCII capacity changed three rows across the two ASCII paths. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/bytes/simple-case/review` from the repository root.

The author checked 6,506 boundaries and 60,000 input samples at seeds 17, 2718, and 65537, and compared all 1,114,112 code point lookups per mapping. Allocation failure and allocator-dependent growth capacity remain outside scope.
