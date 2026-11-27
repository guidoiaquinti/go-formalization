# Independent review of `sync.OnceFunc`

**Verdict: FAITHFUL**

Evidence SHA-256: `27e749061ad60f8c621f0f6dc7781da2f13da7ab70155f721aecff7666880322`

I reviewed evidence commit `8baf056` against Go 1.27.0 `src/sync/once.go` and `src/sync/oncefunc.go`. Creation defers callback execution. In the first call Go recovers a panic, saves its value and immediately rethrows it. Later calls rethrow the saved value without invoking the callback. The pinned overlay checks pointer identity for explicit payloads and the runtime panic category for nil callbacks and panic(nil).

I ran `go run ./difftest/sync/once` on darwin/arm64 and reproduced 2,160 boundaries and 20,000 rows per seed for 17, 2718 and 65537. Six independent literal rows passed with `go run ./difftest/sync/once/review` and distinguished seven erroneous observations. The scope is serial and non-reentrant with panicnil=0. Concurrent publication and visibility, reentrancy, Goexit, and allocation behavior are outside scope. [Independent control](../../../difftest/sync/once/review/main.go)
