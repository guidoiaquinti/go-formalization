# Independent review of `go/constant.i64toi`

**Verdict: FAITHFUL**

Reviewed author commit `342afee` against Go 1.27.0 `src/go/constant/value.go:273–283`. Converts the full signed word exactly and allocates a fresh big.Int.

Pinned Go and Lean matched 140 fixed rows and 60,000 sampled rows at seeds 17, 2718 and 65537. Four independent literal cases matched both. The control detected narrowing to int32 in 2 cases. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/constant/numeric-conversions/review` after building `godiff-go-constant-numeric-conversions`.

Lean represents exact values, float precision and opaque identity. Native checks actual fresh allocation, rounding mode, zero sign and interface identity. Integer-to-rational conversion uses a nonnil initialized source.

Evidence SHA-256: `b2e2b13fb6040f32374c015c123b3dac3d600f6451d236fa984fd68fe3a94a25`
