# Independent review of `compress/flate.token.literal`

**Verdict: FAITHFUL**

Evidence SHA-256: `2fe49363025fe326a8b981ebf78484c55ef31f430755a780f074f8f065b9657a`

I checked evidence commit `baa6682` against pinned Go 1.27.0 `src/compress/flate/token.go`. The source converts the token to uint8, retaining the low byte even for a non-literal type. The theorem identifies that byte within the offset field. The control rejects using the entire low 16-bit word.

The pinned same-package Go comparison and Lean replay at `go run ./difftest/compress/flate/token-fields` passed 7,168 fixed words and 20,000 sampled words each for seeds 17, 2718 and 65537 on darwin/arm64. Four independent fixtures passed with `go run ./difftest/compress/flate/token-fields/review`. Scope covers every packed uint32 word. Token construction and semantic validity are outside scope. [Independent control](../../../../difftest/compress/flate/token-fields/review/main.go)
