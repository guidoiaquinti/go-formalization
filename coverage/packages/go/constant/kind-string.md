# Independent review of `go/constant.Kind.String`

**Verdict: FAITHFUL**

Reviewed author commit `c483761` against Go 1.27.0 `src/go/constant/kind_string.go`. Uses the generated name offsets for codes zero through five, and a signed decimal fallback outside that interval. Theorems check the table and negative fallback. Direct tests cover the above-table branch and the native minimum and maximum.

Author evidence matched 150 fixed and 60,000 seeded rows at seeds 17, 2718 and 65537. Nine independent literal fixtures matched pinned Go and compiled Lean. The control rejected always returning the Unknown name on 8 fixtures. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/constant/scalars/review` from the repository root after building `godiff-go-constant-scalars`. The replay used native 64-bit int codes.

Evidence SHA-256: `a028c4351bb5d6c4b4c2073560086b6be8672508107ebca80e86c02f61ed5e9d`
