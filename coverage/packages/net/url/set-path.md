# Independent review of `net/url.URL.setPath`

**Verdict: FAITHFUL**

Evidence SHA-256: `379969d63941462e8adf1802df1e2019ebc0b5a1a9fc1580162d1eb5cbd0e83c`

Reviewed author commit `3f4c447` against pinned Go 1.27.0 `src/net/url/url.go:660–672`. Decodes before touching the receiver, updates only the path fields, and retains a noncanonical or differently cased spelling. Invalid input preserves the prior state and returns an exact typed escape payload.

Eight literal fixtures passed against the pinned native entries and compiled Lean. They include arbitrary raw octets, canonical and alternative spellings, malformed escapes and nil receivers. Always retaining the raw hint changed three fixtures across both setters. Checking nil before decoding changed one. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/net/url/path-fragment/review` from the repository root.

The author replay passed 66,336 boundaries and 60,000 sampled cases at seeds 17, 2718, and 65537. Concurrency, string storage layout, and allocation failure remain outside scope.
