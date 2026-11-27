# Independent review of `net/netip.uint128.not`

**Verdict: FAITHFUL**

Evidence SHA-256: `41aa8c56938c2d419e7cc83ef3277f896595da7b23339271bd5c3d76f3d23921`

I checked evidence commit `ab8244f` against Go 1.27.0 `net/netip/uint128.go`. Complement acts on each full unsigned limb. An all-zero low limb rejects a mutation leaving the low limb unchanged.

I ran `go run ./difftest/net/netip/uint128` on darwin/arm64 and reproduced 1,796 boundaries and 20,000 cases per seed 17, 2718 and 65537. The compiled Lean and direct pinned private Go observations agreed with the arbitrary-precision reference. Five separate literal fixtures passed both engines with `go run ./difftest/net/netip/uint128/review`, with eleven fixed mutation checks across the group. The scope records panic classification and direct field-address behavior where applicable. Exact runtime fault text and concurrent unsafe access are outside scope. [Independent control](../../../../difftest/net/netip/uint128/review/main.go)
