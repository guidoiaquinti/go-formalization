# Independent review of `go/constant.i64tor`

**Verdict: FAITHFUL**

Reviewed author commit `342afee` against Go 1.27.0 `src/go/constant/value.go:273–283`. Converts the full signed word to an exact fresh rational with denominator one.

Pinned Go and Lean matched 140 fixed rows and 60,000 sampled rows at seeds 17, 2718 and 65537. Four independent literal cases matched both. The control detected treating the input as unsigned in 2 cases. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/constant/numeric-conversions/review` after building `godiff-go-constant-numeric-conversions`.

Lean represents exact values, float precision and opaque identity. Native checks actual fresh allocation, rounding mode, zero sign and interface identity. Integer-to-rational conversion uses a nonnil initialized source.

Evidence SHA-256: `b79e4bf8ecb979c1c65c82b1bec0b64e022aee2173faa08966035665af96779e`
