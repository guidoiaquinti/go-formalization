# Independent review of `go/types.NewMap`

**Verdict: FAITHFUL**

Reviewed author commit `2bb9659` against Go 1.27.0 `src/go/types/map.go`. Each call allocates a fresh identity and stores both interfaces in their original positions, even for unusual keys.

Pinned Go and Lean matched 300 representative profile rows and 60,000 samples using seeds 17, 2718 and 65537. Four independent literal fixtures matched both. The reviewer control detected storing the element in the key field in 4 rows. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/types/container-types/review` after building `godiff-go-types-container-types`.

Lean uses caller-provided identities to represent identity and reference preservation. The pinned native driver separately checks fresh allocations, typed-nil inputs and receiver behavior.

Evidence SHA-256: `6a3dae0f6d246f1136d362a4787f83f063d68ea5addb92b7d8b818abe663f84b`
