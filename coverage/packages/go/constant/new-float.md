# Independent review of `go/constant.newFloat`

**Verdict: FAITHFUL**

Reviewed author commit `342afee` against Go 1.27.0 `src/go/constant/value.go:273–283`. Allocates positive zero at 512 bits with default round-to-nearest-even and exact accuracy. Native checks the sign, mode and precision.

Pinned Go and Lean matched 140 fixed rows and 60,000 sampled rows at seeds 17, 2718 and 65537. Four independent literal cases matched both. The control detected substituting binary64 precision in 4 cases. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/constant/numeric-conversions/review` after building `godiff-go-constant-numeric-conversions`.

Lean represents exact values, float precision and opaque identity. Native checks actual fresh allocation, rounding mode, zero sign and interface identity. Integer-to-rational conversion uses a nonnil initialized source.

Evidence SHA-256: `24936db9d9842eb6fc6c54b29cf2134e35bad9ef31a8aee3afbcb2a8b28387ba`
