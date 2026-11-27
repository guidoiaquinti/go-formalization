# Independent review of `net/url.PathEscape`

**Verdict: FAITHFUL**

Evidence SHA-256: `3a2178981cc821b53435d8d5947e03db14b9045761d95d4f2b821121b65bafd2`

I checked the pinned Go 1.27.0 declaration and lookup table against the Lean model and exact native comparison. The path segment wrapper preserves an input plus, escapes a slash and percent-encodes spaces. It accepts raw strings including malformed UTF-8.

I ran `go run ./difftest/net/url/escape/review` on darwin/arm64. Nine independent literal rows matched all six native declarations and compiled Lean. Four mutation witnesses distinguish overly broad hex classification, plus conversion for a combined mask, escaping the path-segment plus and lowercase hex output. The author comparison passed all 65,536 byte and mode combinations, 100 additional boundaries and 60,000 samples at seeds 17, 2718 and 65537. [Reviewer control](../../../../difftest/net/url/escape/review/main.go)

The claim concerns bytes and observable output. Allocation failure and output storage identity are outside scope.
