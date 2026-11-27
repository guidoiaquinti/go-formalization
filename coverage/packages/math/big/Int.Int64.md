# Independent review of `math/big.Int.Int64`

**Verdict: FAITHFUL**

Evidence SHA-256: `008424070a1e7b9bcf55855aea6fd731fd0da81831da2e81057b4bbcc32a37d8`

I reviewed evidence commit `4a6cdff` against Go 1.27.0 `src/math/big/int.go`. Pinned source reconstructs the low 64 magnitude bits, converts to int64 and applies the stored sign. The -1 theorem checks negation and the native fixture with magnitude 2^64+1 distinguishes low-word truncation and sign application. Results for unrepresentable values describe pinned Go behavior.

I ran `go run ./difftest/math/big/int-words` on darwin/arm64. It passed 1,428 boundaries and 20,000 cases per seed for 17, 2718 and 65537 on each modeled width. Native replay ran at 64 bits. The independent control at `go run ./difftest/math/big/int-words/review` checked five native cases and five Lean cases per width and rejected ten incorrect observations. Magnitudes are normalized. Native 32-bit replay, allocation capacity and timing are outside this review. [Independent control](../../../../difftest/math/big/int-words/review/main.go)
