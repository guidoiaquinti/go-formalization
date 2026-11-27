# Independent review of `compress/flate.huffmanEncoder.bitLength`

**Verdict: FAITHFUL**

Evidence SHA-256: `8d705a7a6097eeeddeaa18a423b693a9afee6ae1f08ffdfbffe5db02550912d1`

Zero frequencies skip table lookup even with a nil receiver or an index beyond the table. Nonzero frequencies multiply the low-byte length by the frequency, with native-width accumulation. The author replay matched 60,014 observations per width at seeds 17, 2718 and 65537. Native execution matched the 64-bit observations and Lean checked both 32-bit and 64-bit arithmetic. [Model](../../../lean/GoLean/Compress/Flate/EncoderCost.lean) · [Comparison](../../../difftest/compress/flate/encoder-cost/check.go)

The independent command `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/compress/flate/encoder-cost/review` passed seven Lean and six native fixtures and rejected seven altered outcomes. The controls include a modular overflow witness at width 32 and preserve the native 64-bit result. Native 32-bit execution remains outside scope. [Reviewer control](../../../difftest/compress/flate/encoder-cost/review/main.go)
