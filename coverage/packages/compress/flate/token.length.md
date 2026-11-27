# Independent review of `compress/flate.token.length`

**Verdict: FAITHFUL**

Evidence SHA-256: `4934f4b557605a707ddec8a9c8e564c13c18c57dbc8d5acd801955e1e3567072`

I checked evidence commit `baa6682` against pinned Go 1.27.0 `src/compress/flate/token.go`. The source shifts the token right by 22 and converts to uint8, discarding the two type bits. The arithmetic decomposition reconstructs the input word from the type, length and offset fields. The independent all-ones fixture detects omitting the highest length bit.

The pinned same-package Go comparison and Lean replay at `go run ./difftest/compress/flate/token-fields` passed 7,168 fixed words and 20,000 sampled words each for seeds 17, 2718 and 65537 on darwin/arm64. Four independent fixtures passed with `go run ./difftest/compress/flate/token-fields/review`. Scope covers every packed uint32 word. Token construction and semantic validity are outside scope. [Independent control](../../../../difftest/compress/flate/token-fields/review/main.go)
