# Independent review of `net/url.Values.Get`

**Verdict: FAITHFUL**

Evidence SHA-256: `6ef16e276f5e8eebd3c4d21e1fe3f446498c5e6cadb4d20036003c9a150bd2fb`

Reviewed author commit `91920d8` against pinned Go 1.27.0 `src/net/url/url.go:898–904`. Reads only the first value at the exact key. Absent keys and nil or empty value lists return the empty string.

Six independent literal fixtures passed against native Go and compiled Lean. Returning the last value changed one fixture. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/net/url/values/review` from the repository root.

The author replay passed 1,350 boundaries and 60,000 sampled cases at seeds 17, 2718, and 65537. Scope is sequential finite maps with unique raw keys and feasible allocations. Capacity, concurrent access, and arbitrary external alias-state modeling remain outside scope.
