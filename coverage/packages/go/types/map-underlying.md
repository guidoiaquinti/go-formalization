# Independent review of `go/types.Map.Underlying`

**Verdict: FAITHFUL**

Reviewed author commit `2bb9659` against Go 1.27.0 `src/go/types/map.go`. Returns the original map pointer wrapped as Type, including a typed nil for a nil receiver.

Pinned Go and Lean matched 300 representative profile rows and 60,000 samples using seeds 17, 2718 and 65537. Four independent literal fixtures matched both. The reviewer control detected returning an untyped nil in 4 rows. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/types/container-types/review` after building `godiff-go-types-container-types`.

Lean uses caller-provided identities to represent identity and reference preservation. The pinned native driver separately checks fresh allocations, typed-nil inputs and receiver behavior.

Evidence SHA-256: `a0e12abb83d7592e2ad11fc77a62c2a0c8a8cbf3926e1914f0f0b46b0ce10e09`
