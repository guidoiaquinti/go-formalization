# Independent review of `net/netip.mask6`

**Verdict: FAITHFUL**

Evidence SHA-256: `532122198a18205588949172e6dad8f9f4ca3d26d75973cdf4e21a3cfa949115`

I checked evidence commit `ab8244f` against Go 1.27.0 `net/netip/uint128.go`. The two pinned shifts construct a high-prefix mask across the limb boundary. The literal result at 64 rejects reversed limbs, 65 isolates the top low-word bit, and out-of-range signed lengths produce a negative-shift panic.

I ran `go run ./difftest/net/netip/uint128` on darwin/arm64 and reproduced 1,796 boundaries and 20,000 cases per seed 17, 2718 and 65537. The compiled Lean and direct pinned private Go observations agreed with the arbitrary-precision reference. Five separate literal fixtures passed both engines with `go run ./difftest/net/netip/uint128/review`, with eleven fixed mutation checks across the group. The scope records panic classification and direct field-address behavior where applicable. Exact runtime fault text and concurrent unsafe access are outside scope. [Independent control](../../../../difftest/net/netip/uint128/review/main.go)
