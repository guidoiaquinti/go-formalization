# Independent review of `net/url.validEncoded`

**Verdict: FAITHFUL**

Evidence SHA-256: `07e9aa7672631b8d988fe24c7476bd4a31db576aa8decbfd953e5b61dcc2a2f4`

Reviewed author commit `3f4c447` against pinned Go 1.27.0 `src/net/url/url.go:703–724`. Checks the raw byte allowlist, including percent and square brackets. It accepts a lone percent or malformed percent escape as a raw byte sequence. The accessors separately attempt decoding.

Eight literal fixtures passed against the pinned native entries and compiled Lean. They include arbitrary raw octets, canonical and alternative spellings, malformed escapes and nil receivers. Removing the percent exception changed four fixtures. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/net/url/path-fragment/review` from the repository root.

The author replay passed 66,336 boundaries and 60,000 sampled cases at seeds 17, 2718, and 65537. Concurrency, string storage layout, and allocation failure remain outside scope.
