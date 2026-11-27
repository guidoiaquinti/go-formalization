# Independent review of `bytes.ToTitle`

**Verdict: FAITHFUL**

Evidence SHA-256: `3ef02d6e5078afe9ad14623a3bc895eb550b7ae7933bb222a02e5e54da2e2b2b`

Reviewed author commit `ab89f24` against pinned Go 1.27.0 `src/bytes/bytes.go:739 and 597–623`. Simple titlecase is applied to each decoded rune. A titlecase digraph can differ from both uppercase and lowercase, and malformed input uses the same one-byte advancement as the pinned decoder.

Eleven fixed byte fixtures and five Unicode mapping triplets passed against native Go and compiled Lean. The native observer checks nonnil results, exact capacity where modeled, and untouched source and spare cells after mutating every returned byte. Substituting uppercase for titlecase changed the two-digraph fixture. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/bytes/simple-case/review` from the repository root.

The author checked 6,506 boundaries and 60,000 input samples at seeds 17, 2718, and 65537, and compared all 1,114,112 code point lookups per mapping. Allocation failure and allocator-dependent growth capacity remain outside scope.
