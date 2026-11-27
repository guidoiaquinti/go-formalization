# Independent review of `bytes.ToLower`

**Verdict: FAITHFUL**

Evidence SHA-256: `44904e163472c432d9e6f95fecd5c512d4b884c2481353fc590bd48f59fab4b3`

Reviewed author commit `ab89f24` against pinned Go 1.27.0 `src/bytes/bytes.go:710–735`. Changed ASCII is copied into a clipped buffer and unchanged ASCII is copied via append. The Unicode path preserves the input-sized allocation when the result fits. The fixture for Kelvin sign and dotted capital I shrinks five input bytes to two output bytes and retains capacity five.

Eleven fixed byte fixtures and five Unicode mapping triplets passed against native Go and compiled Lean. The native observer checks nonnil results, exact capacity where modeled, and untouched source and spare cells after mutating every returned byte. A mutant allocating by output length changed that fixture, and a mutant enlarging the changed ASCII capacity changed three rows across the two ASCII paths. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/bytes/simple-case/review` from the repository root.

The author checked 6,506 boundaries and 60,000 input samples at seeds 17, 2718, and 65537, and compared all 1,114,112 code point lookups per mapping. Allocation failure and allocator-dependent growth capacity remain outside scope.
