# Independent review of `sync.OnceValue`

**Verdict: FAITHFUL**

Evidence SHA-256: `a840c66fad6b8dc065a1ba559a5794dd3ff8d7d2694ccb61eb3ff2ddfd3f4f89`

I reviewed evidence commit `8baf056` against Go 1.27.0 `src/sync/once.go` and `src/sync/oncefunc.go`. Creation defers callback execution and the return value is cached after successful completion. The Lean theorem proves reuse of the cached value and the native fixtures cover both saved values and repeated pointer panic payloads. The result and callback count distinguish reuse from recomputation.

I ran `go run ./difftest/sync/once` on darwin/arm64 and reproduced 2,160 boundaries and 20,000 rows per seed for 17, 2718 and 65537. Six independent literal rows passed with `go run ./difftest/sync/once/review` and distinguished seven erroneous observations. The scope is serial and non-reentrant with panicnil=0. Concurrent publication and visibility, reentrancy, Goexit, and allocation behavior are outside scope. [Independent control](../../../difftest/sync/once/review/main.go)
