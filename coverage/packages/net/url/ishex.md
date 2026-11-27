# Independent review of `net/url.ishex`

**Verdict: FAITHFUL**

Evidence SHA-256: `e2ed3eeb710daee6e0bc3579ca5d0fab85a9f8d95467d3c4b819b7afc9ee63c7`

I checked the pinned Go 1.27.0 declaration and lookup table against the Lean model and exact native comparison. The hex flag classifies ASCII decimal digits and the first six upper- and lowercase letters. The fixture rejects treating g as hex.

I ran `go run ./difftest/net/url/escape/review` on darwin/arm64. Nine independent literal rows matched all six native declarations and compiled Lean. Four mutation witnesses distinguish overly broad hex classification, plus conversion for a combined mask, escaping the path-segment plus and lowercase hex output. The author comparison passed all 65,536 byte and mode combinations, 100 additional boundaries and 60,000 samples at seeds 17, 2718 and 65537. [Reviewer control](../../../../difftest/net/url/escape/review/main.go)

The claim concerns bytes and observable output. Allocation failure and output storage identity are outside scope.
