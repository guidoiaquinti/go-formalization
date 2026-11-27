# Independent review of `net/url.URL.EscapedFragment`

**Verdict: FAITHFUL**

Evidence SHA-256: `24f5aed7053dc6fd4c136fadac5a63737ce594071b0b9ede4ef5c81ad3db6fac`

Reviewed author commit `3f4c447` against pinned Go 1.27.0 `src/net/url/url.go:750–759`. The raw hint must survive both byte validation and decoded-byte equality. The fragment rules permit the literal question mark in the chosen fixture. Malformed escapes and mismatched hints fall back.

Eight literal fixtures passed against the pinned native entries and compiled Lean. They include arbitrary raw octets, canonical and alternative spellings, malformed escapes and nil receivers. Accepting successfully unescaped hints without comparing decoded values changed one fixture. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/net/url/path-fragment/review` from the repository root.

The author replay passed 66,336 boundaries and 60,000 sampled cases at seeds 17, 2718, and 65537. Concurrency, string storage layout, and allocation failure remain outside scope.
