# Independent review of `go/constant.vtoc`

**Verdict: FAITHFUL**

Reviewed author commit `342afee` against Go 1.27.0 `src/go/constant/value.go:273–283`. Preserves the exact real value interface, including nil and typed nil. Native also checks the imaginary component is the compact int64Val zero.

Pinned Go and Lean matched 140 fixed rows and 60,000 sampled rows at seeds 17, 2718 and 65537. Four independent literal cases matched both. The control detected losing the real interface in 3 cases. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/constant/numeric-conversions/review` after building `godiff-go-constant-numeric-conversions`.

Lean represents exact values, float precision and opaque identity. Native checks actual fresh allocation, rounding mode, zero sign and interface identity. Integer-to-rational conversion uses a nonnil initialized source.

Evidence SHA-256: `de0ee6edc148e4861c1302e86e0e7ad2344ee6bde403a45a23351c38474b1b5e`
