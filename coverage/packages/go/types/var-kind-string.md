# Independent review of `go/types.VarKind.String`

**Verdict: FAITHFUL**

Reviewed author commit `f0de01f` against Go 1.27.0 `src/go/types/object.go`. Uses fixed names at codes one through six and unsigned decimal in VarKind(n) elsewhere. Native exhausts the uint8 domain.

Pinned native and Lean matched 1,536 fixed and 60,000 sampled rows at seeds 17, 2718 and 65537. Four independent literal fixtures matched both. The control detected returning VarKind(0) unconditionally in 3 fixtures. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/types/variable-state/review` after building `godiff-go-types-variable-state`.

Lean uses abstract pointer and type identity tokens. Native compares the actual constructor and origin pointers, exact type interfaces including typed nil, raw names and nil outcomes.

Evidence SHA-256: `9eca1f69628b7b3181cadf34df3a673fe2990c673319364a69a482eaa9b50264`
