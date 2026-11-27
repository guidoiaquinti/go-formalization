# Independent review of `time.Time.UnixNano`

**Verdict: FAITHFUL**

Evidence SHA-256: `77cb1bc06bb9a2d50aeb89200c35738f1fb5d163499cb6c55c246d1b3395b7e5`

I reviewed evidence commit `ff84bc7` against pinned Go 1.27.0 `src/time/time.go`. Pinned Go combines decoded seconds and nanoseconds as wrapped int64 arithmetic. The fixtures cover an overflowing nanosecond result and a monotonic-bearing wall word whose seconds come from packed wall storage.

Pinned Go and Lean agreed on 2,560 fixed rows and 20,000 sampled rows per seed for 17, 2718 and 65537 on darwin/arm64 via `go run ./difftest/time/epoch-units`. Five independent native and Lean rows passed with `go run ./difftest/time/epoch-units/review` and detected eleven wrong outcomes across the group. The scoped kernel audit passed. Raw words outside public representable ranges describe the pinned implementation. Lean tracks time words and native comparison observes the exact Local location pointer. [Independent control](../../../difftest/time/epoch-units/review/main.go)
