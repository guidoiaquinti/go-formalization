# Independent review of `net/url.Values.Del`

**Verdict: FAITHFUL**

Evidence SHA-256: `503bcaa0b1fc22a70babe94d9d1069b22f4bdcceac8de6f7753c2ce7a4cf6521`

Reviewed author commit `91920d8` against pinned Go 1.27.0 `src/net/url/url.go:919–921`. Removes the entry entirely, including when it stores an empty list. Missing keys, nil maps, and case-distinct neighbors are preserved.

Six independent literal fixtures passed against native Go and compiled Lean. Keeping the entry changed three fixtures. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/net/url/values/review` from the repository root.

The author replay passed 1,350 boundaries and 60,000 sampled cases at seeds 17, 2718, and 65537. Scope is sequential finite maps with unique raw keys and feasible allocations. Capacity, concurrent access, and arbitrary external alias-state modeling remain outside scope.
