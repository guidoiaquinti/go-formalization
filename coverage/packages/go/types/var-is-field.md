# Independent review of `go/types.Var.IsField`

**Verdict: FAITHFUL**

Reviewed author commit `f0de01f` against Go 1.27.0 `src/go/types/object.go`. Tests exact equality with FieldVar. A nil receiver faults.

Pinned native and Lean matched 1,536 fixed and 60,000 sampled rows at seeds 17, 2718 and 65537. Four independent literal fixtures matched both. The control detected testing the embedded flag instead in 3 fixtures. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/types/variable-state/review` after building `godiff-go-types-variable-state`.

Lean uses abstract pointer and type identity tokens. Native compares the actual constructor and origin pointers, exact type interfaces including typed nil, raw names and nil outcomes.

Evidence SHA-256: `ca046ec1f608649bfff98477363f30154c8cdc885bbe531f3177b46e67b7c64c`
