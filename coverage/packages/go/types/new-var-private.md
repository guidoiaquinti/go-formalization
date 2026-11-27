# Independent review of `go/types.newVar`

**Verdict: FAITHFUL**

Reviewed author commit `f0de01f` against Go 1.27.0 `src/go/types/object.go`. Stores the exact supplied kind, raw name, position, package and type interface. Native also checks fresh receiver identity and zeroed metadata.

Pinned native and Lean matched 1,536 fixed and 60,000 sampled rows at seeds 17, 2718 and 65537. Four independent literal fixtures matched both. The control detected forcing PackageVar in 4 fixtures. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/types/variable-state/review` after building `godiff-go-types-variable-state`.

Lean uses abstract pointer and type identity tokens. Native compares the actual constructor and origin pointers, exact type interfaces including typed nil, raw names and nil outcomes.

Evidence SHA-256: `9c4bada59e5104b1b06151d75c7a318921573e4e16952e3a28de409a8faa079f`
