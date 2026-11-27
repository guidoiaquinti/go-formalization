# Independent review of `net/url.Values.Set`

**Verdict: FAITHFUL**

Evidence SHA-256: `01dccd4a18aced5d95d1b2b0d9529a649c7822628b4baa92df28a45cf56ec1f2`

Reviewed author commit `91920d8` against pinned Go 1.27.0 `src/net/url/url.go:908–910`. Replaces the old list with a singleton, including when the written value is the empty string. A nil map raises the expected assignment panic.

Six independent literal fixtures passed against native Go and compiled Lean. Appending to the old list changed one fixture. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/net/url/values/review` from the repository root.

The author replay passed 1,350 boundaries and 60,000 sampled cases at seeds 17, 2718, and 65537. Scope is sequential finite maps with unique raw keys and feasible allocations. Capacity, concurrent access, and arbitrary external alias-state modeling remain outside scope.
