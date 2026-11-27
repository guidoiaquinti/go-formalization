# Independent review of `go/constant.unknownVal.String`

**Verdict: FAITHFUL**

Reviewed author commit `c483761` against Go 1.27.0 `src/go/constant/value.go`. Returns the lowercase literal unknown for the single receiver value.

Author evidence matched 150 fixed and 60,000 seeded rows at seeds 17, 2718 and 65537. Nine independent literal fixtures matched pinned Go and compiled Lean. The control rejected capitalizing unknown on 9 fixtures. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/constant/scalars/review` from the repository root after building `godiff-go-constant-scalars`. The replay used native 64-bit int codes.

Evidence SHA-256: `a5c14c7649c5f00684f1367dd44ac1f7963490c3ff8439cc51fd457776a661e8`
