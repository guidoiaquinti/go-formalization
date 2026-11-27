# Independent review of `go/constant.unknownVal.ExactString`

**Verdict: FAITHFUL**

Reviewed author commit `c483761` against Go 1.27.0 `src/go/constant/value.go`. Delegates to the ordinary unknown rendering and returns the same lowercase literal.

Author evidence matched 150 fixed and 60,000 seeded rows at seeds 17, 2718 and 65537. Nine independent literal fixtures matched pinned Go and compiled Lean. The control rejected returning an empty exact string on 9 fixtures. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/constant/scalars/review` from the repository root after building `godiff-go-constant-scalars`. The replay used native 64-bit int codes.

Evidence SHA-256: `39b003cb345c00b95dcf4d079b7707fb900fbad877ab3ae573944118c2a9c751`
