# Independent review of `bytes.Title`

**Verdict: FAITHFUL**

Evidence SHA-256: `a49f5d8bda22b2ae4d57e90a0b75a95ebd6e3da0d2ce134ba865bc1e7f68cb57`

Reviewed author commit `5efb474` against pinned Go 1.27.0 `src/bytes/bytes.go:789–832`. The model starts after a space and carries the previous original rune. ASCII punctuation can mark a boundary. ASCII digits and underscores do not. The non-ASCII branch uses Unicode whitespace, so an em dash and a nonbreaking space have different effects. Invalid UTF-8 is re-encoded with one RuneError per invalid byte.

Twelve independent fixtures passed against native Go and compiled Lean. The reviewer checks literal results, initial capacity when the result fits, and source and spare backing after writes to the result. Resetting the prior rune to a space after every rune changed six fixtures. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/bytes/repeat-title/review` from the repository root.

Author evidence passed 638 boundaries and 60,000 sampled cases at seeds 17, 2718, and 65537. Allocator-dependent growth capacity and allocation failure remain outside scope.
