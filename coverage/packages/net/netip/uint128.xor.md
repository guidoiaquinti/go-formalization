# Independent review of `net/netip.uint128.xor`

**Verdict: FAITHFUL**

Evidence SHA-256: `cc570c22557c520d121c16034aa07489110624b7ef14c8e3a180d67603ec2903`

I checked evidence commit `ab8244f` against Go 1.27.0 `net/netip/uint128.go`. XOR is applied to corresponding limbs. Overlapping one bits in the high limb reject accidental OR.

I ran `go run ./difftest/net/netip/uint128` on darwin/arm64 and reproduced 1,796 boundaries and 20,000 cases per seed 17, 2718 and 65537. The compiled Lean and direct pinned private Go observations agreed with the arbitrary-precision reference. Five separate literal fixtures passed both engines with `go run ./difftest/net/netip/uint128/review`, with eleven fixed mutation checks across the group. The scope records panic classification and direct field-address behavior where applicable. Exact runtime fault text and concurrent unsafe access are outside scope. [Independent control](../../../../difftest/net/netip/uint128/review/main.go)
