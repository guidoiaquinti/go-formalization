# Independent review of `net/url.QueryEscape`

**Verdict: FAITHFUL**

Evidence SHA-256: `56cb4072bd27757636bb9ea33705a7cc46741406cd980a4f8774e2c1b3cba0e8`

I checked the pinned Go 1.27.0 declaration and lookup table against the Lean model and exact native comparison. The query wrapper converts spaces to plus and percent-encodes an input plus. It percent-encodes raw non-ASCII bytes independently.

I ran `go run ./difftest/net/url/escape/review` on darwin/arm64. Nine independent literal rows matched all six native declarations and compiled Lean. Four mutation witnesses distinguish overly broad hex classification, plus conversion for a combined mask, escaping the path-segment plus and lowercase hex output. The author comparison passed all 65,536 byte and mode combinations, 100 additional boundaries and 60,000 samples at seeds 17, 2718 and 65537. [Reviewer control](../../../../difftest/net/url/escape/review/main.go)

The claim concerns bytes and observable output. Allocation failure and output storage identity are outside scope.
