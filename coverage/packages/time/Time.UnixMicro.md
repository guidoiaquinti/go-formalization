# Independent review of `time.Time.UnixMicro`

**Verdict: FAITHFUL**

Evidence SHA-256: `b043d4e443c93a8b8b757433f2ad46b96c33bfd6ba0da2d322fc5a914ae13c33`

I reviewed evidence commit `ff84bc7` against pinned Go 1.27.0 `src/time/time.go`. Pinned Go decodes the second and nanosecond words, multiplies seconds by 1,000,000 and divides the fractional contribution by 1,000. Fixed pre-epoch and high-second cases distinguish rounding and signed word overflow.

Pinned Go and Lean agreed on 2,560 fixed rows and 20,000 sampled rows per seed for 17, 2718 and 65537 on darwin/arm64 via `go run ./difftest/time/epoch-units`. Five independent native and Lean rows passed with `go run ./difftest/time/epoch-units/review` and detected eleven wrong outcomes across the group. The scoped kernel audit passed. Raw words outside public representable ranges describe the pinned implementation. Lean tracks time words and native comparison observes the exact Local location pointer. [Independent control](../../../difftest/time/epoch-units/review/main.go)
