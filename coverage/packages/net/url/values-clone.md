# Independent review of `net/url.Values.Clone`

**Verdict: FAITHFUL**

Evidence SHA-256: `dde97d98869af40dfaa3f93ea92a1a1f589110abbb45635e792510df46245df0`

Reviewed author commit `91920d8` against pinned Go 1.27.0 `src/net/url/url.go:930–939`. Preserves map and value-list nilness. A separate native test starts two keys on one backing slice, then verifies clone mutation affects neither the original nor the other cloned value list. It also checks two-way map and backing isolation.

Six independent literal fixtures passed against native Go and compiled Lean. Changing nil lists to empty lists changed three fixtures. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/net/url/values/review` from the repository root.

The author replay passed 1,350 boundaries and 60,000 sampled cases at seeds 17, 2718, and 65537. Scope is sequential finite maps with unique raw keys and feasible allocations. Capacity, concurrent access, and arbitrary external alias-state modeling remain outside scope.
