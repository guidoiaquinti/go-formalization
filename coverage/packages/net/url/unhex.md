# Independent review of `net/url.unhex`

**Verdict: FAITHFUL**

Evidence SHA-256: `04f1eebba5bd95e4c7813f096fc085b11a27a35a36883c4ac23f0ab0d126dd40`

I checked the pinned Go 1.27.0 declaration and lookup table against the Lean model and exact native comparison. On valid hex digits, the arithmetic yields the exact nibble. Capital and lowercase F produce fifteen. Native outputs for non-hex bytes are observed as an implementation detail.

I ran `go run ./difftest/net/url/escape/review` on darwin/arm64. Nine independent literal rows matched all six native declarations and compiled Lean. Four mutation witnesses distinguish overly broad hex classification, plus conversion for a combined mask, escaping the path-segment plus and lowercase hex output. The author comparison passed all 65,536 byte and mode combinations, 100 additional boundaries and 60,000 samples at seeds 17, 2718 and 65537. [Reviewer control](../../../../difftest/net/url/escape/review/main.go)

The claim concerns bytes and observable output. Allocation failure and output storage identity are outside scope.
