# Independent review of `go/types.Var.Kind`

**Verdict: FAITHFUL**

Reviewed author commit `f0de01f` against Go 1.27.0 `src/go/types/object.go`. Reads the raw kind byte without validating enum membership. A nil receiver faults.

Pinned native and Lean matched 1,536 fixed and 60,000 sampled rows at seeds 17, 2718 and 65537. Four independent literal fixtures matched both. The control detected returning zero unconditionally in 4 fixtures. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/types/variable-state/review` after building `godiff-go-types-variable-state`.

Lean uses abstract pointer and type identity tokens. Native compares the actual constructor and origin pointers, exact type interfaces including typed nil, raw names and nil outcomes.

Evidence SHA-256: `edf2a9a5132abeee2aebc5d33153a410112ab5e43c97874ddcf83f2697b9fbcb`
