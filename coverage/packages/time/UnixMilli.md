# Independent review of `time.UnixMilli`

**Verdict: FAITHFUL**

Evidence SHA-256: `fb929e4c60175e706e23fad11d1faff6cf72c1418af05f4a331ff2f90ce56db0`

I reviewed evidence commit `ff84bc7` against pinned Go 1.27.0 `src/time/time.go`. Pinned Go delegates after dividing signed milliseconds into seconds and fractional milliseconds. The constructors preserve a nonnegative fraction for negative input. The controls cover -1 millisecond and the minimum signed unit count.

Pinned Go and Lean agreed on 2,560 fixed rows and 20,000 sampled rows per seed for 17, 2718 and 65537 on darwin/arm64 via `go run ./difftest/time/epoch-units`. Five independent native and Lean rows passed with `go run ./difftest/time/epoch-units/review` and detected eleven wrong outcomes across the group. The scoped kernel audit passed. Raw words outside public representable ranges describe the pinned implementation. Lean tracks time words and native comparison observes the exact Local location pointer. [Independent control](../../../difftest/time/epoch-units/review/main.go)
