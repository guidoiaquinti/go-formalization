# Independent review of `math/big.Int.IsInt64`

**Verdict: FAITHFUL**

Evidence SHA-256: `9269d74633e44dfbc5b188ea762322b449c8a31191337921ba2b13ba4c622278`

I reviewed evidence commit `4a6cdff` against Go 1.27.0 `src/math/big/int.go`. For a magnitude occupying at most 64 native bits, pinned Go accepts a clear sign bit or the negative minimum. The theorem and control distinguish the two signs at exactly 2^63. The larger-magnitude fixture rejects accidental acceptance after discarding high words.

I ran `go run ./difftest/math/big/int-words` on darwin/arm64. It passed 1,428 boundaries and 20,000 cases per seed for 17, 2718 and 65537 on each modeled width. Native replay ran at 64 bits. The independent control at `go run ./difftest/math/big/int-words/review` checked five native cases and five Lean cases per width and rejected ten incorrect observations. Magnitudes are normalized. Native 32-bit replay, allocation capacity and timing are outside this review. [Independent control](../../../../difftest/math/big/int-words/review/main.go)
