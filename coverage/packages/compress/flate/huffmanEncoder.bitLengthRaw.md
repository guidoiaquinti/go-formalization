# Independent review of `compress/flate.huffmanEncoder.bitLengthRaw`

**Verdict: FAITHFUL**

Evidence SHA-256: `f710cccfe6ece464e729c5a26d0cc265e9785df21166893ccdf5aef1b2258592`

Every input byte indexes the table and is charged at least one bit, including a packed zero or a nonzero word with zero length. A short table panics even when the corresponding frequency list has only zeros. The author replay matched 60,014 observations per width at seeds 17, 2718 and 65537. Native execution matched the 64-bit observations and Lean checked both 32-bit and 64-bit arithmetic. [Model](../../../lean/GoLean/Compress/Flate/EncoderCost.lean) · [Comparison](../../../difftest/compress/flate/encoder-cost/check.go)

The independent command `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/compress/flate/encoder-cost/review` passed seven Lean and six native fixtures and rejected seven altered outcomes. The controls include a modular overflow witness at width 32 and preserve the native 64-bit result. Native 32-bit execution remains outside scope. [Reviewer control](../../../difftest/compress/flate/encoder-cost/review/main.go)
