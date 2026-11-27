# Independent review of `go/constant.newRat`

**Verdict: FAITHFUL**

Reviewed author commit `342afee` against Go 1.27.0 `src/go/constant/value.go:273–283`. Allocates a fresh rational at mathematical zero with denominator one.

Pinned Go and Lean matched 140 fixed rows and 60,000 sampled rows at seeds 17, 2718 and 65537. Four independent literal cases matched both. The control detected returning a zero denominator in 4 cases. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/constant/numeric-conversions/review` after building `godiff-go-constant-numeric-conversions`.

Lean represents exact values, float precision and opaque identity. Native checks actual fresh allocation, rounding mode, zero sign and interface identity. Integer-to-rational conversion uses a nonnil initialized source.

Evidence SHA-256: `6e473783b4ed9db1127887041cc51c3294d58d4351c1d31f7cff9a0f2e606cf0`
