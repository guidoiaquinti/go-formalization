# Independent review of `math/big.Int.Sign`

**Verdict: FAITHFUL**

Evidence SHA-256: `119e31c4f253533b9149d41e1ccf12511777780a1456b359d0359dcaa934336c`

I reviewed evidence commit `4a6cdff` against Go 1.27.0 `src/math/big/int.go`. The source tests for empty magnitude before reading the sign. The theorem maps an empty magnitude to zero regardless of its sign flag. The native control exercises canonical zero, positive and negative states.

I ran `go run ./difftest/math/big/int-words` on darwin/arm64. It passed 1,428 boundaries and 20,000 cases per seed for 17, 2718 and 65537 on each modeled width. Native replay ran at 64 bits. The independent control at `go run ./difftest/math/big/int-words/review` checked five native cases and five Lean cases per width and rejected ten incorrect observations. Magnitudes are normalized. Native 32-bit replay, allocation capacity and timing are outside this review. [Independent control](../../../../difftest/math/big/int-words/review/main.go)
