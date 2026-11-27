# Independent review of `go/types.Pointer.Underlying`

**Verdict: FAITHFUL**

Reviewed author commit `2bb9659` against Go 1.27.0 `src/go/types/pointer.go`. Returns the original pointer wrapped as Type, including a typed nil for a nil receiver.

Pinned Go and Lean matched 300 representative profile rows and 60,000 samples using seeds 17, 2718 and 65537. Four independent literal fixtures matched both. The reviewer control detected returning an untyped nil in 4 rows. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/types/container-types/review` after building `godiff-go-types-container-types`.

Lean uses caller-provided identities to represent identity and reference preservation. The pinned native driver separately checks fresh allocations, typed-nil inputs and receiver behavior.

Evidence SHA-256: `b1fb4581c49b83c3fc65b9ddb8c1169997f48031057c8ee0c00af7af254b4a8a`
