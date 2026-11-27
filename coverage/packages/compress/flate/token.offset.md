# Independent review of `compress/flate.token.offset`

**Verdict: FAITHFUL**

Evidence SHA-256: `20a822f0d5aadffe48e3db548ec64fa264c1e9d75788b6662a93521b25fcb831`

I checked evidence commit `baa6682` against pinned Go 1.27.0 `src/compress/flate/token.go`. The source masks with 2^22-1. The remainder representation preserves the whole low field, including the embedded offset code. The control rejects trimming it to the low 16 bits.

The pinned same-package Go comparison and Lean replay at `go run ./difftest/compress/flate/token-fields` passed 7,168 fixed words and 20,000 sampled words each for seeds 17, 2718 and 65537 on darwin/arm64. Four independent fixtures passed with `go run ./difftest/compress/flate/token-fields/review`. Scope covers every packed uint32 word. Token construction and semantic validity are outside scope. [Independent control](../../../../difftest/compress/flate/token-fields/review/main.go)
