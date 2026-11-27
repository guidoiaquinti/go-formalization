# Independent review of `go/types.NewPointer`

**Verdict: FAITHFUL**

Reviewed author commit `2bb9659` against Go 1.27.0 `src/go/types/pointer.go`. Each call allocates a fresh identity and stores the precise base interface.

Pinned Go and Lean matched 300 representative profile rows and 60,000 samples using seeds 17, 2718 and 65537. Four independent literal fixtures matched both. The reviewer control detected reusing the same pointer identity in 4 rows. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/types/container-types/review` after building `godiff-go-types-container-types`.

Lean uses caller-provided identities to represent identity and reference preservation. The pinned native driver separately checks fresh allocations, typed-nil inputs and receiver behavior.

Evidence SHA-256: `691c9cd313aeccdbbd0f26e791b4f62a6c85cfb3ca948bee573c5ea749f9f11e`
