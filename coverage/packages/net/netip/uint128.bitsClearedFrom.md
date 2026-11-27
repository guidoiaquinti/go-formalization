# Independent review of `net/netip.uint128.bitsClearedFrom`

**Verdict: FAITHFUL**

Evidence SHA-256: `4075e31614e630d520056eb410c2c66e21b6f4c2da4926d9cfa6d7c9b4e2f9ca`

I checked evidence commit `ab8244f` against Go 1.27.0 `net/netip/uint128.go`. The source ANDs with the high-prefix mask. Bit 65 preserves exactly one low-limb bit position, and the fixture rejects using the mask complement instead. Bit 129 faults.

I ran `go run ./difftest/net/netip/uint128` on darwin/arm64 and reproduced 1,796 boundaries and 20,000 cases per seed 17, 2718 and 65537. The compiled Lean and direct pinned private Go observations agreed with the arbitrary-precision reference. Five separate literal fixtures passed both engines with `go run ./difftest/net/netip/uint128/review`, with eleven fixed mutation checks across the group. The scope records panic classification and direct field-address behavior where applicable. Exact runtime fault text and concurrent unsafe access are outside scope. [Independent control](../../../../difftest/net/netip/uint128/review/main.go)
