# Independent review of `net/url.URL.setFragment`

**Verdict: FAITHFUL**

Evidence SHA-256: `efe3ea0f31561804752d8b81ab9fe879a7346e8e2af2b8f94a402c2bbab8a3f9`

Reviewed author commit `3f4c447` against pinned Go 1.27.0 `src/net/url/url.go:727–739`. Decodes before touching the receiver and clears a canonical fragment spelling. Invalid escapes return their exact typed payload before dereferencing a nil receiver. The native check compares all unrelated fields.

Eight literal fixtures passed against the pinned native entries and compiled Lean. They include arbitrary raw octets, canonical and alternative spellings, malformed escapes and nil receivers. Always retaining the raw hint changed three fixtures across both setters. Checking nil before decoding changed one. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/net/url/path-fragment/review` from the repository root.

The author replay passed 66,336 boundaries and 60,000 sampled cases at seeds 17, 2718, and 65537. Concurrency, string storage layout, and allocation failure remain outside scope.
