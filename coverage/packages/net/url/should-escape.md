# Independent review of `net/url.shouldEscape`

**Verdict: FAITHFUL**

Evidence SHA-256: `a216a4b8142833fcd50d63fa2788a6f06c7c4899c6db3ff7f9cedb0394db62d1`

I checked the pinned Go 1.27.0 declaration and lookup table against the Lean model and exact native comparison. The byte table permits a byte when any bit in its class intersects the supplied mode. The zero mode permits no bytes. A mask containing only the hex bit permits F and f but not g.

I ran `go run ./difftest/net/url/escape/review` on darwin/arm64. Nine independent literal rows matched all six native declarations and compiled Lean. Four mutation witnesses distinguish overly broad hex classification, plus conversion for a combined mask, escaping the path-segment plus and lowercase hex output. The author comparison passed all 65,536 byte and mode combinations, 100 additional boundaries and 60,000 samples at seeds 17, 2718 and 65537. [Reviewer control](../../../../difftest/net/url/escape/review/main.go)

The claim concerns bytes and observable output. Allocation failure and output storage identity are outside scope.
