# Independent review of `time.Time.UnixMilli`

**Verdict: FAITHFUL**

Evidence SHA-256: `148bcd2762715c5e49281dfc79601324332da5930c584f53e56ea54f0a05317a`

I reviewed evidence commit `ff84bc7` against pinned Go 1.27.0 `src/time/time.go`. Pinned Go decodes the second and nanosecond words, multiplies seconds by 1,000 and truncates the fractional contribution. The pre-epoch fixture distinguishes truncating the nonnegative fraction from truncating the combined negative instant. A raw high-second fixture checks signed wrap.

Pinned Go and Lean agreed on 2,560 fixed rows and 20,000 sampled rows per seed for 17, 2718 and 65537 on darwin/arm64 via `go run ./difftest/time/epoch-units`. Five independent native and Lean rows passed with `go run ./difftest/time/epoch-units/review` and detected eleven wrong outcomes across the group. The scoped kernel audit passed. Raw words outside public representable ranges describe the pinned implementation. Lean tracks time words and native comparison observes the exact Local location pointer. [Independent control](../../../difftest/time/epoch-units/review/main.go)
