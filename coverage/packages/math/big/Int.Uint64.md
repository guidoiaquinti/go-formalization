# Independent review of `math/big.Int.Uint64`

**Verdict: FAITHFUL**

Evidence SHA-256: `eee00fbc2f5bca21a1929cc4939be5757e74d17b99b0a1a5b046ed2c37a20660`

I reviewed evidence commit `4a6cdff` against Go 1.27.0 `src/math/big/int.go`. Pinned source reconstructs the low 64 magnitude bits directly. The theorem and independent fixed negative case distinguish magnitude from two's-complement conversion. Results for unrepresentable values describe pinned Go behavior.

I ran `go run ./difftest/math/big/int-words` on darwin/arm64. It passed 1,428 boundaries and 20,000 cases per seed for 17, 2718 and 65537 on each modeled width. Native replay ran at 64 bits. The independent control at `go run ./difftest/math/big/int-words/review` checked five native cases and five Lean cases per width and rejected ten incorrect observations. Magnitudes are normalized. Native 32-bit replay, allocation capacity and timing are outside this review. [Independent control](../../../../difftest/math/big/int-words/review/main.go)
