# Independent review of `net/url.Values.Encode`

**Verdict: FAITHFUL**

Evidence SHA-256: `157448c23e6228f00e16079198b58c59ba1bfeeb66b4708a38575490bb9cc55c`

Reviewed author commit `91920d8` against pinned Go 1.27.0 `src/net/url/url.go:1015–1043`. Orders raw byte keys before escaping, preserves list order, and omits absent items. The fixtures distinguish literal plus from space, include an empty key and value, and require high-byte and NUL escapes.

Six independent literal fixtures passed against native Go and compiled Lean. Emitting the input entry order changed two fixtures. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/net/url/values/review` from the repository root.

The author replay passed 1,350 boundaries and 60,000 sampled cases at seeds 17, 2718, and 65537. Scope is sequential finite maps with unique raw keys and feasible allocations. Capacity, concurrent access, and arbitrary external alias-state modeling remain outside scope.
