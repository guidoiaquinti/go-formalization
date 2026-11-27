# Independent review of `go/constant.unknownVal.Kind`

**Verdict: FAITHFUL**

Reviewed author commit `c483761` against Go 1.27.0 `src/go/constant/value.go`. Returns Unknown with numeric code zero. The pinned method and Lean theorem agree on the single receiver value.

Author evidence matched 150 fixed and 60,000 seeded rows at seeds 17, 2718 and 65537. Nine independent literal fixtures matched pinned Go and compiled Lean. The control rejected returning code one on 9 fixtures. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/constant/scalars/review` from the repository root after building `godiff-go-constant-scalars`. The replay used native 64-bit int codes.

Evidence SHA-256: `8eebc9682668cff1b397576a76c2f32b8374becc4ef267cf012edb541b6dcfd2`
