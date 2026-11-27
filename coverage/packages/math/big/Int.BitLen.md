# Independent review of `math/big.Int.BitLen`

**Verdict: FAITHFUL**

Evidence SHA-256: `7d5bfc8e79183fbc88554c15f48d6a7e76a00938f2cf4b5e84ff708003aa5609`

I reviewed evidence commit `4a6cdff` against Go 1.27.0 `src/math/big/int.go`. The model counts full lower words plus the significant bits in the highest word, using normalized magnitude. The theorem checks 64 bits for the minimum int64 magnitude in 32-bit layout. A 65-bit independent fixture rejects counting only the highest word.

I ran `go run ./difftest/math/big/int-words` on darwin/arm64. It passed 1,428 boundaries and 20,000 cases per seed for 17, 2718 and 65537 on each modeled width. Native replay ran at 64 bits. The independent control at `go run ./difftest/math/big/int-words/review` checked five native cases and five Lean cases per width and rejected ten incorrect observations. Magnitudes are normalized. Native 32-bit replay, allocation capacity and timing are outside this review. [Independent control](../../../../difftest/math/big/int-words/review/main.go)
