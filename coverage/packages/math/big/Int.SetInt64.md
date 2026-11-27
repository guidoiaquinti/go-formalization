# Independent review of `math/big.Int.SetInt64`

**Verdict: FAITHFUL**

Evidence SHA-256: `8e7145248aed2390d90b83c75615b758bdec0fed9d25010f095b1032fd393603`

I reviewed evidence commit `4a6cdff` against Go 1.27.0 `src/math/big/int.go`. Pinned Go stores the normalized magnitude and replaces the sign, returning the receiver. The 32-bit minimum theorem checks the zero low word and high sign-bit magnitude. The native fixture starts from a negative receiver and verifies that setting zero clears that sign. Storage reuse remains outside the model.

I ran `go run ./difftest/math/big/int-words` on darwin/arm64. It passed 1,428 boundaries and 20,000 cases per seed for 17, 2718 and 65537 on each modeled width. Native replay ran at 64 bits. The independent control at `go run ./difftest/math/big/int-words/review` checked five native cases and five Lean cases per width and rejected ten incorrect observations. Magnitudes are normalized. Native 32-bit replay, allocation capacity and timing are outside this review. [Independent control](../../../../difftest/math/big/int-words/review/main.go)
