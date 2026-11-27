# Independent review of `time.unixTime`

**Verdict: FAITHFUL**

Evidence SHA-256: `ac69cdb611d8cfae5095d07008564f755b427a0baf475092fcdd845c6e9034a3`

I reviewed evidence commit `ff84bc7` against pinned Go 1.27.0 `src/time/time.go`. The private constructor converts signed int32 nanoseconds to uint64, adds the Unix-to-internal second offset modulo 64 bits, and stores the native Local pointer. The control distinguishes sign extension from zero extension and tests extreme seconds.

Pinned Go and Lean agreed on 2,560 fixed rows and 20,000 sampled rows per seed for 17, 2718 and 65537 on darwin/arm64 via `go run ./difftest/time/epoch-units`. Five independent native and Lean rows passed with `go run ./difftest/time/epoch-units/review` and detected eleven wrong outcomes across the group. The scoped kernel audit passed. Raw words outside public representable ranges describe the pinned implementation. Lean tracks time words and native comparison observes the exact Local location pointer. [Independent control](../../../difftest/time/epoch-units/review/main.go)
