# Independent review of `time.Unix`

**Verdict: FAITHFUL**

Evidence SHA-256: `2e0e5b432d750f9c7f8dfbaa8c98466d5d580a807fea53e6663c2d3e6316abc0`

I reviewed evidence commit `ff84bc7` against pinned Go 1.27.0 `src/time/time.go`. Pinned Go carries whole nanoseconds into the seconds and borrows for a negative remainder. Lean models the resulting Euclidean remainder and wrapped signed second word. The controls hit both carry and borrow, including overflow of seconds.

Pinned Go and Lean agreed on 2,560 fixed rows and 20,000 sampled rows per seed for 17, 2718 and 65537 on darwin/arm64 via `go run ./difftest/time/epoch-units`. Five independent native and Lean rows passed with `go run ./difftest/time/epoch-units/review` and detected eleven wrong outcomes across the group. The scoped kernel audit passed. Raw words outside public representable ranges describe the pinned implementation. Lean tracks time words and native comparison observes the exact Local location pointer. [Independent control](../../../difftest/time/epoch-units/review/main.go)
