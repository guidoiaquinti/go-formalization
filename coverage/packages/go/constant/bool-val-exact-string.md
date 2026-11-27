# Independent review of `go/constant.boolVal.ExactString`

**Verdict: FAITHFUL**

Reviewed author commit `c483761` against Go 1.27.0 `src/go/constant/value.go`. Delegates to the ordinary Boolean rendering. Both Boolean values retain the lowercase spelling.

Author evidence matched 150 fixed and 60,000 seeded rows at seeds 17, 2718 and 65537. Nine independent literal fixtures matched pinned Go and compiled Lean. The control rejected returning the uppercase TRUE literal on 9 fixtures. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/constant/scalars/review` from the repository root after building `godiff-go-constant-scalars`. The replay used native 64-bit int codes.

Evidence SHA-256: `349c33fa755d30c4883a7db6a2cf18869f02d18e2e0d022a8e7117d5283e94f3`
