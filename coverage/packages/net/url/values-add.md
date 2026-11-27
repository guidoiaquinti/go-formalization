# Independent review of `net/url.Values.Add`

**Verdict: FAITHFUL**

Evidence SHA-256: `78fe234a1505053f71fba74f7d01e37e60b633ec7c63012e276a6f74ed74ab11`

Reviewed author commit `91920d8` against pinned Go 1.27.0 `src/net/url/url.go:914–916`. Appends a value at the end of the existing list and panics on a nil map. Case-distinct neighbors stay intact.

Six independent literal fixtures passed against native Go and compiled Lean. Prepending the value changed one fixture. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/net/url/values/review` from the repository root.

The author replay passed 1,350 boundaries and 60,000 sampled cases at seeds 17, 2718, and 65537. Scope is sequential finite maps with unique raw keys and feasible allocations. Capacity, concurrent access, and arbitrary external alias-state modeling remain outside scope.
