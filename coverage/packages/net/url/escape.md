# Independent review of `net/url.escape`

**Verdict: FAITHFUL**

Evidence SHA-256: `03c94b6e6c00e42257ac3899f24d0a5d95aaf4c56ac49914fb5f4950e41cdd61`

I checked the pinned Go 1.27.0 declaration and lookup table against the Lean model and exact native comparison. The helper converts spaces to plus only in the exact query component mode. A combined mask can instead percent-encode the space. Hex digits in the output are uppercase, even on malformed UTF-8 bytes.

I ran `go run ./difftest/net/url/escape/review` on darwin/arm64. Nine independent literal rows matched all six native declarations and compiled Lean. Four mutation witnesses distinguish overly broad hex classification, plus conversion for a combined mask, escaping the path-segment plus and lowercase hex output. The author comparison passed all 65,536 byte and mode combinations, 100 additional boundaries and 60,000 samples at seeds 17, 2718 and 65537. [Reviewer control](../../../../difftest/net/url/escape/review/main.go)

The claim concerns bytes and observable output. Allocation failure and output storage identity are outside scope.
