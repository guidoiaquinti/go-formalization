# Independent review of `go/types.Map.Key`

**Verdict: FAITHFUL**

Reviewed author commit `2bb9659` against Go 1.27.0 `src/go/types/map.go`. Reads the stored key interface. The paired distinct key and element fixtures check the field selection. Nil receiver faults and zero object returns a nil Type.

Pinned Go and Lean matched 300 representative profile rows and 60,000 samples using seeds 17, 2718 and 65537. Four independent literal fixtures matched both. The reviewer control detected reading the element field instead in 2 rows. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/types/container-types/review` after building `godiff-go-types-container-types`.

Lean uses caller-provided identities to represent identity and reference preservation. The pinned native driver separately checks fresh allocations, typed-nil inputs and receiver behavior.

Evidence SHA-256: `1ebb1c7d30d7606d549ce4e44abfa481abd971d5ef8d63a8b2137874f3ac2d97`
