# Independent review of `net/netip.uint128.or`

**Verdict: FAITHFUL**

Evidence SHA-256: `dd64447b8dda8ef331186d55d227908818b80ce49f3efe2388e3504ae2ded4b5`

I checked evidence commit `ab8244f` against Go 1.27.0 `net/netip/uint128.go`. OR is applied to corresponding limbs. Overlapping one bits in the high limb reject accidental XOR.

I ran `go run ./difftest/net/netip/uint128` on darwin/arm64 and reproduced 1,796 boundaries and 20,000 cases per seed 17, 2718 and 65537. The compiled Lean and direct pinned private Go observations agreed with the arbitrary-precision reference. Five separate literal fixtures passed both engines with `go run ./difftest/net/netip/uint128/review`, with eleven fixed mutation checks across the group. The scope records panic classification and direct field-address behavior where applicable. Exact runtime fault text and concurrent unsafe access are outside scope. [Independent control](../../../../difftest/net/netip/uint128/review/main.go)
