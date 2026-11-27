# Independent review of `go/types.Pointer.Elem`

**Verdict: FAITHFUL**

Reviewed author commit `2bb9659` against Go 1.27.0 `src/go/types/pointer.go`. Reads the existing base field and faults on a nil receiver. A zero pointer object yields the nil interface.

Pinned Go and Lean matched 300 representative profile rows and 60,000 samples using seeds 17, 2718 and 65537. Four independent literal fixtures matched both. The reviewer control detected always returning a nil interface in 3 rows. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/types/container-types/review` after building `godiff-go-types-container-types`.

Lean uses caller-provided identities to represent identity and reference preservation. The pinned native driver separately checks fresh allocations, typed-nil inputs and receiver behavior.

Evidence SHA-256: `4d1ea3d779782b463ae7f686b19ecd8934dbd7de8b00608c79736fdf6a5cd343`
