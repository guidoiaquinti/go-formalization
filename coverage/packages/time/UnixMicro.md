# Independent review of `time.UnixMicro`

**Verdict: FAITHFUL**

Evidence SHA-256: `db1e10c7f1bf22e03e4f19fe33e708d80f8392ee41315a60e4324216103e9c77`

I reviewed evidence commit `ff84bc7` against pinned Go 1.27.0 `src/time/time.go`. Pinned Go delegates after dividing signed microseconds into seconds and fractional microseconds. The controls cover negative units across zero and the signed endpoints. The native overlay checks raw Local pointer identity.

Pinned Go and Lean agreed on 2,560 fixed rows and 20,000 sampled rows per seed for 17, 2718 and 65537 on darwin/arm64 via `go run ./difftest/time/epoch-units`. Five independent native and Lean rows passed with `go run ./difftest/time/epoch-units/review` and detected eleven wrong outcomes across the group. The scoped kernel audit passed. Raw words outside public representable ranges describe the pinned implementation. Lean tracks time words and native comparison observes the exact Local location pointer. [Independent control](../../../difftest/time/epoch-units/review/main.go)
