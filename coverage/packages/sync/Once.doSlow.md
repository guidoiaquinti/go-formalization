# Independent review of `sync.Once.doSlow`

**Verdict: FAITHFUL**

Evidence SHA-256: `dbba0f6d5f065ce2179ce46e022cae86bf8f209af93e02632407fd67f69902dd`

I reviewed evidence commit `8baf056` against Go 1.27.0 `src/sync/once.go` and `src/sync/oncefunc.go`. The source locks before checking completion and uses deferred calls to store completion before unlocking. The Lean model covers the serial before-and-after state, callback entry and panic. The same-package Go bridge observes the lock held and completion unset inside an entered callback and verifies that the lock is released afterward. The already-completed native slow-path fixture checks that the callback is skipped. Lock acquisition in that skipped branch follows from source inspection.

I ran `go run ./difftest/sync/once` on darwin/arm64 and reproduced 2,160 boundaries and 20,000 rows per seed for 17, 2718 and 65537. Six independent literal rows passed with `go run ./difftest/sync/once/review` and distinguished seven erroneous observations. The scope is serial and non-reentrant with panicnil=0. Concurrent publication and visibility, reentrancy, Goexit, and allocation behavior are outside scope. [Independent control](../../../difftest/sync/once/review/main.go)
