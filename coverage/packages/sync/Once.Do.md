# Independent review of `sync.Once.Do`

**Verdict: FAITHFUL**

Evidence SHA-256: `93efa43b1325cbd897f4aa544dc74e6c4b535e6a4ef1f04fc8cf2c089502c4c0`

I reviewed evidence commit `8baf056` against Go 1.27.0 `src/sync/once.go` and `src/sync/oncefunc.go`. The completed fast path skips even nil callbacks. The first pending action records completion after either normal return or panic. A later call skips the new callback. The Lean theorem covers the completed fast path, and the native fixtures distinguish panic propagation on the first call from retry on a later call.

I ran `go run ./difftest/sync/once` on darwin/arm64 and reproduced 2,160 boundaries and 20,000 rows per seed for 17, 2718 and 65537. Six independent literal rows passed with `go run ./difftest/sync/once/review` and distinguished seven erroneous observations. The scope is serial and non-reentrant with panicnil=0. Concurrent publication and visibility, reentrancy, Goexit, and allocation behavior are outside scope. [Independent control](../../../difftest/sync/once/review/main.go)
