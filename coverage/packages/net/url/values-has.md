# Independent review of `net/url.Values.Has`

**Verdict: FAITHFUL**

Evidence SHA-256: `20adf3a0cb676c2e0505fde13bf75dbd15ac707beaba40e21e5506736d60e7fb`

Reviewed author commit `91920d8` against pinned Go 1.27.0 `src/net/url/url.go:924–927`. Presence is independent of the value list. The case-distinct keys and both nil and nonnil empty lists have separate observations.

Six independent literal fixtures passed against native Go and compiled Lean. Requiring a nonempty list changed two fixtures. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/net/url/values/review` from the repository root.

The author replay passed 1,350 boundaries and 60,000 sampled cases at seeds 17, 2718, and 65537. Scope is sequential finite maps with unique raw keys and feasible allocations. Capacity, concurrent access, and arbitrary external alias-state modeling remain outside scope.
