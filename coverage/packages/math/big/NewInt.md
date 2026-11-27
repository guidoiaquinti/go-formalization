# Independent review of `math/big.NewInt`

**Verdict: FAITHFUL**

Evidence SHA-256: `43df2ed21949a4cf4b572a6329052a008bd36fbe2dca1eaa03720439000a1bc4`

I reviewed evidence commit `4a6cdff` against Go 1.27.0 `src/math/big/int.go`. The source creates a new Int and uses one word or two 32-bit words for a nonzero magnitude. Its unsigned negation handles the signed minimum. The theorem proves the normalized zero case. The independent native control checks distinct object and nonempty backing addresses for simultaneous constructor results.

I ran `go run ./difftest/math/big/int-words` on darwin/arm64. It passed 1,428 boundaries and 20,000 cases per seed for 17, 2718 and 65537 on each modeled width. Native replay ran at 64 bits. The independent control at `go run ./difftest/math/big/int-words/review` checked five native cases and five Lean cases per width and rejected ten incorrect observations. Magnitudes are normalized. Native 32-bit replay, allocation capacity and timing are outside this review. [Independent control](../../../../difftest/math/big/int-words/review/main.go)
