# Independent review of `go/types.Var.SetKind`

**Verdict: FAITHFUL**

Reviewed author commit `f0de01f` against Go 1.27.0 `src/go/types/object.go`. Changes only the kind byte. Native compares the remaining object, origin and embedding state, including synthetic origins. A nil receiver faults.

Pinned native and Lean matched 1,536 fixed and 60,000 sampled rows at seeds 17, 2718 and 65537. Four independent literal fixtures matched both. The control detected skipping the assignment in 3 fixtures. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/types/variable-state/review` after building `godiff-go-types-variable-state`.

Lean uses abstract pointer and type identity tokens. Native compares the actual constructor and origin pointers, exact type interfaces including typed nil, raw names and nil outcomes.

Evidence SHA-256: `bccdcb3991251798507c0a1b7bda8c58aa2ea799092d43774d5e4e7791854f14`
