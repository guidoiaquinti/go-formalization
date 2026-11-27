# Independent review of `sync.OnceValues`

**Verdict: FAITHFUL**

Evidence SHA-256: `19c390a29d5e63e19627f00ec54d4a4aa95a8b94c2dfcfa41935f82575c64d91`

I reviewed evidence commit `8baf056` against Go 1.27.0 `src/sync/once.go` and `src/sync/oncefunc.go`. Creation defers callback execution. One normal call caches both values together and a panic keeps the saved failure for subsequent calls. The Lean theorem proves the pair is returned again without an invocation. The independent fixtures distinguish retaining the full pair from losing its second value.

I ran `go run ./difftest/sync/once` on darwin/arm64 and reproduced 2,160 boundaries and 20,000 rows per seed for 17, 2718 and 65537. Six independent literal rows passed with `go run ./difftest/sync/once/review` and distinguished seven erroneous observations. The scope is serial and non-reentrant with panicnil=0. Concurrent publication and visibility, reentrancy, Goexit, and allocation behavior are outside scope. [Independent control](../../../difftest/sync/once/review/main.go)
