# Independent review of `compress/flate.maxNode`

**Verdict: FAITHFUL**

Evidence SHA-256: `835872def938ffb2723e55cb9966514e25ff3effb8a6b5b53f124a8f47c2b24d`

The sentinel contains the maximum value in both uint16 fields. Native observations confirmed both fields and the control rejects changing one to zero. The pinned Go and Lean comparison passed 64,104 author rows at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Compress/Flate/EncoderInit.lean) · [Comparison](../../../difftest/compress/flate/encoder-init/check.go)

The independent command `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/compress/flate/encoder-init/review` passed ten Lean and six native fixtures and rejected ten altered results across the group. Native execution used darwin/arm64. [Reviewer control](../../../difftest/compress/flate/encoder-init/review/main.go)
