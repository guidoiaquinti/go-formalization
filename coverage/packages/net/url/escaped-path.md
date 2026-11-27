# Independent review of `net/url.URL.EscapedPath`

**Verdict: FAITHFUL**

Evidence SHA-256: `9cc6290fa02a9640b4c4796fcd384d3b4c6a81c7ea7c399e01ff40ef096ee4f1`

Reviewed author commit `3f4c447` against pinned Go 1.27.0 `src/net/url/url.go:687–699`. An accepted raw escape that decodes to the path takes precedence over the asterisk exception. A raw spelling may be disallowed even when the decoded bytes match. The observer checks every URL field after an accessor.

Eight literal fixtures passed against the pinned native entries and compiled Lean. They include arbitrary raw octets, canonical and alternative spellings, malformed escapes and nil receivers. Handling asterisk before the accepted raw hint changed one fixture. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/net/url/path-fragment/review` from the repository root.

The author replay passed 66,336 boundaries and 60,000 sampled cases at seeds 17, 2718, and 65537. Concurrency, string storage layout, and allocation failure remain outside scope.
