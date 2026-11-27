# Independent review of `compress/flate.newHuffmanEncoder`

**Verdict: FAITHFUL**

Evidence SHA-256: `60c5443957ac48b10f3b8248f74aebe73a8af4c46c8f648dc4a64412a5b507f8`

The source rounds capacity using the bit length of uint32(size−1). The fixed native controls distinguish exact powers from just-over powers, confirm zero code and auxiliary storage, and check rejection of negative lengths. Zero on a 64-bit Go profile computes capacity 2^32 and was checked only in Lean to avoid that allocation. Synthetic 32-bit capacity overflow and 64-bit uint32 truncation cases were also checked only in Lean. Runtime allocation failure and exact allocation identity are outside scope. The pinned Go and Lean comparison passed 64,104 author rows at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Compress/Flate/EncoderInit.lean) · [Comparison](../../../difftest/compress/flate/encoder-init/check.go)

The independent command `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/compress/flate/encoder-init/review` passed ten Lean and six native fixtures and rejected ten altered results across the group. Native execution used darwin/arm64. [Reviewer control](../../../difftest/compress/flate/encoder-init/review/main.go)
