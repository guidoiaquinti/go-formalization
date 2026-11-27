# Independent review of `net/netip.uint128.isZero`

**Verdict: FAITHFUL**

Evidence SHA-256: `c8f021d3b68352a4c79c62ecaee7fa7f1654a4b0ddfbe2d753d27f1b66916fb7`

I checked evidence commit `ab8244f` against Go 1.27.0 `net/netip/uint128.go`. Both limbs must be zero. A nonzero high limb with zero low limb rejects checking only the low word.

I ran `go run ./difftest/net/netip/uint128` on darwin/arm64 and reproduced 1,796 boundaries and 20,000 cases per seed 17, 2718 and 65537. The compiled Lean and direct pinned private Go observations agreed with the arbitrary-precision reference. Five separate literal fixtures passed both engines with `go run ./difftest/net/netip/uint128/review`, with eleven fixed mutation checks across the group. The scope records panic classification and direct field-address behavior where applicable. Exact runtime fault text and concurrent unsafe access are outside scope. [Independent control](../../../../difftest/net/netip/uint128/review/main.go)
