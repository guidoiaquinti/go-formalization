# Independent review of `go/types.NewSlice`

**Verdict: FAITHFUL**

Reviewed author commit `2bb9659` against Go 1.27.0 `src/go/types/slice.go`. Each call allocates a fresh identity and stores the precise element interface.

Pinned Go and Lean matched 300 representative profile rows and 60,000 samples using seeds 17, 2718 and 65537. Four independent literal fixtures matched both. The reviewer control detected reusing the same slice identity in 4 rows. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/types/container-types/review` after building `godiff-go-types-container-types`.

Lean uses caller-provided identities to represent identity and reference preservation. The pinned native driver separately checks fresh allocations, typed-nil inputs and receiver behavior.

Evidence SHA-256: `c2efad5c9e480b5c024796adf52359bbd6dcc7459e9bd8d0aeec410ed2d6601c`
